# Fridge API Contract

> BE 작성 기준 계약서. 변경 시 FE와 협의 후 수정.

## 공통

**Headers:** 모든 API에 `Authorization: Bearer {access_token}` 필수

**공통 에러**

| Status | error | 상황 |
|--------|-------|------|
| 401 | `UNAUTHORIZED` | 토큰 없음 / 만료 |
| 403 | `FORBIDDEN` | 다른 사용자의 재료 접근 시도 |
| 404 | `INGREDIENT_NOT_FOUND` | 해당 id 재료 없음 |

---

## GET /api/v1/fridge

내 재료 목록 조회. 기본 정렬: 유통기한 오름차순.

**Query Parameters**

| 파라미터 | 타입 | 설명 | 예시 |
|---------|------|------|------|
| `category` | string | 카테고리 필터 | `VEGETABLE` |
| `expiry_status` | string | 유통기한 상태 필터 | `URGENT` / `SOON` / `NORMAL` / `NO_EXPIRY` |
| `q` | string | 재료명 prefix 검색 | `당근` |

**Response 200**

```json
{
  "items": [
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0d1",
      "name": "계란",
      "category": "DAIRY",
      "quantity": 8,
      "unit": "개",
      "expiry_date": "2025-03-05",
      "expiry_status": "URGENT",
      "source": "manual",
      "added_at": "2025-03-01T09:00:00Z",
      "updated_at": "2025-03-01T09:00:00Z"
    },
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0d2",
      "name": "올리브오일",
      "category": "CONDIMENT",
      "quantity": 1,
      "unit": "병",
      "expiry_date": null,
      "expiry_status": "NO_EXPIRY",
      "source": "manual",
      "added_at": "2025-02-20T10:00:00Z",
      "updated_at": "2025-02-20T10:00:00Z"
    }
  ],
  "total": 2
}
```

---

## POST /api/v1/fridge

재료 추가.

**Request**

```json
{
  "name": "당근",
  "category": "VEGETABLE",
  "quantity": 3,
  "unit": "개",
  "expiry_date": "2025-03-10"
}
```

`expiry_date`는 생략 가능 (생략 시 `null` 저장).

**Response 201**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0d3",
  "name": "당근",
  "category": "VEGETABLE",
  "quantity": 3,
  "unit": "개",
  "expiry_date": "2025-03-10",
  "expiry_status": "NORMAL",
  "source": "manual",
  "added_at": "2025-03-01T11:00:00Z",
  "updated_at": "2025-03-01T11:00:00Z"
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | 필수 필드 누락, 허용되지 않는 category, 수량 0 이하 |
| 409 | `DUPLICATE_INGREDIENT` | 동일 이름 재료 이미 존재 (합산 여부는 FE에서 선택 후 재요청) |

```json
// 409
{
  "error": "DUPLICATE_INGREDIENT",
  "message": "이미 등록된 재료입니다. 기존 수량에 합산하려면 PATCH를 사용하세요.",
  "existing_id": "64f1a2b3c4d5e6f7a8b9c0d1"
}
```

---

## GET /api/v1/fridge/:id

재료 단건 조회.

**Response 200**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0d1",
  "name": "계란",
  "category": "DAIRY",
  "quantity": 8,
  "unit": "개",
  "expiry_date": "2025-03-05",
  "expiry_status": "URGENT",
  "source": "manual",
  "added_at": "2025-03-01T09:00:00Z",
  "updated_at": "2025-03-01T09:00:00Z"
}
```

---

## PATCH /api/v1/fridge/:id

재료 수정. 변경할 필드만 포함.

**Request**

```json
{
  "quantity": 5,
  "expiry_date": "2025-03-08"
}
```

**Response 200**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0d1",
  "name": "계란",
  "category": "DAIRY",
  "quantity": 5,
  "unit": "개",
  "expiry_date": "2025-03-08",
  "expiry_status": "SOON",
  "source": "manual",
  "added_at": "2025-03-01T09:00:00Z",
  "updated_at": "2025-03-02T15:00:00Z"
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | 수량 0 이하, 허용되지 않는 category |

---

## DELETE /api/v1/fridge/:id

재료 단건 삭제.

**Response 204**

본문 없음.

---

## DELETE /api/v1/fridge

재료 다수 일괄 삭제. (요리 완료 후 사용한 재료 일괄 제거에 사용)

**Request**

```json
{
  "ids": [
    "64f1a2b3c4d5e6f7a8b9c0d1",
    "64f1a2b3c4d5e6f7a8b9c0d2"
  ]
}
```

**Response 200**

```json
{
  "deleted_count": 2
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | ids 빈 배열 |

> 존재하지 않는 id가 포함되어도 있는 것만 삭제하고 200 반환 (silent ignore).

---

## GET /api/v1/fridge/summary

냉장고 재료 현황 요약. 알림 도메인 및 홈 대시보드에서 사용.

**Response 200**

```json
{
  "total": 15,
  "urgent": 2,
  "soon": 3,
  "normal": 8,
  "no_expiry": 2
}
```