// Types for Wails
declare global {
    interface Window {
        go: any;
        runtime: any;
    }
}

// State
let inputPath = "";
let nValue = 2;
let lastPageCount = 0;

function askUserForRecommendation(message: string): Promise<boolean> {
    const modal = document.getElementById('confirmModal') as HTMLDivElement;
    const msgEl = document.getElementById('confirmModalMsg') as HTMLParagraphElement;
    const approveBtn = document.getElementById('confirmApproveBtn') as HTMLButtonElement;
    const declineBtn = document.getElementById('confirmDeclineBtn') as HTMLButtonElement;

    if (!modal || !msgEl || !approveBtn || !declineBtn) {
        return Promise.resolve(false);
    }

    msgEl.innerHTML = message;
    modal.style.display = 'flex';

    return new Promise((resolve) => {
        const cleanup = (value: boolean) => {
            approveBtn.removeEventListener('click', handleApprove);
            declineBtn.removeEventListener('click', handleDecline);
            modal.style.display = 'none';
            resolve(value);
        };

        const handleApprove = () => cleanup(true);
        const handleDecline = () => cleanup(false);

        approveBtn.addEventListener('click', handleApprove);
        declineBtn.addEventListener('click', handleDecline);
    });
}

// Prevent default browser behavior for drag & drop navigation (prevents opening files in webview)
// Using capture phase (true) to intercept the event before webview default navigation can trigger.
window.addEventListener('dragenter', (e) => {
    e.preventDefault();
}, true);

window.addEventListener('dragover', (e) => {
    e.preventDefault();
}, true);

window.addEventListener('drop', (e) => {
    e.preventDefault();
}, true);

// Initialization
async function init() {
    console.log("Booklet Pro: Initializing...");
    
    const elements = {
        dropZone: document.getElementById('dropZone'),
        fileList: document.getElementById('fileList'),
        layoutGrid: document.getElementById('layoutGrid'),
        createBtn: document.getElementById('createBtn') as HTMLButtonElement,
        formSize: document.getElementById('formSize') as HTMLSelectElement,
        margin: document.getElementById('margin') as HTMLInputElement,
        btype: document.getElementById('btype') as HTMLSelectElement,
        multifolio: document.getElementById('multifolio') as HTMLInputElement,
        folioSize: document.getElementById('folioSize') as HTMLInputElement,
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

    const { dropZone, fileList, layoutGrid, createBtn, formSize, margin, btype, multifolio, folioSize, guides, binding } = elements as any;

    // Handle Wails native drag and drop
    const setupDragAndDrop = () => {
        if (window.runtime && window.runtime.OnFileDrop) {
            window.runtime.OnFileDrop((_x: number, _y: number, paths: string[]) => {
                if (paths && paths.length > 0) {
                    const droppedPath = paths[0];
                    if (droppedPath.toLowerCase().endsWith('.pdf')) {
                        handleFileSelected(droppedPath, fileList, createBtn);
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
                        handleFileSelected(droppedPath, fileList, createBtn);
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
                await handleFileSelected(result, fileList, createBtn);
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
            if (lastPageCount > 0) {
                applySmartRecommendations(lastPageCount);
            }
        });
    });

    // Toggle folio size group based on multifolio checkbox
    const folioSizeGroup = document.getElementById('folioSizeGroup') as HTMLDivElement;
    multifolio.addEventListener('change', () => {
        if (folioSizeGroup) {
            folioSizeGroup.style.display = multifolio.checked ? 'block' : 'none';
        }
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
                Margin:     parseFloat(margin.value || '10'),
                Binding:    binding.checked ? 'long' : 'short',
                BType:      btype.value,
                Multifolio: multifolio.checked,
                FolioSize:  parseInt(folioSize.value || '6'),
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

function updateFileList(fileList: HTMLElement, pageCount?: number) {
    if (!inputPath) {
        fileList.innerHTML = '<div style="color: var(--text-dim); text-align: center; font-size: 13px; padding: 20px;">No files selected</div>';
        return;
    }

    const fileName = inputPath.split(/[\\/]/).pop();
    const pageCountText = pageCount ? ` (${pageCount}페이지)` : '';
    fileList.innerHTML = `
        <div class="file-item">
            <div class="file-info">
                <svg width="16" height="16" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z"></path>
                </svg>
                <span class="file-name">${fileName}${pageCountText}</span>
            </div>
            <span class="remove-btn" id="removeFile">&times;</span>
        </div>
    `;

    document.getElementById('removeFile')?.addEventListener('click', () => {
        inputPath = "";
        lastPageCount = 0;
        const createBtn = document.getElementById('createBtn') as HTMLButtonElement;
        updateFileList(fileList);
        if (createBtn) createBtn.disabled = true;
    });
}

async function handleFileSelected(path: string, fileList: HTMLElement, createBtn: HTMLButtonElement) {
    inputPath = path;
    try {
        const pageCount = await window.go.main.App.GetPageCount(path);
        lastPageCount = pageCount;
        updateFileList(fileList, pageCount);
        await applySmartRecommendations(pageCount);
    } catch (err) {
        console.error("Failed to get page count:", err);
        lastPageCount = 0;
        updateFileList(fileList);
    }
    createBtn.disabled = !inputPath;
}

async function applySmartRecommendations(pageCount: number) {
    try {
        const btypeSelect = document.getElementById('btype') as HTMLSelectElement;
        const multifolioCheckbox = document.getElementById('multifolio') as HTMLInputElement;
        const folioSizeInput = document.getElementById('folioSize') as HTMLInputElement;
        const folioSizeGroup = document.getElementById('folioSizeGroup') as HTMLDivElement;
        const advancedDetails = document.querySelector('details.advanced-settings') as HTMLDetailsElement;

        if (!btypeSelect || !multifolioCheckbox || !folioSizeInput || !folioSizeGroup) {
            console.warn("고급 설정 엘리먼트를 DOM에서 찾을 수 없습니다.");
            return;
        }

        const pagesPerSheet = nValue * 2;
        if (pagesPerSheet <= 0) return;
        const totalSheets = Math.ceil(pageCount / pagesPerSheet);

        // 추천 옵션 결정
        let targetBType = "booklet";
        let targetMultifolio = false;
        let targetFolioSize = 6;
        let message = "";

        if (totalSheets <= 10) {
            targetBType = "booklet";
            targetMultifolio = false;
            if (btypeSelect.value === "booklet" && !multifolioCheckbox.checked) {
                return;
            }
            message = `선택하신 PDF는 <strong>${pageCount}페이지(종이 ${totalSheets}장)</strong>로 비교적 두껍지 않습니다.<br><br>가장 일반적인 형태인 <strong>단일 소책자(Saddle Stitch) 모드</strong>로 변환 설정을 변경하시겠습니까?`;
        } else if (totalSheets <= 30) {
            targetBType = "booklet";
            targetMultifolio = true;
            targetFolioSize = 6;
            if (btypeSelect.value === "booklet" && multifolioCheckbox.checked && parseInt(folioSizeInput.value) === 6) {
                return;
            }
            message = `선택하신 PDF는 <strong>${pageCount}페이지(종이 ${totalSheets}장)</strong>로 다소 두껍습니다.<br><br>접힘 불량 및 페이지 잘림을 방지하기 위해 6시트 단위의 <strong>시그니처 제본 모드(Multifolio)</strong>로 설정을 변경하시겠습니까?`;
        } else {
            targetBType = "perfectbound";
            targetMultifolio = false;
            if (btypeSelect.value === "perfectbound" && !multifolioCheckbox.checked) {
                return;
            }
            message = `선택하신 PDF는 <strong>${pageCount}페이지(종이 ${totalSheets}장)</strong>로 매우 두껍습니다.<br><br>반으로 접는 제본이 불가하므로 책등에 풀칠을 하는 <strong>무선 제본(Perfect Bound) 모드</strong>로 설정을 변경하시겠습니까?`;
        }

        const approved = await askUserForRecommendation(message);
        if (approved) {
            btypeSelect.value = targetBType;
            multifolioCheckbox.checked = targetMultifolio;
            folioSizeInput.value = targetFolioSize.toString();
            folioSizeGroup.style.display = targetMultifolio ? "block" : "none";
            if (advancedDetails) {
                advancedDetails.open = true;
            }
        }
    } catch (err) {
        console.error("applySmartRecommendations exception:", err);
    }
}

// 초기화 실행
init().catch(err => console.error("Init failed:", err));
