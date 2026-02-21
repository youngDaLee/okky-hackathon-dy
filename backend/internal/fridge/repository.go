package fridge

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ExpiryStatusCount struct {
	Urgent int
	Soon   int
	Total  int
}

type IngredientRepository interface {
	FindAllByUserID(ctx context.Context, userID bson.ObjectID, filter ListFilter) ([]*Ingredient, error)
	FindByID(ctx context.Context, id, userID bson.ObjectID) (*Ingredient, error)
	Create(ctx context.Context, ingredient *Ingredient) error
	Update(ctx context.Context, id, userID bson.ObjectID, update *Ingredient) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error
	BulkDelete(ctx context.Context, ids []bson.ObjectID, userID bson.ObjectID) error
	CountByExpiryStatus(ctx context.Context, userID bson.ObjectID) (*ExpiryStatusCount, error)
}

type ingredientRepository struct {
	col *mongo.Collection
}

func NewIngredientRepository(col *mongo.Collection) IngredientRepository {
	return &ingredientRepository{col: col}
}

func (r *ingredientRepository) FindAllByUserID(ctx context.Context, userID bson.ObjectID, filter ListFilter) ([]*Ingredient, error) {
	return nil, nil
}

func (r *ingredientRepository) FindByID(ctx context.Context, id, userID bson.ObjectID) (*Ingredient, error) {
	return nil, nil
}

func (r *ingredientRepository) Create(ctx context.Context, ingredient *Ingredient) error {
	return nil
}

func (r *ingredientRepository) Update(ctx context.Context, id, userID bson.ObjectID, update *Ingredient) error {
	return nil
}

func (r *ingredientRepository) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	return nil
}

func (r *ingredientRepository) BulkDelete(ctx context.Context, ids []bson.ObjectID, userID bson.ObjectID) error {
	return nil
}

func (r *ingredientRepository) CountByExpiryStatus(ctx context.Context, userID bson.ObjectID) (*ExpiryStatusCount, error) {
	return nil, nil
}
