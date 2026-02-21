## Context

추천 엔진(`recommendation`)이 동작하려면 `recipes` 컬렉션에 레시피가 있어야 한다.
현재 비어 있어 만개의레시피(10000recipe.com)에서 인기순 100개를 크롤링하여 시드 데이터를 확보한다.
크롤러는 Python으로 작성하며, 기존 Go 백엔드와 독립적으로 동작한다.

## Goals / Non-Goals

**Goals:**
- 만개의레시피에서 인기순 100개 레시피를 크롤링하여 MongoDB `recipes` 컬렉션에 저장
- 재료명을 정규화하여 추천 엔진 매칭에 사용 가능하도록 처리
- 원문 재료(수량 포함)도 함께 저장하여 프론트에서 표시 가능
- `source_url` 기준 upsert로 주기적 재실행 시 중복 방지

**Non-Goals:**
- 조리 단계(steps) 크롤링
- 레시피 이미지 다운로드/저장
- 크롤러를 API 서버로 노출
- 전체 사이트 크롤링

## Decisions

### D1. 크롤러 언어: Python

**선택**: Python (requests + BeautifulSoup4 + pymongo)
**대안**: Go (colly/goquery)
**이유**: 크롤링 생태계가 풍부하고, 한글 처리 및 프로토타이핑이 빠름.
배치 스크립트이므로 Go 서버와 통합할 필요 없음.

### D2. 데이터 추출 방식: Schema.org JSON-LD 우선

**선택**: HTML 내 `<script type="application/ld+json">` 에서 Schema.org Recipe 구조화 데이터 추출
**대안**: HTML 요소 직접 파싱
**이유**: JSON-LD는 구조화된 JSON이라 파싱이 안정적이고 사이트 UI 변경에 덜 취약함.
JSON-LD가 없는 경우 HTML fallback 처리.

### D3. 재료 정규화: 크롤링 시점에 수행

**선택**: 크롤러가 원문 재료에서 수량/단위를 제거하여 `required_ingredients`에 저장
**저장 구조**:
- `raw_ingredients`: 원문 그대로 (예: `"달걀 2개"`)
- `required_ingredients`: 정규화된 이름만 (예: `"달걀"`)
**정규화 로직**: 정규식으로 뒤쪽 숫자+단위 패턴 제거
  - 단위 사전: `g, kg, ml, L, 개, 대, 컵, 큰술, 작은술, 줄기, 장, 쪽, 톨, 봉, 근, 모, 포기, 마리` 등
  - 수량 표현: `약간, 적당량, 조금, 적당히, 한줌, 반개` 등

### D4. MainIngredient: 첫 번째 재료

**선택**: `recipeIngredient[0]`을 정규화하여 `main_ingredient`로 사용
**이유**: 만개의레시피는 주재료를 명시적 필드로 제공하지 않으며,
관행적으로 첫 번째 재료가 주재료인 경우가 많음.

### D5. 중복 제거: source_url unique index + upsert

**선택**: `source_url`에 unique index를 걸고, pymongo의 `update_one(upsert=True)` 사용
**이유**: 주기적 동기화 시 동일 레시피를 덮어쓰기(갱신)로 처리.
신규 레시피는 insert, 기존은 update.

### D6. 프로젝트 구조

```
crawler/
├── requirements.txt    # requests, beautifulsoup4, pymongo
├── config.py           # MONGO_URI, DB name, 크롤링 설정
├── scraper.py          # 목록 페이지 + 상세 페이지 크롤링
├── normalizer.py       # 재료명 정규화
├── store.py            # MongoDB upsert
└── main.py             # CLI 엔트리포인트
```

### D7. Rate Limiting

요청 간 1~2초 딜레이 (`time.sleep`). robots.txt에서 recipe 경로는 허용되어 있으나
서버 부담 최소화를 위해 적용.

## Risks / Trade-offs

- **[사이트 구조 변경]** → JSON-LD 파싱 실패 시 에러 로깅하고 해당 레시피 skip. HTML fallback은 MVP 이후.
- **[재료 정규화 정확도]** → 정규식 기반이라 100% 정확하지 않음. "양파 1/2개" 같은 분수 표현은 별도 처리 필요. 초기엔 best-effort로 운영.
- **[CookingTime 파싱]** → "5분 이내", "30분 이내", "60분 이내", "2시간 이상" 등 텍스트 → 분 단위 정수 변환 필요. 매핑 테이블로 처리.
- **[Category 매핑]** → 만개의레시피의 카테고리 체계와 현재 모델의 카테고리(한식/중식/일식/양식/간식/기타)가 다름. 목록 페이지에서 카테고리 정보가 제한적이므로 상세 페이지에서 추출하거나 기본값 "기타" 처리.

## Open Questions

- 주기적 동기화 주기와 실행 방법 (cron / 수동 / GitHub Actions) — 추후 결정
