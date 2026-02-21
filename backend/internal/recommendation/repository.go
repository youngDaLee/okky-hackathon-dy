package recommendation

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RecipeRepository interface {
	FindAll(ctx context.Context, filter RecommendationRequest) ([]*Recipe, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*Recipe, error)
	Search(ctx context.Context, keyword, category string, limit int) ([]*Recipe, error)
	FindCandidatesByIngredients(ctx context.Context, ingredientNames []string) ([]*Recipe, error)
}

type recipeRepository struct {
	col *mongo.Collection
}

func NewRecipeRepository(col *mongo.Collection) RecipeRepository {
	return &recipeRepository{col: col}
}

func (r *recipeRepository) FindAll(ctx context.Context, filter RecommendationRequest) ([]*Recipe, error) {
	return nil, nil
}

func (r *recipeRepository) FindByID(ctx context.Context, id bson.ObjectID) (*Recipe, error) {
	return nil, nil
}

func (r *recipeRepository) Search(ctx context.Context, keyword, category string, limit int) ([]*Recipe, error) {
	return nil, nil
}

func (r *recipeRepository) FindCandidatesByIngredients(ctx context.Context, ingredientNames []string) ([]*Recipe, error) {
	return nil, nil
}
