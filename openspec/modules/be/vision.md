# Module: internal/vision

## 패키지 경로
`backend/internal/vision`

## 책임
이미지(영수증/냉장고 사진) 업로드 및 Vertex AI 기반 식재료 추출.
비동기 처리 후 사용자 확인을 거쳐 fridge 도메인에 재료 등록.

---

## 파일 구조

```
internal/vision/
├── model.go
├── repository.go
├── service.go
├── handler.go
└── worker.go
```

---

## model.go

### 구조체

| 구조체              | 역할                                         |
|---------------------|----------------------------------------------|
| `VisionJob`         | MongoDB `vision_jobs` 컬렉션 문서            |
| `ExtractedItem`     | 추출 결과 항목 (내장 도큐먼트)               |
| `CreateJobRequest`  | 이미지 업로드 요청 (multipart/form-data)     |
| `JobResponse`       | VisionJob 상태·결과 응답 DTO                 |
| `ConfirmRequest`    | 사용자 확인 요청 DTO (selected items 목록)   |

### VisionJob 필드
```
ID            primitive.ObjectID  `bson:"_id"`
UserID        primitive.ObjectID  `bson:"user_id"`
Type          string              `bson:"type"`          // RECEIPT / FRIDGE
ImageURL      string              `bson:"image_url"`     // GCS URL
Status        string              `bson:"status"`        // PENDING/PROCESSING/DONE/FAILED
RawResult     bson.M              `bson:"raw_result"`    // Vertex AI 원본
Extracted     []ExtractedItem     `bson:"extracted"`
ErrorMessage  *string             `bson:"error_message"`
CreatedAt     time.Time           `bson:"created_at"`
CompletedAt   *time.Time          `bson:"completed_at"`
```

### VisionJob Status enum
```
StatusPending    = "PENDING"
StatusProcessing = "PROCESSING"
StatusDone       = "DONE"
StatusFailed     = "FAILED"
```

### ExtractedItem 필드
```
Name            string   `bson:"name"`
NormalizedName  string   `bson:"normalized_name"`
Quantity        *float64 `bson:"quantity"`
Unit            *string  `bson:"unit"`
Confidence      float64  `bson:"confidence"`
Selected        bool     `bson:"selected"`
```

---

## repository.go

### 인터페이스 (VisionRepository)
```
Create(ctx, job) error
FindByID(ctx, id, userID) (*VisionJob, error)
UpdateStatus(ctx, id, status, result) error
DeleteExpired(ctx, before time.Time) error   // 24시간 이후 정리
```

### MongoDB 컬렉션
- `vision_jobs`
  - TTL index: `{created_at: 1}` with expireAfterSeconds=86400 (24시간)
  - index: `{user_id: 1, created_at: -1}`

---

## service.go

### 인터페이스 (VisionService)
```
CreateJob(ctx, userID, jobType, file multipart.File, filename string) (*JobResponse, error)
GetJob(ctx, id, userID) (*JobResponse, error)
ConfirmJob(ctx, id, userID, req ConfirmRequest) error
```

### 핵심 로직 포인트
- CreateJob:
  1. GCS 업로드 (`pkg/gcs`)
  2. VisionJob 생성 (PENDING)
  3. worker 채널에 job ID 전달 (비동기 트리거)
- ConfirmJob:
  1. Job 상태가 DONE인지 검증
  2. selected=true 항목만 `fridge.BulkAddIngredients` 호출

---

## handler.go

### 엔드포인트 등록
```
POST   /api/v1/vision/jobs              (JWT 필요, multipart/form-data)
GET    /api/v1/vision/jobs/:id          (JWT 필요)
POST   /api/v1/vision/jobs/:id/confirm  (JWT 필요)
```

---

## worker.go

### 역할
- 백그라운드 고루틴으로 실행
- `jobCh chan string` 채널로 job ID 수신
- Vertex AI 호출 → 결과 파싱 → VisionJob 업데이트

### 고루틴 설계
```
StartWorker(ctx, repo VisionRepository, visionClient *vertexai.Client)
  └── for job := range jobCh
        1. UpdateStatus(PROCESSING)
        2. vertexai.Analyze(imageURL, type)
        3. parseResult(rawResult) → []ExtractedItem
        4. UpdateStatus(DONE, extracted)
           또는 UpdateStatus(FAILED, errorMessage)
```

### 재료명 정규화
- normalizer 함수: RECEIPT 타입은 영수증 품목 텍스트 정제
- FRIDGE 타입은 Cloud Vision label → 한글 식품명 매핑 (내부 사전)

---

## 외부 의존성

| 의존 대상          | 용도                       |
|--------------------|----------------------------|
| `internal/fridge`  | ConfirmJob 시 재료 등록    |
| `pkg/gcs`          | 이미지 GCS 업로드          |
| `pkg/vertexai`     | Cloud Vision API 호출      |
| `pkg/database`     | MongoDB 접근               |
