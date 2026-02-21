package recommendation

import "context"

type RecommendationService interface {
	GetRecommendations(ctx context.Context, userID string, req RecommendationRequest) ([]RecommendationResult, error)
	GetTodayRecommendations(ctx context.Context, userID string) ([]RecommendationResult, error)
	SearchRecipes(ctx context.Context, keyword, category string, limit int) ([]*Recipe, error)
	GetRecipeByID(ctx context.Context, id string) (*Recipe, error)
}

type recommendationService struct {
	repo     RecipeRepository
	external ExternalSearcher
}

func NewRecommendationService(repo RecipeRepository, external ExternalSearcher) RecommendationService {
	return &recommendationService{repo: repo, external: external}
}

func (s *recommendationService) GetRecommendations(ctx context.Context, userID string, req RecommendationRequest) ([]RecommendationResult, error) {
	return nil, nil
}

func (s *recommendationService) GetTodayRecommendations(ctx context.Context, userID string) ([]RecommendationResult, error) {
	return nil, nil
}

func (s *recommendationService) SearchRecipes(ctx context.Context, keyword, category string, limit int) ([]*Recipe, error) {
	return nil, nil
}

func (s *recommendationService) GetRecipeByID(ctx context.Context, id string) (*Recipe, error) {
	return nil, nil
}
