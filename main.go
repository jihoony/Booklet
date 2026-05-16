package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"Booklet/pkg/booklet"
)

func main() {
	if len(os.Args) > 1 {
		runCLI()
	} else {
		runGUI()
	}
}

func runCLI() {
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
	flag.StringVar(&formSize, "formsize", "A4", "용지 크기 (A4, A5, Letter 등, 기본: A4)")
	flag.StringVar(&guides, "guides", "on", "접기/자르기 가이드라인 표시 (on/off, 기본: on)")
	flag.Float64Var(&margin, "margin", 10, "여백 크기 (포인트 단위, 기본: 10)")
	flag.StringVar(&binding, "binding", "long", "제본 방향 (long/short, 기본: long)")
	flag.StringVar(&btype, "btype", "booklet", "booklet 유형 (booklet/advanced/perfectbound, 기본: booklet)")
	flag.StringVar(&multifolio, "multifolio", "off", "시그니처 모드 (on/off, 기본: off)")
	flag.IntVar(&folioSize, "foliosize", 6, "한 시그니처당 시트 수 (multifolio=on일 때 사용, 기본: 6)")

	// 도움말
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "\nPDF를 소책자(booklet) 형태로 변환하는 도구\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "사용법:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  go run main.go -i input.pdf -o output.pdf [옵션들]\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "옵션:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "\n예시:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  go run main.go -i doc.pdf -o booklet.pdf -n 4 -formsize A4 -guides on\n")
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

	// 옵션 구조체 생성
	opts := booklet.Options{
		Input:      input,
		Output:     output,
		N:          n,
		FormSize:   formSize,
		Guides:     guides == "on",
		Margin:     margin,
		Binding:    binding,
		BType:      btype,
		Multifolio: multifolio == "on",
		FolioSize:  folioSize,
	}

	fmt.Printf("입력: %s\n출력: %s\nn=%d\n\n변환 중...\n", opts.Input, opts.Output, opts.N)

	// 프로세스 실행
	if err := booklet.Process(opts); err != nil {
		log.Fatalf("오류 발생: %v", err)
	}

	fmt.Printf("성공! 소책자 PDF가 생성되었습니다: %s\n", output)
}
