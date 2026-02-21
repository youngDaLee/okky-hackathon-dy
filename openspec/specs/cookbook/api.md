# Cookbook API Contract

> BE 작성 기준 계약서. 변경 시 FE와 협의 후 수정.

## 공통

**Headers:** 모든 API에 `Authorization: Bearer {access_token}` 필수

**공통 에러**

| Status | error | 상황 |
|--------|-------|------|
| 401 | `UNAUTHORIZED` | 토큰 없음 / 만료 |
| 403 | `FORBIDDEN` | 다른 사용자의 저장 레시피 접근 시도 |
| 404 | `SAVED_RECIPE_NOT_FOUND` | 해당 id 저장 레시피 없음 |

---

## GET /api/v1/cookbook

저장된 레시피 목록 조회. 기본 정렬: 저장 시각 내림차순.

**Query Parameters**

| 파라미터 | 타입 | 설명 | 예시 |
|---------|------|------|------|
| `label` | string | 라벨 필터 | `주말요리` |
| `q` | string | 레시피 제목 prefix 검색 | `김치` |

**Response 200**

```json
{
  "items": [
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0d1",
      "recipe_snapshot": {
        "title": "계란볶음밥",
        "source_url": "https://www.youtube.com/watch?v=abc123",
        "thumbnail_url": "https://img.youtube.com/vi/abc123/hqdefault.jpg",
        "source_type": "EXTERNAL",
        "main_ingredient": "계란",
        "category": "한식"
      },
      "label": "자취 필수",
      "note": "간장 조금 덜 넣기",
      "rating": 4,
      "saved_at": "2025-03-01T09:00:00Z"
    },
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0d2",
      "recipe_snapshot": {
        "title": "두부 된장찌개",
        "source_url": null,
        "thumbnail_url": "https://storage.googleapis.com/fridge-master/recipes/abc.jpg",
        "source_type": "INTERNAL",
        "main_ingredient": "두부",
        "category": "한식"
      },
      "label": "미분류",
      "note": null,
      "rating": null,
      "saved_at": "2025-03-02T11:00:00Z"
    }
  ],
  "total": 2
}
```

---

## POST /api/v1/cookbook

레시피 저장.

**Request**

```json
{
  "recipe_id": "64f1a2b3c4d5e6f7a8b9c0e1",
  "label": "주말요리",
  "note": "마늘 두 배로"
}
```

`recipe_id`는 내부 레시피 저장 시 사용. 외부 레시피 저장 시 `recipe_snapshot` 직접 전달:

```json
{
  "recipe_snapshot": {
    "title": "계란볶음밥",
    "source_url": "https://www.youtube.com/watch?v=abc123",
    "thumbnail_url": "https://img.youtube.com/vi/abc123/hqdefault.jpg",
    "source_type": "EXTERNAL",
    "main_ingredient": "계란",
    "category": "한식"
  },
  "label": "미분류"
}
```

`label` 생략 시 `"미분류"` 기본값. `note`, `rating` 생략 가능.

**Response 201**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0d1",
  "recipe_snapshot": {
    "title": "계란볶음밥",
    "source_url": "https://www.youtube.com/watch?v=abc123",
    "thumbnail_url": "https://img.youtube.com/vi/abc123/hqdefault.jpg",
    "source_type": "EXTERNAL",
    "main_ingredient": "계란",
    "category": "한식"
  },
  "label": "미분류",
  "note": null,
  "rating": null,
  "saved_at": "2025-03-03T10:00:00Z"
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | recipe_id와 recipe_snapshot 모두 없음, label 20자 초과 |
| 409 | `ALREADY_SAVED` | 동일 레시피(내부 ID 또는 source_url) 이미 저장됨 |
| 429 | `SAVE_LIMIT_EXCEEDED` | 저장 레시피 500개 초과 |

```json
// 409
{
  "error": "ALREADY_SAVED",
  "message": "이미 저장된 레시피입니다.",
  "existing_id": "64f1a2b3c4d5e6f7a8b9c0d1"
}
```

---

## GET /api/v1/cookbook/:id

저장 레시피 단건 조회.

**Response 200**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0d1",
  "recipe_snapshot": {
    "title": "계란볶음밥",
    "source_url": "https://www.youtube.com/watch?v=abc123",
    "thumbnail_url": "https://img.youtube.com/vi/abc123/hqdefault.jpg",
    "source_type": "EXTERNAL",
    "main_ingredient": "계란",
    "category": "한식"
  },
  "label": "자취 필수",
  "note": "간장 조금 덜 넣기",
  "rating": 4,
  "saved_at": "2025-03-01T09:00:00Z",
  "updated_at": "2025-03-02T14:00:00Z"
}
```

---

## PATCH /api/v1/cookbook/:id

라벨 / 메모 / 평점 수정. 변경할 필드만 포함.

**Request**

```json
{
  "label": "다이어트",
  "note": "저염으로 조절",
  "rating": 5
}
```

평점 초기화 시 `"rating": null` 전달.

**Response 200**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0d1",
  "label": "다이어트",
  "note": "저염으로 조절",
  "rating": 5,
  "updated_at": "2025-03-05T10:00:00Z"
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | rating 1~5 범위 벗어남, label 20자 초과, note 500자 초과 |

---

## DELETE /api/v1/cookbook/:id

저장 레시피 단건 삭제.

**Response 204**

본문 없음.

---

## GET /api/v1/cookbook/labels

사용자의 라벨 목록 및 각 라벨의 저장 수 반환.

**Response 200**

```json
{
  "labels": [
    { "label": "미분류", "count": 10 },
    { "label": "주말요리", "count": 5 },
    { "label": "다이어트", "count": 3 }
  ]
}
```