## 1. Repository

- [x] 1.1 `internal/recommendation/repository.go` — MongoDB `recipes` 컬렉션 인덱스 생성 (text: title+tags, main_ingredient, required_ingredients)
- [x] 1.2 `internal/recommendation/repository.go` — `FindCandidatesByIngredients` 구현 (`required_ingredients: {$in: names}` 쿼리)
- [x] 1.3 `internal/recommendation/repository.go` — `FindByID` 구현 (ObjectID 파싱 포함)
- [x] 1.4 `internal/recommendation/repository.go` — `Search` 구현 (keyword $text 검색, category 필터, limit)
- [x] 1.5 `internal/recommendation/repository.go` — `FindAll` 구현 (category 필터, limit)

## 2. Service — 매칭 엔진

- [x] 2.1 `internal/recommendation/service.go` — `normalize(s string) string` 구현 (소문자 + 공백 제거)
- [x] 2.2 `internal/recommendation/service.go` — `matchRecipe(recipe, userIngMap, urgentMap)` 구현
  - match_score 계산, Tier 1/2/3 분류, urgency_bonus 판별
- [x] 2.3 `internal/recommendation/service.go` — FridgeService 주입 및 사용자 재료 조회
  - userID → 재료 목록 + ExpiryStatus 조회, URGENT 재료 별도 추출
- [x] 2.4 `internal/recommendation/service.go` — `GetRecommendations` 구현
  - 후보 조회 → matchRecipe → Tier별 정렬 (urgency_bonus DESC, match_score DESC)
  - tier/category/max_missing/limit 필터 적용
- [x] 2.5 `internal/recommendation/service.go` — `GetTodayRecommendations` 구현
  - URGENT 재료 기반 Tier 1 우선, 없으면 전체 재료 Tier 1/2 fallback
- [x] 2.6 `internal/recommendation/service.go` — `SearchRecipes` 구현 (repo.Search 위임)
- [x] 2.7 `internal/recommendation/service.go` — `GetRecipeByID` 구현 (404 처리)

## 3. Handler — HTTP 핸들러

- [ ] 3.1 `internal/recommendation/handler.go` — `GET /recommendations` 핸들러 (쿼리 파라미터 바인딩, tier별 그룹 응답)
- [ ] 3.2 `internal/recommendation/handler.go` — `GET /recommendations/today` 핸들러
- [ ] 3.3 `internal/recommendation/handler.go` — `GET /recipes` 핸들러 (keyword/category/limit)
- [ ] 3.4 `internal/recommendation/handler.go` — `GET /recipes/:id` 핸들러 (404 처리)

## 4. Router 연결 및 서버 부트스트랩

- [x] 4.1 `internal/server/router.go` — recommendation 라우트 4개 등록 (auth 미들웨어 적용)
- [x] 4.2 `cmd/server/main.go` — RecipeRepository, RecommendationService 생성 및 주입 (FridgeService 공유)
