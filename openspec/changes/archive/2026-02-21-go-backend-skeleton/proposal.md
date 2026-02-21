## Why

냉털마스터 Go 백엔드의 첫 번째 단계로, 기능 구현 전 프로젝트 뼈대를 세운다.
디렉토리 구조, 패키지 레이아웃, 빈 인터페이스/구조체 파일을 확립하여 팀 병렬 개발의 기반을 마련한다.

## What Changes

- `backend/` 디렉토리 신규 생성 (Go 모듈 루트)
- `cmd/server/main.go` — 진입점 (빈 main 함수)
- `internal/server/` — Gin 라우터 및 미들웨어 스켈레톤
- `internal/auth/` — model, repository, service, handler (빈 시그니처)
- `internal/fridge/` — model, repository, service, handler (빈 시그니처)
- `internal/vision/` — model, repository, service, handler, worker (빈 시그니처)
- `internal/recommendation/` — model, repository, service, handler, external (빈 시그니처)
- `internal/cookbook/` — model, repository, service, handler (빈 시그니처)
- `internal/notification/` — model, repository, service, handler, scheduler (빈 시그니처)
- `pkg/config/`, `pkg/database/`, `pkg/jwt/`, `pkg/gcs/`, `pkg/vertexai/` — 공유 유틸 스켈레톤
- `go.mod` / `go.sum` — 의존성 선언

## Capabilities

### New Capabilities

- `be-project-skeleton`: Go 백엔드 프로젝트 전체 디렉토리·파일 뼈대 (package 선언 + 빈 타입/인터페이스/함수만 포함, 기능 구현 없음)

### Modified Capabilities

(없음 — 기존 스펙 요구사항 변경 없음)

## Impact

- `backend/` 신규 디렉토리 전체 생성
- 기존 코드 영향 없음
- 이후 각 도메인 기능 구현 Change의 기반이 됨

### Non-goals

- 실제 비즈니스 로직 구현
- DB 연결 및 쿼리 작성
- 인증/JWT 동작 구현
- 외부 API (Vertex AI, YouTube 등) 실제 연동
- 테스트 코드 작성
