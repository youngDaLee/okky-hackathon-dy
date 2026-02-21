package cookbook

import "context"

type CookbookService interface {
	ListSavedRecipes(ctx context.Context, userID string, filter ListFilter) ([]*SavedRecipeResponse, error)
	SaveRecipe(ctx context.Context, userID string, req SaveRecipeRequest) (*SavedRecipeResponse, error)
	GetSavedRecipe(ctx context.Context, id, userID string) (*SavedRecipeResponse, error)
	UpdateSavedRecipe(ctx context.Context, id, userID string, req UpdateRecipeRequest) (*SavedRecipeResponse, error)
	RemoveSavedRecipe(ctx context.Context, id, userID string) error
	GetLabels(ctx context.Context, userID string) ([]LabelSummary, error)
}

type cookbookService struct {
	repo CookbookRepository
}

func NewCookbookService(repo CookbookRepository) CookbookService {
	return &cookbookService{repo: repo}
}

func (s *cookbookService) ListSavedRecipes(ctx context.Context, userID string, filter ListFilter) ([]*SavedRecipeResponse, error) {
	return nil, nil
}

func (s *cookbookService) SaveRecipe(ctx context.Context, userID string, req SaveRecipeRequest) (*SavedRecipeResponse, error) {
	return nil, nil
}

func (s *cookbookService) GetSavedRecipe(ctx context.Context, id, userID string) (*SavedRecipeResponse, error) {
	return nil, nil
}

func (s *cookbookService) UpdateSavedRecipe(ctx context.Context, id, userID string, req UpdateRecipeRequest) (*SavedRecipeResponse, error) {
	return nil, nil
}

func (s *cookbookService) RemoveSavedRecipe(ctx context.Context, id, userID string) error {
	return nil
}

func (s *cookbookService) GetLabels(ctx context.Context, userID string) ([]LabelSummary, error) {
	return nil, nil
}
