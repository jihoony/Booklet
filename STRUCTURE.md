# 프로젝트 구조 (Project Structure)

## 1. 디렉토리 구조 (Directory Tree)
```text
Booklet/
├── cmd/                # (향후 확장용) 별도 실행부
├── pkg/
│   └── booklet/        # 핵심 비즈니스 로직 (Core Engine)
│       ├── booklet.go  # 메인 프로세스 및 옵션 정의
│       ├── guide.go    # 가이드라인 생성 및 오버레이
│       └── logic.go    # 페이지 재배치 알고리즘
├── frontend/           # 데스크탑 앱 프런트엔드 (Vite + TS)
│   ├── index.html      # UI 메인 구조 및 스타일 정의
│   └── src/
│       └── main.ts     # UI 인터랙션 및 백엔드 연동
├── app.go              # 백엔드 브릿지 및 다이얼로그 기능
├── main.go             # 하이브리드 모드 엔트리 포인트
├── main_wails.go       # GUI 실행 설정 및 리소스 임베딩
├── wails.json          # Wails 프로젝트 설정 파일
└── go.mod              # Go 모듈 및 의존성 관리
```

## 2. 파일별 상세 역할 (File Mapping)

| 파일/폴더 | 주요 역할 | 비고 |
| :--- | :--- | :--- |
| `main.go` | 실행 인자를 확인하여 CLI 또는 GUI 모드로 분기 | 엔트리 포인트 |
| `app.go` | 파일 선택 창, 폴더 열기 등 OS 네이티브 기능 브릿지 | Frontend <-> Go |
| `pkg/booklet` | PDF를 소책자용으로 재배치하고 가이드라인을 삽입하는 핵심 엔진 | CLI/GUI 공통 사용 |
| `frontend/` | 사용자 설정을 수집하고 진행 상태를 보여주는 모던 UI | 웹 기술 기반 |
| `main_wails.go`| Wails 설정 로드 및 데스크탑 윈도우 생성 | GUI 전용 |

## 3. 데이터 흐름 (Data Flow)
1. **사용자**: GUI에서 PDF 파일을 선택하고 레이아웃 설정을 완료한 후 "생성" 버튼 클릭
2. **Frontend (`main.ts`)**: `window.go.main.App.ProcessBooklet(opts)` 호출
3. **Bridge (`app.go`)**: 호출을 받아 `pkg/booklet.Process(opts)` 실행
4. **Engine (`pkg/booklet`)**: PDF 수정 작업 완료 후 결과 파일 저장
5. **Bridge (`app.go`)**: 성공 응답 반환 및 사용자에게 폴더 열기 여부 확인
6. **사용자**: 폴더 열기 선택 시 OS 탐색기 구동
