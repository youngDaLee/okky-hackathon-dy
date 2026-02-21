## 1. Go 모델 확장 (백엔드)

- [x] 1.1 `recommendation/model.go`의 Recipe struct에 `RawIngredients []string` 필드 추가
- [x] 1.2 `recommendation/model.go`의 RecipeSnap 또는 관련 응답 struct에 `raw_ingredients` JSON 태그 반영

## 2. Python 크롤러 프로젝트 셋업

- [x] 2.1 `crawler/` 디렉토리 생성 및 `requirements.txt` 작성 (requests, beautifulsoup4, pymongo)
- [x] 2.2 `crawler/config.py` 작성 — MONGO_URI, DB_NAME, COLLECTION, 크롤링 설정(페이지 수, 딜레이)

## 3. 크롤링 핵심 로직 (Python)

- [x] 3.1 `crawler/scraper.py` — 목록 페이지 크롤링: `order=reco`로 page 1~4 순회하여 레시피 URL 100개 수집
- [x] 3.2 `crawler/scraper.py` — 상세 페이지 크롤링: JSON-LD `@type: Recipe` 파싱으로 데이터 추출
- [x] 3.3 `crawler/scraper.py` — 조리시간/난이도 텍스트 파싱 (분 단위 변환, difficulty 매핑)

## 4. 재료 정규화 (Python)

- [x] 4.1 `crawler/normalizer.py` — 수량/단위 제거 정규식 작성 (g, 개, 큰술, 약간, 1/2 등)
- [x] 4.2 `crawler/normalizer.py` — 괄호 부가설명 제거 및 공백 trim

## 5. MongoDB 저장 (Python)

- [x] 5.1 `crawler/store.py` — MongoDB 연결 및 `source_url` unique index 생성
- [x] 5.2 `crawler/store.py` — `source_url` 기준 upsert 함수 구현

## 6. 엔트리포인트 및 통합

- [x] 6.1 `crawler/main.py` — CLI 엔트리포인트: 수집 → 파싱 → 정규화 → 저장 파이프라인 연결
- [x] 6.2 실행 테스트: `python crawler/main.py` 로 100개 크롤링 → MongoDB 확인
