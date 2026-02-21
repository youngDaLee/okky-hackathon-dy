# Module: internal/cookbook

## 패키지 경로
`backend/internal/cookbook`

## 책임
"나만의 요리책" — 사용자가 마음에 드는 레시피를 저장하고,
라벨로 분류하며 만족도(평점)를 기록하는 개인 보관함.

---

## 파일 구조

```
internal/cookbook/
├── model.go
├── repository.go
├── service.go
└── handler.go
```

---

## model.go

### 구조체

| 구조체              | 역할                                         |
|---------------------|----------------------------------------------|
| `SavedRecipe`       | MongoDB `saved_recipes` 컬렉션 문서          |
| `RecipeSnap`        | 저장 시점 레시피 스냅샷 (내장 도큐먼트)      |
| `SaveRecipeRequest` | POST /cookbook 요청 DTO                      |
| `UpdateRecipeRequest` | PATCH /cookbook/:id 요청 DTO               |
| `SavedRecipeResponse` | 응답 DTO                                   |
| `LabelSummary`      | GET /cookbook/labels 응답 항목              |

### SavedRecipe 필드
```
ID            primitive.ObjectID  `bson:"_id"`
UserID        primitive.ObjectID  `bson:"user_id"`
RecipeID      *primitive.ObjectID `bson:"recipe_id"`      // INTERNAL이면 존재
RecipeSnap    RecipeSnap          `bson:"recipe_snapshot"`
Label         string              `bson:"label"`          // 기본: "미분류"
Note          string              `bson:"note"`
Rating        *int                `bson:"rating"`         // 1-5, nil이면 미평가
SavedAt       time.Time           `bson:"saved_at"`
UpdatedAt     time.Time           `bson:"updated_at"`
```

### RecipeSnap 필드
```
Title           string  `bson:"title"`
SourceURL       string  `bson:"source_url"`
ThumbnailURL    string  `bson:"thumbnail_url"`
SourceType      string  `bson:"source_type"`
MainIngredient  string  `bson:"main_ingredient"`
Category        string  `bson:"category"`
```

---

## repository.go

### 인터페이스 (CookbookRepository)
```
FindAllByUserID(ctx, userID, filter) ([]*SavedRecipe, error)
FindByID(ctx, id, userID) (*SavedRecipe, error)
FindBySourceURL(ctx, sourceURL, userID) (*SavedRecipe, error)  // 중복 방지
FindByRecipeID(ctx, recipeID, userID) (*SavedRecipe, error)    // 중복 방지
Create(ctx, saved) error
Update(ctx, id, userID, update) error
Delete(ctx, id, userID) error
GetLabelSummary(ctx, userID) ([]LabelSummary, error)          // aggregate
```

### MongoDB 컬렉션
- `saved_recipes`
  - unique index: `{user_id: 1, recipe_id: 1}` (partial: recipe_id != null)
  - unique index: `{user_id: 1, "recipe_snapshot.source_url": 1}`
  - index: `{user_id: 1, saved_at: -1}`
  - index: `{user_id: 1, label: 1}`

---

## service.go

### 인터페이스 (CookbookService)
```
ListSavedRecipes(ctx, userID, filter) ([]*SavedRecipeResponse, error)
SaveRecipe(ctx, userID, req SaveRecipeRequest) (*SavedRecipeResponse, error)
GetSavedRecipe(ctx, id, userID) (*SavedRecipeResponse, error)
UpdateSavedRecipe(ctx, id, userID, req UpdateRecipeRequest) (*SavedRecipeResponse, error)
RemoveSavedRecipe(ctx, id, userID) error
GetLabels(ctx, userID) ([]LabelSummary, error)
```

### 핵심 로직 포인트
- SaveRecipe: 중복 저장 방지
  - SourceType == INTERNAL: recipe_id로 중복 체크
  - SourceType == EXTERNAL: source_url로 중복 체크
  - 중복 시 409 Conflict 반환
- Label 미지정 시 기본값 `"미분류"` 자동 설정

---

## handler.go

### 엔드포인트 등록
```
GET    /api/v1/cookbook              (JWT 필요)
POST   /api/v1/cookbook              (JWT 필요)
GET    /api/v1/cookbook/labels       (JWT 필요)
GET    /api/v1/cookbook/:id          (JWT 필요)
PATCH  /api/v1/cookbook/:id         (JWT 필요)
DELETE /api/v1/cookbook/:id          (JWT 필요)
```

---

## 외부 의존성

| 의존 대상               | 용도                              |
|-------------------------|-----------------------------------|
| `internal/recommendation` | SaveRecipeRequest에서 Recipe 정보 참조 (선택적) |
| `pkg/database`          | MongoDB 접근                      |
