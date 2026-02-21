# Vision API Contract

> BE 작성 기준 계약서. 변경 시 FE와 협의 후 수정.

## 공통

**Headers:** 모든 API에 `Authorization: Bearer {access_token}` 필수

**공통 에러**

| Status | error | 상황 |
|--------|-------|------|
| 401 | `UNAUTHORIZED` | 토큰 없음 / 만료 |
| 404 | `JOB_NOT_FOUND` | 해당 id 작업 없음 |

**VisionJob 상태 흐름**

```
PENDING → PROCESSING → DONE
                     → FAILED
```

---

## POST /api/v1/vision/jobs

이미지 업로드 및 AI 분석 작업 생성. `multipart/form-data`.

**Request (multipart/form-data)**

| 필드 | 타입 | 설명 |
|------|------|------|
| `image` | file | JPEG / PNG / WEBP, 최대 10MB |
| `type` | string | `RECEIPT` 또는 `FRIDGE` |

**Response 202**

작업이 생성되었고 비동기 처리 중. 결과는 폴링으로 확인.

```json
{
  "job_id": "64f1a2b3c4d5e6f7a8b9c0a1",
  "status": "PENDING",
  "type": "RECEIPT",
  "created_at": "2025-03-05T10:00:00Z"
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | type 없음 / 지원하지 않는 파일 형식 |
| 413 | `FILE_TOO_LARGE` | 10MB 초과 |

```json
// 413
{
  "error": "FILE_TOO_LARGE",
  "message": "이미지 파일은 10MB 이하여야 합니다."
}
```

---

## GET /api/v1/vision/jobs/:id

작업 상태 및 결과 조회. DONE 전까지 폴링.

**폴링 권장 간격:** 2초

**Response 200 — PENDING / PROCESSING**

```json
{
  "job_id": "64f1a2b3c4d5e6f7a8b9c0a1",
  "status": "PROCESSING",
  "type": "RECEIPT",
  "extracted": [],
  "created_at": "2025-03-05T10:00:00Z",
  "completed_at": null
}
```

**Response 200 — DONE**

```json
{
  "job_id": "64f1a2b3c4d5e6f7a8b9c0a1",
  "status": "DONE",
  "type": "RECEIPT",
  "extracted": [
    {
      "name": "계란10구",
      "normalized_name": "계란",
      "quantity": 10,
      "unit": "개",
      "confidence": 0.95,
      "selected": true
    },
    {
      "name": "당근 3개",
      "normalized_name": "당근",
      "quantity": 3,
      "unit": "개",
      "confidence": 0.88,
      "selected": true
    },
    {
      "name": "불명확품목",
      "normalized_name": "불명확품목",
      "quantity": null,
      "unit": null,
      "confidence": 0.32,
      "selected": false
    }
  ],
  "created_at": "2025-03-05T10:00:00Z",
  "completed_at": "2025-03-05T10:00:08Z"
}
```

> `confidence < 0.5`인 항목은 `selected: false` 기본값. FE에서 사용자가 수동 체크 가능.

**Response 200 — FAILED**

```json
{
  "job_id": "64f1a2b3c4d5e6f7a8b9c0a1",
  "status": "FAILED",
  "type": "FRIDGE",
  "extracted": [],
  "error_message": "Vertex AI 처리 타임아웃 (30초 초과)",
  "created_at": "2025-03-05T10:00:00Z",
  "completed_at": "2025-03-05T10:00:30Z"
}
```

---

## POST /api/v1/vision/jobs/:id/confirm

사용자가 결과를 확인하고 선택한 항목을 냉장고에 등록.
내부적으로 `POST /api/v1/fridge`를 호출.

**Request**

FE에서 사용자가 수정/선택한 최종 목록 전달:

```json
{
  "items": [
    {
      "name": "계란",
      "category": "DAIRY",
      "quantity": 10,
      "unit": "개",
      "expiry_date": "2025-03-15"
    },
    {
      "name": "당근",
      "category": "VEGETABLE",
      "quantity": 3,
      "unit": "개",
      "expiry_date": null
    }
  ]
}
```

**Response 200**

```json
{
  "registered_count": 2,
  "ingredients": [
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0d1",
      "name": "계란",
      "category": "DAIRY",
      "quantity": 10,
      "unit": "개",
      "expiry_date": "2025-03-15",
      "expiry_status": "NORMAL"
    },
    {
      "id": "64f1a2b3c4d5e6f7a8b9c0d2",
      "name": "당근",
      "category": "VEGETABLE",
      "quantity": 3,
      "unit": "개",
      "expiry_date": null,
      "expiry_status": "NO_EXPIRY"
    }
  ]
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | items 빈 배열, 필수 필드 누락 |
| 409 | `JOB_NOT_DONE` | 작업이 아직 DONE 상태가 아님 |

```json
// 409
{
  "error": "JOB_NOT_DONE",
  "message": "아직 분석이 완료되지 않았습니다. 잠시 후 다시 시도해주세요.",
  "current_status": "PROCESSING"
}
```