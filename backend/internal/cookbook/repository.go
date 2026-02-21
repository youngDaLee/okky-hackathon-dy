package cookbook

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CookbookRepository interface {
	FindAllByUserID(ctx context.Context, userID bson.ObjectID, filter ListFilter) ([]*SavedRecipe, error)
	FindByID(ctx context.Context, id, userID bson.ObjectID) (*SavedRecipe, error)
	FindBySourceURL(ctx context.Context, sourceURL string, userID bson.ObjectID) (*SavedRecipe, error)
	FindByRecipeID(ctx context.Context, recipeID, userID bson.ObjectID) (*SavedRecipe, error)
	Create(ctx context.Context, saved *SavedRecipe) error
	Update(ctx context.Context, id, userID bson.ObjectID, update *SavedRecipe) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error
	GetLabelSummary(ctx context.Context, userID bson.ObjectID) ([]LabelSummary, error)
}

type cookbookRepository struct {
	col *mongo.Collection
}

func NewCookbookRepository(col *mongo.Collection) CookbookRepository {
	return &cookbookRepository{col: col}
}

func (r *cookbookRepository) FindAllByUserID(ctx context.Context, userID bson.ObjectID, filter ListFilter) ([]*SavedRecipe, error) {
	return nil, nil
}

func (r *cookbookRepository) FindByID(ctx context.Context, id, userID bson.ObjectID) (*SavedRecipe, error) {
	return nil, nil
}

func (r *cookbookRepository) FindBySourceURL(ctx context.Context, sourceURL string, userID bson.ObjectID) (*SavedRecipe, error) {
	return nil, nil
}

func (r *cookbookRepository) FindByRecipeID(ctx context.Context, recipeID, userID bson.ObjectID) (*SavedRecipe, error) {
	return nil, nil
}

func (r *cookbookRepository) Create(ctx context.Context, saved *SavedRecipe) error {
	return nil
}

func (r *cookbookRepository) Update(ctx context.Context, id, userID bson.ObjectID, update *SavedRecipe) error {
	return nil
}

func (r *cookbookRepository) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	return nil
}

func (r *cookbookRepository) GetLabelSummary(ctx context.Context, userID bson.ObjectID) ([]LabelSummary, error) {
	return nil, nil
}
