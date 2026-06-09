package booklet

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Options 소책자 변환 설정을 담는 구조체
type Options struct {
	Input      string
	Output     string
	N          int
	FormSize   string
	Guides     bool
	Margin     float64
	Binding    string
	BType      string
	Multifolio bool
	FolioSize  int
}

// buildConfigString 소책자 설정을 문자열로 빌드합니다. (유닛 테스트에 활용)
func buildConfigString(opts Options, pageCount int) string {
	multifolioStr := "off"
	if opts.Multifolio {
		multifolioStr = "on"
	} else {
		pagesPerSheet := opts.N * 2
		if pagesPerSheet > 0 {
			totalSheets := (pageCount + pagesPerSheet - 1) / pagesPerSheet
			// 종이 10장 이상이면 multifolio 자동 활성화
			if totalSheets > 10 {
				multifolioStr = "on"
			}
		}
	}

	// 옵션 문자열 조합
	descParts := []string{
		fmt.Sprintf("formsize:%s", opts.FormSize),
		"guides:off", // 가이드라인은 후처리로 직접 그림
		fmt.Sprintf("margin:%.1f", opts.Margin),
		fmt.Sprintf("binding:%s", opts.Binding),
		fmt.Sprintf("btype:%s", opts.BType),
	}

	if multifolioStr == "on" {
		descParts = append(descParts, "multifolio:on", fmt.Sprintf("foliosize:%d", opts.FolioSize))
	}

	return strings.Join(descParts, ", ")
}

// Process 실제 소책자 변환 프로세스를 실행합니다.
func Process(opts Options) error {
	// 파일 존재 여부 확인 겸 페이지 수 체크
	ctx, err := api.ReadContextFile(opts.Input)
	if err != nil {
		return fmt.Errorf("입력 파일 분석 실패: %v", err)
	}
	pageCount := ctx.PageCount

	desc := buildConfigString(opts, pageCount)

	// booklet 생성
	conf := model.NewDefaultConfiguration()
	nupVal, err := api.PDFBookletConfig(opts.N, desc, conf)
	if err != nil {
		return fmt.Errorf("설정 오류: %v", err)
	}

	err = api.BookletFile([]string{opts.Input}, opts.Output, nil, nupVal, conf)
	if err != nil {
		return fmt.Errorf("booklet 생성 실패: %v", err)
	}

	// 가이드라인 후처리
	if opts.Guides {
		if err := addGuides(opts.Output, opts.N, opts.Binding); err != nil {
			return fmt.Errorf("가이드라인 추가 실패: %v", err)
		}
	}

	return nil
}

func addGuides(pdfPath string, n int, binding string) error {
	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return fmt.Errorf("PDF 읽기 실패: %v", err)
	}

	if ctx.PageCount == 0 {
		return fmt.Errorf("페이지가 없습니다")
	}

	_, _, inhPAttrs, err := ctx.PageDict(1, false)
	if err != nil {
		return fmt.Errorf("페이지 정보 읽기 실패: %v", err)
	}
	mediaBox := inhPAttrs.MediaBox
	if mediaBox == nil {
		return fmt.Errorf("MediaBox를 찾을 수 없습니다")
	}

	width := mediaBox.Width()
	height := mediaBox.Height()

	guidePdfPath := pdfPath + ".guides.tmp.pdf"
	if err := createGuidePDF(guidePdfPath, width, height, n, binding); err != nil {
		return err
	}
	defer os.Remove(guidePdfPath)

	wm, err := api.PDFWatermark(guidePdfPath, "scale:1.0, pos:c, off:0 0, rot:0", true, false, types.POINTS)
	if err != nil {
		return err
	}

	conf := model.NewDefaultConfiguration()
	tempOut := pdfPath + ".tmp.pdf"
	err = api.AddWatermarksFile(pdfPath, tempOut, nil, wm, conf)
	if err != nil {
		return err
	}

	return os.Rename(tempOut, pdfPath)
}

func createGuidePDF(path string, w, h float64, n int, binding string) error {
	type Line struct {
		x1, y1, x2, y2 float64
		style          string // "solid" or "dashed"
	}
	var lines []Line

	addH := func(y float64, style string) {
		lines = append(lines, Line{0, y, w, y, style})
	}
	addV := func(x float64, style string) {
		lines = append(lines, Line{x, 0, x, h, style})
	}

	switch n {
	case 2:
		addH(h/2, "dashed")
	case 4:
		if binding == "long" {
			addH(h/2, "solid")
			addV(w/2, "dashed")
		} else {
			addH(h/2, "dashed")
			addV(w/2, "solid")
		}
	case 6:
		addH(h/3, "solid")
		addH(h*2/3, "solid")
		addV(w/2, "dashed")
	case 8:
		if binding == "long" {
			addH(h/2, "solid")
			addH(h/4, "dashed")
			addH(h*3/4, "dashed")
			addV(w/2, "solid")
		} else {
			addH(h/4, "solid")
			addH(h/2, "solid")
			addH(h*3/4, "solid")
			addV(w/2, "dashed")
		}
	}

	var buf bytes.Buffer
	offsets := make([]int, 5)
	buf.WriteString("%PDF-1.7\n")
	offsets[1] = buf.Len()
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	offsets[2] = buf.Len()
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	offsets[3] = buf.Len()
	fmt.Fprintf(&buf, "3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.2f %.2f] /Contents 4 0 R /Resources << /ProcSet [/PDF] >> >>\nendobj\n", w, h)

	var content strings.Builder
	content.WriteString("q\n0.5 G\n1 w\n")
	for _, l := range lines {
		if l.style == "dashed" {
			content.WriteString("[3] 0 d\n")
		} else {
			content.WriteString("[] 0 d\n")
		}
		fmt.Fprintf(&content, "%.2f %.2f m %.2f %.2f l S\n", l.x1, l.y1, l.x2, l.y2)
	}
	content.WriteString("Q\n")

	offsets[4] = buf.Len()
	fmt.Fprintf(&buf, "4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", content.Len(), content.String())

	xrefOffset := buf.Len()
	buf.WriteString("xref\n0 5\n0000000000 65535 f \n")
	for i := 1; i < 5; i++ {
		fmt.Fprintf(&buf, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&buf, "trailer\n<< /Size 5 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xrefOffset)

	return os.WriteFile(path, buf.Bytes(), 0644)
}
