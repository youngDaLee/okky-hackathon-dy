# Recommendation API Contract

> BE 작성 기준 계약서. 변경 시 FE와 협의 후 수정.

## 공통

**Headers:** 모든 API에 `Authorization: Bearer {access_token}` 필수

**공통 에러**

| Status | error | 상황 |
|--------|-------|------|
| 401 | `UNAUTHORIZED` | 토큰 없음 / 만료 |
| 404 | `RECIPE_NOT_FOUND` | 해당 id 레시피 없음 |

---

## GET /api/v1/recommendations

냉장고 재료 기반 레시피 추천. 사용자 allergens/dietary_prefs 자동 필터 적용.

**Query Parameters**

| 파라미터 | 타입 | 설명 | 기본값 |
|---------|------|------|--------|
| `tier` | string | `1` / `2` / `3` / `all` | `all` |
| `category` | string | 한식 / 중식 / 일식 / 양식 / 간식 / 기타 | - |
| `max_missing` | int | 부족 재료 최대 허용 개수 | - |
| `limit` | int | 결과 최대 개수 | `20` |

**Response 200**

```json
{
  "items": [
    {
      "tier": 1,
      "match_score": 1.0,
      "matched_ingredients": ["계란", "당근", "밥"],
      "missing_ingredients": [],
      "urgency_bonus": true,
      "recipe": {
        "id": "64f1a2b3c4d5e6f7a8b9c0e1",
        "title": "계란볶음밥",
        "description": "냉장고 파먹기 베스트 메뉴",
        "main_ingredient": "계란",
        "source_type": "INTERNAL",
        "source_url": null,
        "thumbnail_url": "https://storage.googleapis.com/fridge-master/recipes/egg-fried-rice.jpg",
        "category": "한식",
        "tags": ["초간단", "10분완성"],
        "cooking_time_min": 10,
        "difficulty": "EASY"
      }
    },
    {
      "tier": 2,
      "match_score": 0.75,
      "matched_ingredients": ["두부", "된장"],
      "missing_ingredients": ["애호박"],
      "urgency_bonus": false,
      "recipe": {
        "id": "64f1a2b3c4d5e6f7a8b9c0e2",
        "title": "두부 된장찌개",
        "description": "국민 찌개",
        "main_ingredient": "두부",
        "source_type": "EXTERNAL",
        "source_url": "https://www.youtube.com/watch?v=xyz789",
        "thumbnail_url": "https://img.youtube.com/vi/xyz789/hqdefault.jpg",
        "category": "한식",
        "tags": ["찌개"],
        "cooking_time_min": 20,
        "difficulty": "EASY"
      }
    }
  ],
  "total": 2,
  "fridge_ingredient_count": 15
}
```

**재료 0개일 때 Response 200**

```json
{
  "items": [],
  "total": 0,
  "fridge_ingredient_count": 0,
  "message": "냉장고에 재료를 추가하면 레시피를 추천해드려요!"
}
```

---

## GET /api/v1/recommendations/today

오늘의 추천. URGENT 재료(D-0~D-2) 기반 Tier 1 레시피 우선 노출.

**Response 200**

```json
{
  "urgent_ingredients": ["계란", "두부"],
  "items": [
    {
      "tier": 1,
      "match_score": 1.0,
      "matched_ingredients": ["계란"],
      "missing_ingredients": [],
      "urgency_bonus": true,
      "recipe": {
        "id": "64f1a2b3c4d5e6f7a8b9c0e1",
        "title": "계란찜",
        "description": "계란 빠르게 소비하기",
        "main_ingredient": "계란",
        "source_type": "INTERNAL",
        "source_url": null,
        "thumbnail_url": "https://storage.googleapis.com/fridge-master/recipes/steamed-egg.jpg",
        "category": "한식",
        "tags": ["계란요리", "다이어트"],
        "cooking_time_min": 15,
        "difficulty": "EASY"
      }
    }
  ],
  "total": 1
}
```

URGENT 재료가 없을 때 — 빈 `urgent_ingredients`와 일반 Tier 1 추천 반환:

```json
{
  "urgent_ingredients": [],
  "items": [...],
  "total": 5
}
```

---

## GET /api/v1/recipes

내부 레시피 검색.

**Query Parameters**

| 파라미터 | 타입 | 설명 | 예시 |
|---------|------|------|------|
| `q` | string | 레시피명 검색 | `김치` |
| `category` | string | 카테고리 필터 | `한식` |
| `tags` | string | 쉼표 구분 태그 필터 | `초간단,10분완성` |
| `difficulty` | string | `EASY` / `MEDIUM` / `HARD` | `EASY` |
| `limit` | int | 최대 개수 | `20` |
| `offset` | int | 페이지네이션 오프셋 | `0` |

**Response 200**

```json
{
  "items": [
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0e1",
      "title": "김치볶음밥",
      "main_ingredient": "김치",
      "source_type": "INTERNAL",
      "thumbnail_url": "https://storage.googleapis.com/fridge-master/recipes/kimchi-rice.jpg",
      "category": "한식",
      "tags": ["초간단", "볶음밥"],
      "cooking_time_min": 15,
      "difficulty": "EASY"
    }
  ],
  "total": 1,
  "offset": 0,
  "limit": 20
}
```

---

## GET /api/v1/recipes/:id

레시피 상세 조회(내부 레시피 전용).

**Response 200**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0e1",
  "title": "계란볶음밥",
  "description": "냉장고 파먹기 베스트 메뉴",
  "required_ingredients": ["계란", "밥", "간장"],
  "optional_ingredients": ["당근", "파"],
  "main_ingredient": "계란",
  "source_type": "INTERNAL",
  "source_url": null,
  "thumbnail_url": "https://storage.googleapis.com/fridge-master/recipes/egg-fried-rice.jpg",
  "category": "한식",
  "tags": ["초간단", "10분완성"],
  "cooking_time_min": 10,
  "difficulty": "EASY",
  "created_at": "2025-01-01T00:00:00Z"
}
```