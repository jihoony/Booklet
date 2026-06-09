// Types for Wails
declare global {
    interface Window {
        go: any;
        runtime: any;
    }
}

// State
let inputPath = "";
let nValue = 4;

// Prevent default browser behavior for drag & drop navigation (prevents opening files in webview)
window.addEventListener('dragover', (e) => {
    e.preventDefault();
}, false);

window.addEventListener('drop', (e) => {
    e.preventDefault();
}, false);

// Initialization
async function init() {
    console.log("Booklet Pro: Initializing...");
    
    const elements = {
        dropZone: document.getElementById('dropZone'),
        fileList: document.getElementById('fileList'),
        layoutGrid: document.getElementById('layoutGrid'),
        createBtn: document.getElementById('createBtn') as HTMLButtonElement,
        formSize: document.getElementById('formSize') as HTMLSelectElement,
        guides: document.getElementById('guides') as HTMLInputElement,
        binding: document.getElementById('binding') as HTMLInputElement,
    };

    // 요소 존재 확인
    for (const [key, el] of Object.entries(elements)) {
        if (!el) {
            console.error(`Booklet Pro: Element '${key}' not found!`);
            return; 
        }
    }

    const { dropZone, fileList, layoutGrid, createBtn, formSize, guides, binding } = elements as any;

    // Handle Wails native drag and drop
    const setupDragAndDrop = () => {
        if (window.runtime && window.runtime.OnFileDrop) {
            window.runtime.OnFileDrop((_x: number, _y: number, paths: string[]) => {
                if (paths && paths.length > 0) {
                    const droppedPath = paths[0];
                    if (droppedPath.toLowerCase().endsWith('.pdf')) {
                        inputPath = droppedPath;
                        updateFileList(fileList);
                        createBtn.disabled = !inputPath;
                    } else {
                        alert("PDF 파일만 지원됩니다.");
                    }
                }
            }, true);
            console.log("Booklet Pro: Wails native OnFileDrop registered.");
        } else if (window.runtime) {
            window.runtime.EventsOn("wails:file-drop", (_x: number, _y: number, paths: string[]) => {
                if (paths && paths.length > 0) {
                    const droppedPath = paths[0];
                    if (droppedPath.toLowerCase().endsWith('.pdf')) {
                        inputPath = droppedPath;
                        updateFileList(fileList);
                        createBtn.disabled = !inputPath;
                    } else {
                        alert("PDF 파일만 지원됩니다.");
                    }
                }
            });
            console.log("Booklet Pro: Wails wails:file-drop event registered.");
        } else {
            console.log("Booklet Pro: Wails runtime not ready yet, retrying drag & drop registration in 100ms...");
            setTimeout(setupDragAndDrop, 100);
        }
    };
    setupDragAndDrop();

    // Drop Zone Events
    dropZone.addEventListener('click', async () => {
        try {
            if (!window.go) {
                alert("Wails runtime not loaded yet.");
                return;
            }
            const result = await window.go.main.App.SelectFile();
            if (result) {
                inputPath = result;
                updateFileList(fileList);
                createBtn.disabled = !inputPath;
            }
        } catch (err) {
            console.error("File selection error:", err);
        }
    });
    
    // Layout Selector
    const layoutOpts = layoutGrid.querySelectorAll('.layout-opt');
    layoutOpts.forEach((opt: any) => {
        opt.addEventListener('click', () => {
            layoutOpts.forEach((o: any) => o.classList.remove('selected'));
            opt.classList.add('selected');
            nValue = parseInt(opt.getAttribute('data-value') || '4');
            console.log("Layout changed to:", nValue);
        });
    });

    // Create Button
    createBtn.addEventListener('click', async () => {
        if (!inputPath) return;
        
        try {
            const outputPath = await window.go.main.App.SelectSaveFile(inputPath);
            if (!outputPath) return;

            const opts = {
                Input:      inputPath,
                Output:     outputPath,
                N:          nValue,
                FormSize:   formSize.value,
                Guides:     guides.checked,
                Binding:    binding.checked ? 'long' : 'short',
                BType:      "booklet",
                Multifolio: false,
                FolioSize:  6,
            };
            
            createBtn.disabled = true;
            createBtn.textContent = "Processing...";
            
            const result = await window.go.main.App.ProcessBooklet(opts);
            
            if (result === "Success") {
                const modal = document.getElementById('successModal') as HTMLDivElement;
                const openBtn = document.getElementById('modalOpenBtn') as HTMLButtonElement;
                const closeBtn = document.getElementById('modalCloseBtn') as HTMLButtonElement;

                if (modal && openBtn && closeBtn) {
                    modal.style.display = 'flex';
                    
                    // Open Folder handler
                    const handleOpen = async () => {
                        await window.go.main.App.OpenFolder(opts.Output);
                        modal.style.display = 'none';
                        cleanup();
                    };
                    
                    // Close handler
                    const handleClose = () => {
                        modal.style.display = 'none';
                        cleanup();
                    };

                    const cleanup = () => {
                        openBtn.removeEventListener('click', handleOpen);
                        closeBtn.removeEventListener('click', handleClose);
                    };

                    openBtn.addEventListener('click', handleOpen);
                    closeBtn.addEventListener('click', handleClose);
                }
            } else {
                alert(result);
            }
        } catch (err) {
            alert("Error: " + err);
        } finally {
            createBtn.disabled = false;
            createBtn.textContent = "Booklet 생성";
        }
    });

    console.log("Booklet Pro: Ready.");
    document.getElementById('loading-check')?.remove();
}

function updateFileList(fileList: HTMLElement) {
    if (!inputPath) {
        fileList.innerHTML = '<div style="color: var(--text-dim); text-align: center; font-size: 13px; padding: 20px;">No files selected</div>';
        return;
    }

    const fileName = inputPath.split(/[\\/]/).pop();
    fileList.innerHTML = `
        <div class="file-item">
            <div class="file-info">
                <svg width="16" height="16" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z"></path>
                </svg>
                <span class="file-name">${fileName}</span>
            </div>
            <span class="remove-btn" id="removeFile">&times;</span>
        </div>
    `;

    document.getElementById('removeFile')?.addEventListener('click', () => {
        inputPath = "";
        const createBtn = document.getElementById('createBtn') as HTMLButtonElement;
        updateFileList(fileList);
        if (createBtn) createBtn.disabled = true;
    });
}

// 초기화 실행
init().catch(err => console.error("Init failed:", err));
