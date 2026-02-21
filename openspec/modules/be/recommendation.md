# Module: internal/recommendation

## 패키지 경로
`backend/internal/recommendation`

## 책임
3-Tier 레시피 추천 엔진.
내부 DB 레시피 우선 매칭 + 외부 YouTube/Google 검색 보완.
유통기한 긴급도 가중치 적용.

---

## 파일 구조

```
internal/recommendation/
├── model.go
├── repository.go
├── service.go
├── handler.go
└── external.go
```

---

## model.go

### 구조체

| 구조체                   | 역할                                           |
|--------------------------|------------------------------------------------|
| `Recipe`                 | MongoDB `recipes` 컬렉션 문서 (내부 DB)        |
| `RecommendationResult`   | 추천 결과 응답 DTO (DB 미저장)                 |
| `RecommendationRequest`  | GET /recommendations 쿼리 파라미터 바인딩      |
| `ExternalRecipe`         | 외부 검색 결과 임시 구조체                     |

### Recipe 필드
```
ID                   primitive.ObjectID  `bson:"_id"`
Title                string              `bson:"title"`
Description          string              `bson:"description"`
RequiredIngredients  []string            `bson:"required_ingredients"` // 정규화 이름
OptionalIngredients  []string            `bson:"optional_ingredients"`
MainIngredient       string              `bson:"main_ingredient"`
SourceType           string              `bson:"source_type"`  // INTERNAL / EXTERNAL
SourceURL            string              `bson:"source_url"`
ThumbnailURL         string              `bson:"thumbnail_url"`
Category             string              `bson:"category"`
Tags                 []string            `bson:"tags"`
CookingTimeMin       int                 `bson:"cooking_time_min"`
Difficulty           string              `bson:"difficulty"`   // EASY/MEDIUM/HARD
CreatedAt            time.Time           `bson:"created_at"`
```

### RecommendationResult 필드 (응답 전용)
```
Recipe               Recipe
Tier                 int       // 1, 2, 3
MatchScore           float64   // 0.0 ~ 1.0
MatchedIngredients   []string
MissingIngredients   []string
UrgencyBonus         bool
```

---

## repository.go

### 인터페이스 (RecipeRepository)
```
// 내부 DB 레시피 조회
FindAll(ctx, filter) ([]*Recipe, error)
FindByID(ctx, id) (*Recipe, error)
Search(ctx, keyword, category string, limit int) ([]*Recipe, error)

// 주재료 기준 후보 레시피 조회 (매칭 엔진용)
FindCandidatesByIngredients(ctx, ingredientNames []string) ([]*Recipe, error)
```

### MongoDB 컬렉션
- `recipes`
  - text index: `{title: "text", tags: "text"}`
  - index: `{main_ingredient: 1}`
  - index: `{required_ingredients: 1}`

---

## service.go

### 인터페이스 (RecommendationService)
```
GetRecommendations(ctx, userID string, req RecommendationRequest) ([]RecommendationResult, error)
GetTodayRecommendations(ctx, userID string) ([]RecommendationResult, error)
SearchRecipes(ctx, keyword, category string, limit int) ([]*Recipe, error)
GetRecipeByID(ctx, id string) (*Recipe, error)
```

### 핵심 로직 포인트 (service 내부)

#### match() — Tier 분류 함수
```
입력: recipe *Recipe, userIngredients map[string]bool, urgentNames map[string]bool

matchCount = required_ingredients ∩ userIngredients 개수
matchScore = matchCount / len(required_ingredients)

if matchScore == 1.0           → Tier 1
if mainIngredient ∈ user AND matchScore >= 0.6 → Tier 2
if matchScore >= 0.3           → Tier 3
else                           → 제외

urgencyBonus = required_ingredients ∩ urgentNames 개수 > 0
```

#### 외부 보완 기준
```
if len(tier1Results) < 3  → external.SearchRecipes(주재료, limit=5) 호출
if len(tier2Results) < 5  → external.SearchRecipes(주재료, limit=3) 호출
```

---

## handler.go

### 엔드포인트 등록
```
GET    /api/v1/recommendations          (JWT 필요)
GET    /api/v1/recommendations/today    (JWT 필요)
GET    /api/v1/recipes                  (JWT 필요)
GET    /api/v1/recipes/:id              (JWT 필요)
```

### GET /api/v1/recommendations 쿼리 파라미터
```
tier        int    (1/2/3, 없으면 전체)
category    string
max_missing int
limit       int    (기본 20)
```

---

## external.go

### 역할
외부 API (YouTube Data API v3, Google Custom Search) 연동.
내부 레시피가 부족할 때 service.go에서 호출.

### 인터페이스 (ExternalSearcher)
```
SearchRecipes(ctx, query string, limit int) ([]*ExternalRecipe, error)
```

### ExternalRecipe 구조체
```
Title        string
SourceURL    string
ThumbnailURL string
SourceType   string  // "YOUTUBE" / "BLOG"
```

### 검색 전략
- YouTube: `query + " 레시피"` → YouTube Data API v3 search.list
- Google: Custom Search API → 요리 블로그 도메인 필터

### 캐싱
- 쿼리 문자열 기준 1시간 메모리 캐시 (sync.Map 또는 go-cache)

---

## 외부 의존성

| 의존 대상         | 용도                                    |
|-------------------|-----------------------------------------|
| `internal/fridge` | 사용자 재료 목록 + URGENT 재료 조회     |
| `internal/auth`   | dietary_prefs, allergens 조회           |
| `pkg/database`    | MongoDB 접근 (recipes 컬렉션)           |
| YouTube Data API  | 외부 레시피 링크 검색                   |
| Google Custom Search | 블로그 레시피 검색                   |
