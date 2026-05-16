# 릴리즈 및 배포 가이드 (Release Guide)

본 문서는 `Booklet Pro`를 각 운영체제별 릴리즈 버전으로 빌드하고 배포하는 방법을 설명합니다.

## 1. 사전 준비
빌드 전 프런트엔드 의존성이 최신 상태인지 확인합니다.
```bash
cd frontend
npm install
cd ..
```

## 2. 플랫폼별 빌드 명령어

### Windows (64bit)
윈도우 사용자를 위한 단일 실행 파일(`.exe`)을 생성합니다.
```bash
wails build -platform windows/amd64 -clean
```
- `-clean`: 이전 빌드 캐시를 삭제하고 새로 빌드합니다.
- 결과물: `build/bin/Booklet.exe`

### Linux (Ubuntu/Mint 등)
현재 개발 환경과 동일한 라이브러리 버전을 사용하는 단일 바이너리를 생성합니다.
```bash
wails build -platform linux/amd64 -clean -tags webkit2_41
```
- 결과물: `build/bin/Booklet`

### macOS (Universal)
Intel과 Apple Silicon(M1/M2) 모두 지원하는 앱 번들을 생성합니다.
```bash
wails build -platform darwin/universal -clean
```
- 결과물: `build/bin/Booklet.app`

## 4. 보안 및 난독화 빌드 (Obfuscation)
프로그램의 소스 코드나 임베딩된 에셋을 보호하기 위해 난독화 옵션을 사용할 수 있습니다.

### 기본 난독화
에셋 암호화 및 기본적인 바이너리 보호를 적용합니다.
```bash
wails build -clean -obfuscated
```

### 강력한 난독화 (Garble 사용)
Go 소스 코드의 함수명, 변수명 등을 알아보기 어렵게 만듭니다.
```bash
# garble 설치 필요: go install mvdan.cc/garble@latest
wails build -clean -obfuscated -garbleargs "-literals -tiny -seed=random"
```

## 5. 배포 시 주의사항
- **의존성**: Wails 앱은 OS의 네이티브 웹뷰 엔진을 사용합니다.
  - Windows: WebView2 (Win 10/11 기본 내장)
  - Linux: libwebkit2gtk-4.0 또는 4.1
  - macOS: WKWebView (기본 내장)
- **아이콘 변경**: 배포용 아이콘을 적용하려면 `build/appicon.png` 파일을 교체한 후 다시 빌드하세요.
- **압축**: 배포 시에는 `build/bin/` 폴더에 생성된 실행 파일을 ZIP 등으로 압축하여 전달하는 것이 좋습니다.

## 4. 버전 관리
`wails.json` 파일의 `info` 섹션에서 버전을 수정할 수 있습니다.
```json
"info": {
  "productVersion": "1.0.0",
  "copyright": "Copyright © 2026"
}
```
