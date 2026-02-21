## Why

Issue #6: 냉털마스터의 핵심 가치인 레시피 추천 기능이 없어 MVP 동작 불가.
스켈레톤으로 인터페이스만 정의된 recommendation 도메인에 3-Tier 매칭 엔진과 내부 DB + 외부 검색 연동을 구현한다.

## What Changes

- `internal/recommendation/repository.go` — MongoDB `recipes` 컬렉션 CRUD + 인덱스 생성
- `internal/recommendation/service.go` — 3-Tier 매칭 로직, urgency_bonus, 외부 보완 오케스트레이션
- `internal/recommendation/handler.go` — 4개 엔드포인트 HTTP 핸들러 구현
- `internal/recommendation/external.go` — YouTube / Google Custom Search 연동 (메모리 캐시 포함)
- `internal/server/router.go` — recommendation 라우트 등록
- `internal/fridge` — 의존: userID 기준 재료 목록 + URGENT 재료 조회 (기존 FridgeService 재사용)

## Capabilities

### New Capabilities

(없음 — 기존 recommendation capability에 구현 추가)

### Modified Capabilities

- `recommendation`: 스켈레톤 → 3-Tier 매칭 엔진 + 외부 검색 연동 완전 구현

## Impact

- `backend/internal/recommendation/` 4개 파일 전면 구현
- MongoDB `recipes` 컬렉션 인덱스 생성 포함
- recommendation 서비스가 fridge 서비스에 의존 (userID로 재료 목록 조회)
- 외부 API 키 없을 경우 graceful fallback (외부 검색 생략, 내부 결과만 반환)

### Non-goals

- 레시피 DB 시딩 (별도 데이터 작업)
- 식이 제한(allergens) 필터링 — auth 도메인 User 모델 미구현으로 MVP 제외
- Cookbook 도메인 연동 (별도 이슈 #7)
- 레시피 CRUD 관리자 API
- 테스트 코드 작성 (MVP 범위 외)
