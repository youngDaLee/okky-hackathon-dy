## Why

Issue #4: 사용자의 냉장고 재료를 관리하는 CRUD API가 없어 MVP 핵심 기능 동작 불가.
스켈레톤으로 인터페이스만 정의된 fridge 도메인에 실제 비즈니스 로직과 MongoDB 연동을 구현한다.

## What Changes

- `internal/fridge/repository.go` — MongoDB CRUD 구현 (ingredients 컬렉션)
- `internal/fridge/service.go` — 비즈니스 로직 구현 (ExpiryStatus 계산, 200개 제한, 중복 처리)
- `internal/fridge/handler.go` — HTTP 핸들러 구현 (요청 파싱, 응답 직렬화, 에러 처리)
- `internal/server/router.go` — fridge 라우트 등록
- `pkg/config/config.go` — 환경변수 로드 구현

## Capabilities

### New Capabilities

(없음 — 기존 fridge capability에 구현 추가)

### Modified Capabilities

- `fridge`: 스켈레톤에서 실제 CRUD 구현으로 요구사항 변경 없음 (spec 그대로 유지, 구현만 추가)

## Impact

- `backend/internal/fridge/` 3개 파일 전면 구현
- MongoDB `ingredients` 컬렉션 인덱스 생성 포함
- Auth 미들웨어 의존 (#2 선행 조건 — JWT 미들웨어가 Context에 user_id를 주입한다고 가정하고 구현)

### Non-goals

- Auth 미들웨어 실제 구현 (별도 이슈 #2)
- Vision 도메인 연동 (별도 이슈 #6)
- Notification 트리거 (별도 이슈 #7)
- 테스트 코드 작성 (MVP 범위 외)
