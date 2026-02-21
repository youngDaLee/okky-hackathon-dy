package fridge

import "context"

type FridgeService interface {
	ListIngredients(ctx context.Context, userID string, filter ListFilter) ([]*IngredientResponse, error)
	GetIngredient(ctx context.Context, id, userID string) (*IngredientResponse, error)
	AddIngredient(ctx context.Context, userID string, req CreateIngredientReq) (*IngredientResponse, error)
	UpdateIngredient(ctx context.Context, id, userID string, req UpdateIngredientReq) (*IngredientResponse, error)
	RemoveIngredient(ctx context.Context, id, userID string) error
	BulkRemove(ctx context.Context, ids []string, userID string) error
	GetSummary(ctx context.Context, userID string) (*FridgeSummary, error)
	BulkAddIngredients(ctx context.Context, userID string, items []CreateIngredientReq) ([]*IngredientResponse, error)
}

type fridgeService struct {
	repo IngredientRepository
}

func NewFridgeService(repo IngredientRepository) FridgeService {
	return &fridgeService{repo: repo}
}

func (s *fridgeService) ListIngredients(ctx context.Context, userID string, filter ListFilter) ([]*IngredientResponse, error) {
	return nil, nil
}

func (s *fridgeService) GetIngredient(ctx context.Context, id, userID string) (*IngredientResponse, error) {
	return nil, nil
}

func (s *fridgeService) AddIngredient(ctx context.Context, userID string, req CreateIngredientReq) (*IngredientResponse, error) {
	return nil, nil
}

func (s *fridgeService) UpdateIngredient(ctx context.Context, id, userID string, req UpdateIngredientReq) (*IngredientResponse, error) {
	return nil, nil
}

func (s *fridgeService) RemoveIngredient(ctx context.Context, id, userID string) error {
	return nil
}

func (s *fridgeService) BulkRemove(ctx context.Context, ids []string, userID string) error {
	return nil
}

func (s *fridgeService) GetSummary(ctx context.Context, userID string) (*FridgeSummary, error) {
	return nil, nil
}

func (s *fridgeService) BulkAddIngredients(ctx context.Context, userID string, items []CreateIngredientReq) ([]*IngredientResponse, error) {
	return nil, nil
}
