# Module: internal/fridge

## 패키지 경로
`backend/internal/fridge`

## 책임
사용자 냉장고 재료 목록 CRUD.
유통기한 추적 및 ExpiryStatus 계산, 재료 현황 요약 제공.

---

## 파일 구조

```
internal/fridge/
├── model.go
├── repository.go
├── service.go
└── handler.go
```

---

## model.go

### 구조체

| 구조체                | 역할                                     |
|-----------------------|------------------------------------------|
| `Ingredient`          | MongoDB `ingredients` 컬렉션 문서        |
| `CreateIngredientReq` | POST /fridge 요청 DTO                    |
| `UpdateIngredientReq` | PATCH /fridge/:id 요청 DTO               |
| `IngredientResponse`  | 응답 DTO (ExpiryStatus 계산 필드 포함)   |
| `FridgeSummary`       | GET /fridge/summary 응답 DTO             |
| `BulkDeleteReq`       | DELETE /fridge body DTO                  |

### Ingredient 필드
```
ID          primitive.ObjectID  `bson:"_id"`
UserID      primitive.ObjectID  `bson:"user_id"`
Name        string              `bson:"name"`
Category    string              `bson:"category"`   // enum
Quantity    float64             `bson:"quantity"`
Unit        string              `bson:"unit"`
ExpiryDate  *time.Time          `bson:"expiry_date"` // nil 허용
Source      string              `bson:"source"`     // manual/receipt/vision
AddedAt     time.Time           `bson:"added_at"`
UpdatedAt   time.Time           `bson:"updated_at"`
```

### Category enum 상수
```
CategoryVegetable = "VEGETABLE"
CategoryFruit     = "FRUIT"
CategoryMeat      = "MEAT"
CategorySeafood   = "SEAFOOD"
CategoryDairy     = "DAIRY"
CategoryGrain     = "GRAIN"
CategoryCondiment = "CONDIMENT"
CategoryFrozen    = "FROZEN"
CategoryOther     = "OTHER"
```

### ExpiryStatus (계산, DB 미저장)
```
ExpiryStatusUrgent   = "URGENT"    // D-0 ~ D-2
ExpiryStatusSoon     = "SOON"      // D-3 ~ D-5
ExpiryStatusNormal   = "NORMAL"    // D-6 이후
ExpiryStatusNoExpiry = "NO_EXPIRY" // expiry_date == nil
```

---

## repository.go

### 인터페이스 (IngredientRepository)
```
FindAllByUserID(ctx, userID, filter) ([]*Ingredient, error)
FindByID(ctx, id, userID) (*Ingredient, error)
Create(ctx, ingredient) error
Update(ctx, id, userID, update) error
Delete(ctx, id, userID) error
BulkDelete(ctx, ids []ObjectID, userID) error
CountByExpiryStatus(ctx, userID) (*ExpiryStatusCount, error)
```

### MongoDB 컬렉션
- `ingredients`
  - index: `{user_id: 1, added_at: -1}`
  - index: `{user_id: 1, expiry_date: 1}`
  - index: `{user_id: 1, name: "text"}`

---

## service.go

### 인터페이스 (FridgeService)
```
ListIngredients(ctx, userID, filter) ([]*IngredientResponse, error)
GetIngredient(ctx, id, userID) (*IngredientResponse, error)
AddIngredient(ctx, userID, req) (*IngredientResponse, error)
UpdateIngredient(ctx, id, userID, req) (*IngredientResponse, error)
RemoveIngredient(ctx, id, userID) error
BulkRemove(ctx, ids, userID) error
GetSummary(ctx, userID) (*FridgeSummary, error)

// vision 도메인에서 호출하는 배치 등록
BulkAddIngredients(ctx, userID, items []CreateIngredientReq) ([]*IngredientResponse, error)
```

### 핵심 로직 포인트
- ExpiryStatus 계산: service 레이어에서 `time.Now()` 기준으로 계산 후 Response에 포함
- 동일 name 재료 중복 등록: service에서 감지하여 클라이언트에 409 Conflict 반환 (병합 여부는 클라이언트 결정)

---

## handler.go

### 엔드포인트 등록
```
GET    /api/v1/fridge              (JWT 필요)
POST   /api/v1/fridge              (JWT 필요)
GET    /api/v1/fridge/summary      (JWT 필요)
GET    /api/v1/fridge/:id          (JWT 필요)
PATCH  /api/v1/fridge/:id         (JWT 필요)
DELETE /api/v1/fridge/:id          (JWT 필요)
DELETE /api/v1/fridge              (JWT 필요, body: ids[])
```

---

## 외부 의존성

| 의존 대상       | 용도                |
|-----------------|---------------------|
| `internal/auth` | user_id Context 추출 (미들웨어를 통해) |
| `pkg/database`  | MongoDB 접근        |
