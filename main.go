package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func main() {
	// 명령줄 플래그 정의
	var (
		input      string
		output     string
		n          int
		formSize   string
		guides     string
		margin     float64
		binding    string
		btype      string
		multifolio string
		folioSize  int
	)

	// 필수 인자
	flag.StringVar(&input, "in", "", "입력 PDF 파일 경로 (필수)")
	flag.StringVar(&input, "i", "", "입력 PDF 파일 경로 (단축, 필수)")

	flag.StringVar(&output, "out", "", "출력 booklet PDF 파일 경로 (필수)")
	flag.StringVar(&output, "o", "", "출력 booklet PDF 파일 경로 (단축, 필수)")

	flag.IntVar(&n, "n", 4, "한 면에 배치할 페이지 수 (2,4,6,8 지원, 기본값: 4)")

	// booklet 옵션들 (기본값 설정)
	flag.StringVar(&formSize, "formsize", "A4", "용지 크기 (A4, A5, Letter 등, 기본: A4)")
	flag.StringVar(&guides, "guides", "on", "접기/자르기 가이드라인 표시 (on/off, 기본: on)")
	flag.Float64Var(&margin, "margin", 10, "여백 크기 (포인트 단위, 기본: 10)")
	flag.StringVar(&binding, "binding", "long", "제본 방향 (long/short, 기본: long)")
	flag.StringVar(&btype, "btype", "booklet", "booklet 유형 (booklet/advanced/perfectbound, 기본: booklet)")

	// 긴 문서용 옵션
	flag.StringVar(&multifolio, "multifolio", "off", "시그니처 모드 (on/off, 기본: off)")
	flag.IntVar(&folioSize, "foliosize", 8, "한 시그니처당 시트 수 (multifolio=on일 때 사용, 기본: 8)")

	// 도움말
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "\nPDF를 소책자(booklet) 형태로 변환하는 도구\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "사용법:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  go run main.go -i input.pdf -o output.pdf [옵션들]\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "옵션:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "\n예시:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  go run main.go -i doc.pdf -o booklet.pdf -n 4 -formsize A4 -guides on\n")
		_, _ = fmt.Fprintf(os.Stderr, "  go run main.go -i long.pdf -o long_booklet.pdf -multifolio on -foliosize 16\n")
	}

	flag.Parse()

	// 필수 인자 검사
	if input == "" || output == "" {
		_, _ = fmt.Fprintln(os.Stderr, "오류: 입력 파일(-i)과 출력 파일(-o)은 필수입니다.")
		flag.Usage()
		os.Exit(1)
	}

	// n 값 검증
	if n != 2 && n != 4 && n != 6 && n != 8 {
		log.Fatalf("오류: -n 값은 2, 4, 6, 8 중 하나여야 합니다. (현재: %d)", n)
	}

	// guides 옵션 처리: 텍스트 없는 가이드라인을 위해 pdfcpu에는 off로 전달하고 후처리로 그린다.
	guidesOn := false
	if guides == "on" {
		guidesOn = true
		guides = "off"
	}

	// 페이지 수 확인 및 multifolio 자동 설정
	if multifolio == "off" {
		// 파일 존재 여부 확인 겸 페이지 수 체크
		ctx, err := api.ReadContextFile(input)
		if err != nil {
			log.Fatalf("입력 파일 분석 실패: %v", err)
		}
		pageCount := ctx.PageCount

		// 양면 인쇄 기준 한 장당 페이지 수 (n-up * 2)
		pagesPerSheet := n * 2

		// 필요한 총 종이 장수
		totalSheets := (pageCount + pagesPerSheet - 1) / pagesPerSheet

		// 종이 10장 이상이면 접기 힘드므로 multifolio 자동 활성화
		if totalSheets > 10 {
			multifolio = "on"
			fmt.Printf("알림: 전체 페이지(%d쪽)가 많아 시그니처 모드(multifolio)를 자동으로 활성화합니다.\n", pageCount)
			fmt.Printf("      총 %d장의 종이가 필요하며, %d장씩 묶어서(signature) 출력하도록 설정되었습니다.\n", totalSheets, folioSize)
			fmt.Printf("      (한 묶음당 %d 페이지)\n", folioSize*pagesPerSheet)
		}
	}

	// 옵션 문자열 조합
	descParts := []string{
		fmt.Sprintf("formsize:%s", formSize),
		fmt.Sprintf("guides:%s", guides),
		fmt.Sprintf("margin:%.1f", margin),
		fmt.Sprintf("binding:%s", binding),
		fmt.Sprintf("btype:%s", btype),
	}

	if multifolio == "on" {
		descParts = append(descParts, "multifolio:on", fmt.Sprintf("foliosize:%d", folioSize))
	}

	desc := ""
	if len(descParts) > 0 {
		for i, part := range descParts {
			if i > 0 {
				desc += ", "
			}
			desc += part
		}
	}

	// booklet 생성
	conf := model.NewDefaultConfiguration()
	fmt.Printf("입력: %s\n출력: %s\nn=%d\n옵션: %s\n\n변환 중...\n", input, output, n, desc)

	nupVal, err := api.PDFBookletConfig(n, desc, conf)
	if err != nil {
		log.Fatalf("설정 오류: %v", err)
	}

	err = api.BookletFile([]string{input}, output, nil, nupVal, conf)
	if err != nil {
		log.Fatalf("booklet 생성 실패: %v", err)
	}

	// 가이드라인 후처리
	if guidesOn {
		fmt.Println("가이드라인 추가 중...")
		if err := addGuides(output, n, binding); err != nil {
			log.Printf("경고: 가이드라인 추가 실패: %v", err)
		} else {
			fmt.Println("가이드라인 추가 완료")
		}
	}

	fmt.Printf("성공! 소책자 PDF가 생성되었습니다: %s\n", output)
}

func addGuides(pdfPath string, n int, binding string) error {
	// api.ReadContextFile을 사용하여 파일을 읽고 검증까지 수행
	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return fmt.Errorf("PDF 읽기 실패: %v", err)
	}

	if ctx.PageCount == 0 {
		return fmt.Errorf("페이지가 없습니다")
	}

	// 첫 번째 페이지의 MediaBox 정보 가져오기
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

	guidePdfPath := "temp_guides.pdf"
	if err := createGuidePDF(guidePdfPath, width, height, n, binding); err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(guidePdfPath)
	}()

	// 워터마크(스탬프) 설정
	wm, err := api.PDFWatermark(guidePdfPath, "scale:1.0, pos:c, off:0 0, rot:0", true, false, types.POINTS)
	if err != nil {
		return err
	}

	// 워터마크 추가를 위한 설정
	conf := model.NewDefaultConfiguration()

	// 덮어쓰기 문제 방지를 위해 임시 파일 사용
	tempOut := pdfPath + ".tmp.pdf"
	err = api.AddWatermarksFile(pdfPath, tempOut, nil, wm, conf)
	if err != nil {
		return err
	}

	// 임시 파일을 원본 파일로 이동 (덮어쓰기)
	return os.Rename(tempOut, pdfPath)
}

func createGuidePDF(path string, w, h float64, n int, binding string) error {
	type Line struct {
		x1, y1, x2, y2 float64
		style          string // "solid" or "dashed"
	}
	var lines []Line

	// 수평선 추가 (y 고정)
	addH := func(y float64, style string) {
		lines = append(lines, Line{0, y, w, y, style})
	}
	// 수직선 추가 (x 고정)
	addV := func(x float64, style string) {
		lines = append(lines, Line{x, 0, x, h, style})
	}

	// pdfcpu booklet 로직에 따른 가이드라인
	switch n {
	case 2:
		// 2-up: Horizontal Fold
		addH(h/2, "dashed")
	case 4:
		if binding == "long" {
			// Long Edge: Horizontal Cut, Vertical Fold
			addH(h/2, "solid")
			addV(w/2, "dashed")
		} else {
			// Short Edge: Horizontal Fold, Vertical Cut
			addH(h/2, "dashed")
			addV(w/2, "solid")
		}
	case 6:
		// 6-up: Horizontal Cut (1/3, 2/3), Vertical Fold
		addH(h/3, "solid")
		addH(h*2/3, "solid")
		addV(w/2, "dashed")
	case 8:
		if binding == "long" {
			// Long Edge: Horizontal Cut (1/2), Horizontal Fold (1/4, 3/4), Vertical Cut
			addH(h/2, "solid")
			addH(h/4, "dashed")
			addH(h*3/4, "dashed")
			addV(w/2, "solid")
		} else {
			// Short Edge: Horizontal Cut (1/4, 1/2, 3/4), Vertical Fold
			addH(h/4, "solid")
			addH(h/2, "solid")
			addH(h*3/4, "solid")
			addV(w/2, "dashed")
		}
	}

	var buf bytes.Buffer
	offsets := make([]int, 5)

	// Header
	buf.WriteString("%PDF-1.7\n")

	// Obj 1: Catalog
	offsets[1] = buf.Len()
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	// Obj 2: Pages
	offsets[2] = buf.Len()
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	// Obj 3: Page
	// Resources 딕셔너리 추가: /Resources << /ProcSet [/PDF] >>
	offsets[3] = buf.Len()
	_, _ = fmt.Fprintf(&buf, "3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.2f %.2f] /Contents 4 0 R /Resources << /ProcSet [/PDF] >> >>\nendobj\n", w, h)

	// Stream content construction
	var content strings.Builder
	content.WriteString("q\n0.5 G\n1 w\n")
	for _, l := range lines {
		if l.style == "dashed" {
			content.WriteString("[3] 0 d\n")
		} else {
			content.WriteString("[] 0 d\n")
		}
		_, _ = fmt.Fprintf(&content, "%.2f %.2f m %.2f %.2f l S\n", l.x1, l.y1, l.x2, l.y2)
	}
	content.WriteString("Q\n")

	// Obj 4: Content Stream
	offsets[4] = buf.Len()
	_, _ = fmt.Fprintf(&buf, "4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", content.Len(), content.String())

	// Xref
	xrefOffset := buf.Len()
	buf.WriteString("xref\n0 5\n0000000000 65535 f \n")
	for i := 1; i < 5; i++ {
		_, _ = fmt.Fprintf(&buf, "%010d 00000 n \n", offsets[i])
	}

	// Trailer
	_, _ = fmt.Fprintf(&buf, "trailer\n<< /Size 5 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xrefOffset)

	return os.WriteFile(path, buf.Bytes(), 0644)
}
