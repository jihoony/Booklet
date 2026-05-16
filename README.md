# Booklet Pro (소책자 제조기)

PDF 파일을 소책자 인쇄용 레이아웃(2-up, 4-up 등)으로 자동 변환해주는 멀티플랫폼 데스크탑 애플리케이션 및 CLI 도구입니다.

![Booklet GUI Screenshot](https://raw.githubusercontent.com/jihoony/Booklet/main/screenshot.png) <!-- 실제 스크린샷 파일이 있다면 교체 가능 -->

## 🚀 주요 기능
- **GUI 기반 데스크탑 앱**: 드래그 앤 드롭으로 간편하게 PDF 변환 (Wails v2 기반)
- **강력한 CLI 지원**: 터미널 환경에서도 자동화 및 일괄 처리 가능
- **다양한 레이아웃**: 2-Up, 4-Up, 6-Up, 8-Up 지원
- **스마트 가이드라인**: 접지 및 절단을 돕는 시각적 안내선 삽입 기능
- **멀티폴리오(Multifolio) 지원**: 대량 페이지를 여러 권의 소책자로 분할 처리 가능

## 💻 사용 방법

### 1. 데스크탑 앱 (GUI) 모드
인자 없이 실행하면 현대적인 디자인의 GUI 모드로 실행됩니다.
```bash
./Booklet
```

### 2. CLI 모드
터미널에서 직접 옵션을 지정하여 실행할 수 있습니다.
```bash
./Booklet -i input.pdf -o output.pdf -n 4 -guides on
```

#### CLI 필수 인자:
- `-i`: 입력 PDF 파일 경로
- `-o`: 출력될 소책자 PDF 파일 경로

#### CLI 선택 인자:
- `-n`: 한 면에 배치할 페이지 수 (2, 4, 6, 8 / 기본값: 4)
- `-form`: 출력 용지 크기 (A4, A3 / 기본값: A4)
- `-guides`: 가이드라인 표시 여부 (on, off / 기본값: off)
- `-margin`: 페이지 간 여백 (0~50 / 기본값: 10)
- `-binding`: 제본 방향 (long, short / 기본값: long)
- `-multifolio`: 멀티폴리오 모드 사용 여부 (on, off / 기본값: off)

## 🛠 빌드 및 설치 가이드 (개발자용)

### 시스템 요구사항
- **Go**: 1.25 이상
- **Node.js**: 18.x 이상
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### 빌드 명령어 (Desktop App)
```bash
# 리눅스 환경 (Ubuntu 24.04+)
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
wails build -tags webkit2_41
```

## 📂 프로젝트 구조
- `pkg/booklet`: 소책자 변환 핵심 비즈니스 로직
- `frontend`: Vite + TypeScript 기반의 현대적 UI
- `app.go`: Go 백엔드와 프런트엔드 연결 브릿지
- `main.go`: 실행 환경(CLI/GUI) 감지 및 엔트리 포인트

## 📄 라이선스
MIT License
