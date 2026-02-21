# Notification API Contract

> BE 작성 기준 계약서. 변경 시 FE와 협의 후 수정.

## 공통

**Headers:** 모든 API에 `Authorization: Bearer {access_token}` 필수

**공통 에러**

| Status | error | 상황 |
|--------|-------|------|
| 401 | `UNAUTHORIZED` | 토큰 없음 / 만료 |
| 404 | `NOTIFICATION_NOT_FOUND` | 해당 id 알림 없음 |

> MVP: In-App 알림 전용. 클라이언트 폴링 방식으로 동작.

---

## GET /api/v1/notifications

내 알림 목록 조회. 최근 50개, 생성 시각 내림차순.

**Response 200**

```json
{
  "items": [
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0f1",
      "ingredient_name": "계란",
      "type": "D_DAY",
      "message": "계란 오늘이 마지막이에요. 지금 요리하세요!",
      "is_read": false,
      "scheduled_date": "2025-03-05",
      "sent_at": "2025-03-05T09:00:00Z",
      "created_at": "2025-03-02T09:00:00Z"
    },
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0f2",
      "ingredient_name": "두부",
      "type": "D_3",
      "message": "두부이 3일 후 만료돼요. 오늘 사용해볼까요?",
      "is_read": true,
      "scheduled_date": "2025-03-03",
      "sent_at": "2025-03-03T09:00:00Z",
      "created_at": "2025-02-28T09:00:00Z"
    }
  ],
  "total": 2,
  "unread_count": 1
}
```

**`type` 값 설명**

| type | 의미 |
|------|------|
| `D_3` | 유통기한 3일 전 |
| `D_1` | 유통기한 1일 전 |
| `D_DAY` | 유통기한 당일 |

---

## PATCH /api/v1/notifications/:id/read

알림 단건 읽음 처리.

**Request Body 없음**

**Response 200**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0f1",
  "is_read": true
}
```

---

## POST /api/v1/notifications/read-all

내 알림 전체 읽음 처리.

**Request Body 없음**

**Response 200**

```json
{
  "updated_count": 5
}
```

---

## GET /api/v1/notifications/count

읽지 않은 알림 수. 헤더 뱃지 표시에 사용.

**Response 200**

```json
{
  "unread_count": 3
}
```