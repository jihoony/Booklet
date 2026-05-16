# Booklet Creator

PDF 문서를 소책자(Booklet) 형태로 인쇄할 수 있도록 페이지를 재배치하고 변환해주는 CLI 도구입니다. [pdfcpu](https://github.com/pdfcpu/pdfcpu) 라이브러리를 기반으로 하며, 실제 제본 시 유용한 가이드라인 표시 기능과 대량 페이지 자동 분할(Multifolio) 기능을 제공합니다.

## 주요 기능

*   **소책자 변환**: PDF 페이지를 인쇄 후 접어서 책을 만들 수 있도록 순서를 재배치합니다.
*   **다양한 레이아웃**: 한 면당 2, 4, 6, 8 페이지 배치를 지원합니다.
*   **스마트 가이드라인**: 종이를 자르는 선(실선)과 접는 선(점선)을 구분하여 표시해줍니다. (텍스트 없이 깔끔한 선만 표시)
*   **자동 시그니처 모드**: 페이지 수가 많아 한 번에 접기 힘든 경우, 자동으로 여러 묶음(Signature)으로 나누어 처리합니다. (기본 10장 초과 시 자동 활성화)
*   **다양한 용지 지원**: A4, A5, Letter 등 다양한 용지 크기를 지원합니다.

## 설치 및 빌드

Go 언어(1.25 이상 권장)가 설치되어 있어야 합니다.

```bash
# 저장소 클론
git clone https://github.com/your-repo/booklet.git
cd booklet

# 의존성 다운로드
go mod tidy

# 실행
go run main.go [옵션]
```

또는 바이너리로 빌드하여 사용할 수 있습니다.

```bash
go build -o booklet main.go
./booklet [옵션]
```

## 사용법

기본적인 사용법은 다음과 같습니다.

```bash
go run main.go -i <입력파일.pdf> -o <출력파일.pdf> [옵션]
```

### 필수 옵션

*   `-i`, `-in`: 입력 PDF 파일 경로
*   `-o`, `-out`: 출력 PDF 파일 경로

### 선택 옵션

*   `-n`: 한 면에 배치할 페이지 수 (2, 4, 6, 8 지원, 기본값: 4)
    *   `2`: 2-up (한 장에 4페이지)
    *   `4`: 4-up (한 장에 8페이지) - *가장 일반적인 소책자 형태*
*   `-formsize`: 용지 크기 (A4, A5, Letter 등, 기본값: A4)
*   `-guides`: 접기/자르기 가이드라인 표시 (on/off, 기본값: on)
*   `-margin`: 여백 크기 (포인트 단위, 기본값: 10)
*   `-binding`: 제본 방향 (long/short, 기본값: long)
    *   `long`: 긴 쪽 제본 (일반적인 책)
    *   `short`: 짧은 쪽 제본 (달력 등)
*   `-btype`: 소책자 제본 유형 (booklet/advanced/perfectbound, 기본값: booklet)
    *   `booklet`: 일반적인 중첩 제본 (Saddle Stitch)
    *   `advanced`: 고급 제본 설정 적용
    *   `perfectbound`: 무선 제본 (떡제본)용 레이아웃
*   `-multifolio`: 시그니처 모드 강제 설정 (on/off, 기본값: off)
    *   *참고: 종이 장수가 10장을 넘어가면 자동으로 on으로 설정됩니다.*
*   `-foliosize`: 한 시그니처(묶음)당 시트 수 (multifolio=on일 때 사용, 기본값: 6)

## 사용 예시

**1. 기본 A4 소책자 만들기 (4-up)**
가장 일반적인 형태로, A4 용지 한 면에 4페이지(양면 8페이지)가 인쇄되어 반으로 접고 다시 반으로 접거나 잘라서 책을 만드는 형태입니다.
```bash
go run main.go -i input.pdf -o booklet.pdf
```

**2. A3 용지에 2페이지씩 배치 (2-up)**
A3 용지를 반으로 접어 A4 크기의 책을 만들 때 유용합니다.
```bash
go run main.go -i input.pdf -o booklet.pdf -n 2 -formsize A3
```

**3. 가이드라인 없이 생성**
```bash
go run main.go -i input.pdf -o booklet.pdf -guides off
```

**4. 페이지가 많은 문서 (자동 Multifolio)**
입력 파일이 매우 큰 경우(예: 200페이지), 프로그램이 자동으로 시그니처 모드를 활성화하여 6장(24~48페이지)씩 묶어서 출력합니다. 사용자는 출력된 묶음들을 각각 접어서 합치면 됩니다.

## 라이선스

이 프로젝트는 [pdfcpu](https://github.com/pdfcpu/pdfcpu)를 사용하며, 해당 라이브러리의 라이선스를 따릅니다.
