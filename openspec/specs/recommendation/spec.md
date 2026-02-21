# Recommendation Spec

## 개요

레시피 추천 엔진 도메인.
사용자의 냉장고 재료를 기반으로 3-Tier 매칭 알고리즘을 적용하고,
내부 DB 레시피 + 외부 YouTube/블로그 링크를 혼합하여 추천 리스트를 제공한다.

---

## 레시피 데이터 소스

### 내부 DB (Source: INTERNAL)
- 사전에 큐레이션/크롤링된 레시피 메타데이터
- 재료 목록이 정형화된 데이터 → Tier 분류 정확도 높음
- 직접 레시피 내용 제공 가능

### 외부 링크 (Source: EXTERNAL)
- YouTube 동영상 URL 또는 블로그 포스트 URL
- 재료 목록은 AI 또는 스크래핑으로 추출 (best-effort)
- 신뢰도 낮은 재료 목록 → Tier 분류 보조용

---

## 핵심 엔티티

### Recipe (내부 DB)

| 필드                  | 타입       | 설명                                          |
|-----------------------|------------|-----------------------------------------------|
| id                    | ObjectID   | 레시피 고유 ID                                |
| title                 | string     | 레시피명                                      |
| description           | string     | 간단 설명                                     |
| required_ingredients  | []string   | 필수 재료명 목록 (정규화된 이름)              |
| optional_ingredients  | []string   | 선택 재료명 목록                              |
| main_ingredient       | string     | 대표 주재료 1개                               |
| source_type           | string     | INTERNAL / EXTERNAL                           |
| source_url            | string     | 원본 링크 (YouTube, 블로그 등)                |
| thumbnail_url         | string     | 썸네일 이미지 URL                             |
| category              | string     | 한식 / 중식 / 일식 / 양식 / 간식 / 기타      |
| tags                  | []string   | 검색 태그 (예: 초간단, 10분완성, 다이어트)    |
| cooking_time_min      | int        | 조리 시간 (분)                                |
| difficulty            | string     | EASY / MEDIUM / HARD                          |
| created_at            | time.Time  |                                               |

### RecommendationResult (응답용, DB 저장 안 함)

| 필드                  | 타입     | 설명                                          |
|-----------------------|----------|-----------------------------------------------|
| recipe                | Recipe   | 레시피 정보                                   |
| tier                  | int      | 1 / 2 / 3                                     |
| match_score           | float64  | 보유 재료 일치율 0.0 ~ 1.0                    |
| matched_ingredients   | []string | 일치한 재료 목록                              |
| missing_ingredients   | []string | 부족한 재료 목록                              |
| urgency_bonus         | bool     | 유통기한 임박 재료 포함 여부 (상단 부스트)    |

---

## 3-Tier 매칭 알고리즘

```
입력: user_ingredients = 사용자 냉장고 재료 목록

Tier 1 (100% 매칭)
  조건: required_ingredients ⊆ user_ingredients
  → 추가 구매 없이 완성 가능
  정렬: urgency_bonus DESC, match_score DESC

Tier 2 (주재료 매칭)
  조건: main_ingredient ∈ user_ingredients
        AND match_score >= 0.6 (required 중 60% 이상 보유)
  → 소량 추가만 필요
  정렬: match_score DESC, urgency_bonus DESC

Tier 3 (부분 매칭)
  조건: match_score >= 0.3 (required 중 30% 이상 보유)
  → 참고용, 재료 구매 계획 수립
  정렬: match_score DESC
```

### 유통기한 긴급도 보정 (urgency_bonus)

```
recipe.required_ingredients 중 하나라도
user의 URGENT 재료(D-0 ~ D-2)와 겹치면:
  → urgency_bonus = true
  → Tier 1, Tier 2 내에서 최상단 노출
```

### 식이 제한 필터링
- User.allergens 또는 User.dietary_prefs에 해당하는 재료를 포함한 레시피 제외
- 필터 적용 후 Tier 분류 수행

---

## 외부 콘텐츠 통합 (A+B 혼합)

### 내부 DB 우선
1. 내부 DB에서 Tier 1/2/3 레시피 조회
2. 내부 결과가 충분하지 않으면 외부 검색 보완

### 외부 검색 보완 기준
- Tier 1 내부 결과 < 3개: 외부 검색으로 보완
- Tier 2 내부 결과 < 5개: 외부 검색으로 보완

### 외부 검색 전략
- 쿼리: `"{main_ingredient} 레시피"` YouTube Data API v3, Google Custom Search
- 결과에서 재료 목록 추출 (설명 텍스트 파싱, best-effort)
- 추출 실패 시 Tier 미분류 "참고 레시피"로 별도 노출

---

## API 엔드포인트

```
GET    /api/v1/recommendations          추천 레시피 조회
GET    /api/v1/recommendations/today    오늘의 추천 (URGENT 재료 기반 Tier 1)
GET    /api/v1/recipes                  내부 레시피 검색 (이름/카테고리/태그)
GET    /api/v1/recipes/:id              레시피 상세
```

### GET /api/v1/recommendations 쿼리 파라미터

| 파라미터    | 설명                              | 기본값 |
|-------------|-----------------------------------|--------|
| tier        | 1 / 2 / 3 / all                   | all    |
| category    | 한식 / 양식 / 기타                | -      |
| max_missing | 부족 재료 최대 허용 개수          | -      |
| limit       | 결과 최대 개수                    | 20     |

---

## 제약 조건

- 추천은 항상 로그인 사용자 기준 (재료 목록 필요)
- 재료 목록이 0개인 경우 → 빈 결과 또는 인기 레시피 fallback (구현 시 결정)
- 외부 검색 결과 캐싱: 동일 재료 조합은 1시간 캐시
- 내부 DB 레시피 수: MVP 기준 최소 500개 이상 사전 등록 권장

---

## 의존성

- **auth**: user_id, dietary_prefs, allergens
- **fridge**: 사용자 냉장고 재료 목록 + ExpiryStatus
- YouTube Data API v3 또는 Google Custom Search API (외부 콘텐츠)
