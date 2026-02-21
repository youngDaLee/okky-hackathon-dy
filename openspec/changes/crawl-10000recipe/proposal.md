## Why

추천 엔진이 동작하려면 내부 DB에 레시피 데이터가 필요하다.
현재 `recipes` 컬렉션이 비어 있어 추천 기능이 사실상 작동하지 않는 상태이며,
만개의레시피(10000recipe.com)에서 인기 레시피를 크롤링하여 시드 데이터를 확보한다.

## What Changes

- **Python 크롤러 신규 추가**: `crawler/` 디렉토리에 만개의레시피 크롤러 작성
  - 목록 페이지에서 레시피 URL 수집 (인기순 상위 100개)
  - 상세 페이지에서 Schema.org JSON-LD 파싱으로 레시피 데이터 추출
  - 재료명 정규화 (수량/단위 제거) 후 MongoDB에 upsert
  - `source_url` 기준 중복 제거로 주기적 동기화 지원
- **Go Recipe 모델 확장**: `raw_ingredients` 필드 추가 (원문 재료 저장용)
- **recommendation spec 갱신**: Recipe 엔티티에 `raw_ingredients` 필드 반영

## Non-goals

- 조리 단계(steps) 크롤링 — 원문 사이트 링크로 대체
- 레시피 이미지 자체 저장 — 원본 thumbnail URL 참조만
- 실시간 크롤링 API — 배치 스크립트로만 동작
- 전체 사이트 크롤링 — MVP는 인기순 100개로 한정

## Capabilities

### New Capabilities
- `recipe-crawler`: 만개의레시피 크롤링, 재료 정규화, MongoDB upsert 파이프라인

### Modified Capabilities
- `recommendation`: Recipe 엔티티에 `raw_ingredients []string` 필드 추가 (원문 재료 표시용)

## Impact

- **신규 디렉토리**: `crawler/` (Python, 기존 Go 백엔드와 독립)
- **모델 변경**: `backend/internal/recommendation/model.go` — Recipe struct에 필드 추가
- **DB**: `recipes` 컬렉션에 `source_url` unique index 추가 (upsert용)
- **의존성**: Python 패키지 (requests, beautifulsoup4, pymongo)
- **spec 변경**: `openspec/specs/recommendation/spec.md` — raw_ingredients 반영