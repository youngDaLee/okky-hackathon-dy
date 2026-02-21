## Context

recommendation 도메인 스켈레톤이 이미 존재한다 (`internal/recommendation/{model,repository,service,handler,external}.go`).
MongoDB driver v2, Gin 프레임워크 사용.
Auth 미들웨어가 `gin.Context`에 `"userID"` 키로 문자열 user_id를 주입한다고 가정.
fridge 도메인은 이미 구현 완료 — FridgeService를 통해 재료 목록 + ExpiryStatus 조회 가능.

## Goals / Non-Goals

**Goals:**
- 4개 API 엔드포인트 완전 동작
- 3-Tier 매칭 알고리즘 (urgency_bonus, match_score)
- 내부 DB 결과 부족 시 외부 YouTube/Google 검색 보완
- 외부 검색 결과 1시간 메모리 캐시
- MongoDB 인덱스 생성 (text, main_ingredient, required_ingredients)

**Non-Goals:**
- 식이 제한(allergens) 필터링 (auth User 모델 미구현)
- 레시피 DB 시딩 및 배치 크롤러 (별도 작업)
- 실시간 외부 API 호출 (런타임에서 제거, 배치로 처리)
- 캐시 (DB 조회만 하므로 불필요)

## Decisions

### 1. FridgeService 직접 주입 (도메인 간 의존)

recommendation 서비스가 FridgeService를 생성자에서 받아 직접 호출.

```go
type recommendationService struct {
    repo     RecipeRepository
    external ExternalSearcher
    fridge   fridge.FridgeService   // 재료 목록 조회
}
```

**대안:** fridge HTTP API 호출 → 내부 함수 직접 호출이 더 단순하고 트랜잭션 경계도 명확.

### 2. 재료 이름 정규화 전략

매칭 시 대소문자 무시 + 앞뒤 공백 제거. 한글 특성상 stemming 불필요.

```go
normalize(s string) string
  → strings.ToLower(strings.TrimSpace(s))
```

레시피 `required_ingredients`와 사용자 재료명 양쪽 모두 normalize 후 비교.

### 3. 매칭 후보 조회 전략

전체 레시피를 스캔하지 않고 `main_ingredient` 또는 `required_ingredients`에 사용자 재료가 하나라도 포함된 후보만 조회 후 인메모리에서 Tier 분류.

```
FindCandidatesByIngredients(ingredientNames []string) []*Recipe
  → { required_ingredients: { $in: ingredientNames } }
```

MVP 500개 레시피 기준 인메모리 분류 성능 문제 없음.

### 4. external.go 역할

배치 크롤러 전용 — 런타임 서비스에서 호출하지 않음.
`ExternalSearcher` 인터페이스는 유지하되, `service.go`는 의존하지 않는다.
`recipes` 컬렉션에 `source_type = EXTERNAL`로 저장된 레시피는
배치가 사전 처리한 것으로 간주하고 DB 조회 시 INTERNAL과 동일하게 취급.

### 5. GET /recommendations/today 구현 전략

오늘의 추천 = URGENT(D-0~D-2) 재료 기반 Tier 1 결과 우선, 없으면 Tier 2.

```go
// service내 today 로직
urgentIngredients = filter(userIngredients, status == URGENT)
candidates = FindCandidatesByIngredients(urgentIngredients.names)
tier1 = match(candidates, urgentIngredients) → Tier1만 추출
if len(tier1) < 3 → fallback: 전체 재료로 GetRecommendations
```

### 7. Limit 기본값 및 빈 냉장고 처리

```
limit 기본값: 20
빈 냉장고(재료 0개): 빈 결과 반환 (인기 레시피 fallback은 MVP 외)
```

## Risks / Trade-offs

- **재료 이름 불일치**: 사용자가 "당근" 입력 vs 레시피 "당근(채썬것)" → normalize로 완화 불가. MVP는 수용, 향후 synonym 테이블 고려.
- **외부 레시피 재료 품질**: AI 추출 재료 목록의 정확도에 따라 Tier 분류 정확도 영향 → 배치에서 신뢰도 낮은 경우 `required_ingredients`를 비워두면 Tier 3 이하로 분류됨.
