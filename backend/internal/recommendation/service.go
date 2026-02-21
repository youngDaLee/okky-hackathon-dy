package recommendation

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"okky-hackathon/fridge-master-backend/internal/fridge"
)

var ErrNotFound = errors.New("recipe not found")

type FridgeService interface {
	ListIngredients(ctx context.Context, userID string, filter fridge.ListFilter) ([]*fridge.IngredientResponse, int, error)
}

type RecommendationService interface {
	GetRecommendations(ctx context.Context, userID string, req RecommendationRequest) ([]RecommendationResult, error)
	GetTodayRecommendations(ctx context.Context, userID string) ([]RecommendationResult, error)
	SearchRecipes(ctx context.Context, keyword, category string, limit int) ([]*Recipe, error)
	GetRecipeByID(ctx context.Context, id string) (*Recipe, error)
}

type recommendationService struct {
	repo   RecipeRepository
	fridge FridgeService
}

func NewRecommendationService(repo RecipeRepository, fridge FridgeService) RecommendationService {
	return &recommendationService{repo: repo, fridge: fridge}
}

// Task 2.1 — normalize
func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Task 2.2 — matchRecipe: Tier 분류, match_score, urgency_bonus
func matchRecipe(recipe *Recipe, userIngMap, urgentMap map[string]bool) *RecommendationResult {
	if len(recipe.RequiredIngredients) == 0 {
		return nil
	}

	var matched, missing []string
	for _, ing := range recipe.RequiredIngredients {
		key := normalize(ing)
		if userIngMap[key] {
			matched = append(matched, ing)
		} else {
			missing = append(missing, ing)
		}
	}

	score := float64(len(matched)) / float64(len(recipe.RequiredIngredients))

	var tier int
	switch {
	case score == 1.0:
		tier = 1
	case userIngMap[normalize(recipe.MainIngredient)] && score >= 0.6:
		tier = 2
	case score >= 0.3:
		tier = 3
	default:
		return nil // 매칭 제외
	}

	urgencyBonus := false
	for _, ing := range recipe.RequiredIngredients {
		if urgentMap[normalize(ing)] {
			urgencyBonus = true
			break
		}
	}

	return &RecommendationResult{
		Recipe:             *recipe,
		Tier:               tier,
		MatchScore:         score,
		MatchedIngredients: matched,
		MissingIngredients: missing,
		UrgencyBonus:       urgencyBonus,
	}
}

// Task 2.3 — 사용자 재료 조회 및 URGENT 추출
func (s *recommendationService) getUserIngredients(ctx context.Context, userID string) (userIngMap, urgentMap map[string]bool, err error) {
	items, _, err := s.fridge.ListIngredients(ctx, userID, fridge.ListFilter{})
	if err != nil {
		return nil, nil, err
	}
	userIngMap = make(map[string]bool, len(items))
	urgentMap = make(map[string]bool)
	for _, item := range items {
		key := normalize(item.Name)
		userIngMap[key] = true
		if item.ExpiryStatus == fridge.ExpiryStatusUrgent {
			urgentMap[key] = true
		}
	}
	return userIngMap, urgentMap, nil
}

// sortResults: urgency_bonus DESC → match_score DESC (in-place)
func sortResults(results []RecommendationResult) {
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			a, b := results[i], results[j]
			if (!a.UrgencyBonus && b.UrgencyBonus) ||
				(a.UrgencyBonus == b.UrgencyBonus && a.MatchScore < b.MatchScore) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}

// ingredientKeys returns normalized names as a slice.
func ingredientKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Task 2.4 — GetRecommendations
func (s *recommendationService) GetRecommendations(ctx context.Context, userID string, req RecommendationRequest) ([]RecommendationResult, error) {
	userIngMap, urgentMap, err := s.getUserIngredients(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(userIngMap) == 0 {
		return []RecommendationResult{}, nil
	}

	candidates, err := s.repo.FindCandidatesByIngredients(ctx, ingredientKeys(userIngMap))
	if err != nil {
		return nil, err
	}

	var results []RecommendationResult
	for _, recipe := range candidates {
		if req.Category != "" && recipe.Category != req.Category {
			continue
		}
		r := matchRecipe(recipe, userIngMap, urgentMap)
		if r == nil {
			continue
		}
		if req.Tier != nil && r.Tier != *req.Tier {
			continue
		}
		if req.MaxMissing != nil && len(r.MissingIngredients) > *req.MaxMissing {
			continue
		}
		results = append(results, *r)
	}

	sortResults(results)

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

// Task 2.5 — GetTodayRecommendations
func (s *recommendationService) GetTodayRecommendations(ctx context.Context, userID string) ([]RecommendationResult, error) {
	userIngMap, urgentMap, err := s.getUserIngredients(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(userIngMap) == 0 {
		return []RecommendationResult{}, nil
	}

	// URGENT 재료 기반 후보 조회
	var candidates []*Recipe
	if len(urgentMap) > 0 {
		candidates, err = s.repo.FindCandidatesByIngredients(ctx, ingredientKeys(urgentMap))
		if err != nil {
			return nil, err
		}
	}

	// Tier 1 추출
	var tier1 []RecommendationResult
	for _, recipe := range candidates {
		r := matchRecipe(recipe, userIngMap, urgentMap)
		if r != nil && r.Tier == 1 {
			tier1 = append(tier1, *r)
		}
	}

	if len(tier1) > 0 {
		sortResults(tier1)
		if len(tier1) > 10 {
			tier1 = tier1[:10]
		}
		return tier1, nil
	}

	// URGENT 없거나 Tier 1 없으면 전체 재료로 fallback
	return s.GetRecommendations(ctx, userID, RecommendationRequest{Limit: 10})
}

// Task 2.6 — SearchRecipes
func (s *recommendationService) SearchRecipes(ctx context.Context, keyword, category string, limit int) ([]*Recipe, error) {
	return s.repo.Search(ctx, keyword, category, limit)
}

// Task 2.7 — GetRecipeByID
func (s *recommendationService) GetRecipeByID(ctx context.Context, id string) (*Recipe, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}
	recipe, err := s.repo.FindByID(ctx, oid)
	if err != nil || recipe == nil {
		return nil, ErrNotFound
	}
	return recipe, nil
}
