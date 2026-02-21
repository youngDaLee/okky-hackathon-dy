# Module: internal/notification

## 패키지 경로
`backend/internal/notification`

## 책임
유통기한 임박 재료 알림.
재료 등록/수정 시 D-3, D-1, D-day 알림 스케줄 생성.
매일 09:00 cron 배치 잡으로 당일 발송 처리.

---

## 파일 구조

```
internal/notification/
├── model.go
├── repository.go
├── service.go
├── handler.go
└── scheduler.go
```

---

## model.go

### 구조체

| 구조체                  | 역할                                       |
|-------------------------|--------------------------------------------|
| `ExpiryAlert`           | MongoDB `expiry_alerts` 컬렉션 문서        |
| `AlertResponse`         | 알림 응답 DTO                              |
| `AlertCountResponse`    | GET /notifications/count 응답 DTO          |

### ExpiryAlert 필드
```
ID              primitive.ObjectID  `bson:"_id"`
UserID          primitive.ObjectID  `bson:"user_id"`
IngredientID    primitive.ObjectID  `bson:"ingredient_id"`
IngredientName  string              `bson:"ingredient_name"`  // 스냅샷
Type            string              `bson:"type"`             // D_3/D_1/D_DAY
ScheduledDate   time.Time           `bson:"scheduled_date"`   // 발송 예정일 (자정 기준)
SentAt          *time.Time          `bson:"sent_at"`          // nil이면 미발송
IsRead          bool                `bson:"is_read"`
CreatedAt       time.Time           `bson:"created_at"`
```

### Alert Type enum
```
AlertTypeD3   = "D_3"
AlertTypeD1   = "D_1"
AlertTypeDay  = "D_DAY"
```

---

## repository.go

### 인터페이스 (NotificationRepository)
```
BulkCreate(ctx, alerts []*ExpiryAlert) error
FindByUserID(ctx, userID, limit int) ([]*ExpiryAlert, error)
FindUnsentByDate(ctx, date time.Time) ([]*ExpiryAlert, error)     // 배치 잡용
CountUnread(ctx, userID) (int64, error)
MarkSent(ctx, ids []ObjectID) error
MarkRead(ctx, id, userID) error
MarkAllRead(ctx, userID) error
DeleteByIngredientID(ctx, ingredientID) error                      // 재료 삭제 시
DeleteUnsentByIngredientID(ctx, ingredientID) error                // 재료 수정 시
```

### MongoDB 컬렉션
- `expiry_alerts`
  - unique index: `{ingredient_id: 1, type: 1}` (중복 방지)
  - index: `{user_id: 1, created_at: -1}`
  - index: `{scheduled_date: 1, sent_at: 1}` (배치 잡용)
  - TTL index: `{created_at: 1}` with expireAfterSeconds=2592000 (30일)

---

## service.go

### 인터페이스 (NotificationService)
```
// fridge 도메인에서 호출 (재료 등록/수정/삭제 시)
ScheduleAlerts(ctx, ingredient Ingredient) error
RescheduleAlerts(ctx, ingredient Ingredient) error  // 수정 시: 삭제 후 재생성
CancelAlerts(ctx, ingredientID ObjectID) error      // 삭제 시

// HTTP 핸들러에서 호출
ListAlerts(ctx, userID string, limit int) ([]*AlertResponse, error)
CountUnread(ctx, userID string) (int64, error)
ReadAlert(ctx, id, userID string) error
ReadAllAlerts(ctx, userID string) error

// scheduler에서 호출
ProcessDueAlerts(ctx) error
```

### ScheduleAlerts 로직
```
입력: ingredient (expiry_date 있는 경우만)

for type in [D_3, D_1, D_DAY]:
  scheduledDate = expiry_date - type.days
  if scheduledDate >= today:
    create ExpiryAlert (중복이면 skip, upsert)
```

---

## handler.go

### 엔드포인트 등록
```
GET    /api/v1/notifications           (JWT 필요)
GET    /api/v1/notifications/count     (JWT 필요)
PATCH  /api/v1/notifications/:id/read  (JWT 필요)
POST   /api/v1/notifications/read-all  (JWT 필요)
```

---

## scheduler.go

### 역할
애플리케이션 시작 시 백그라운드 고루틴으로 실행되는 cron 스케줄러.

### 스케줄 설계
```
StartScheduler(ctx, service NotificationService)
  └── cron.New(timezone=Asia/Seoul)
        └── "0 9 * * *"  →  ProcessDueAlerts(ctx)
```

### ProcessDueAlerts 처리 흐름
```
1. FindUnsentByDate(today) → 미발송 알림 목록
2. for each alert:
   - In-App: MarkSent(id)  // DB 업데이트만으로 발송 완료
   (MVP: 실제 Push 발송 없음)
3. 처리 건수 로그 출력
```

---

## 외부 의존성

| 의존 대상        | 용도                                      |
|------------------|-------------------------------------------|
| `internal/fridge` | 재료 등록/수정/삭제 이벤트에서 호출됨   |
| `pkg/database`   | MongoDB 접근                              |
| `github.com/robfig/cron/v3` | cron 스케줄러                 |
