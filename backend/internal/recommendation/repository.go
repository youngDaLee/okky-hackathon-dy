package recommendation

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RecipeRepository interface {
	FindAll(ctx context.Context, filter RecommendationRequest) ([]*Recipe, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*Recipe, error)
	Search(ctx context.Context, keyword, category string, limit int) ([]*Recipe, error)
	FindCandidatesByIngredients(ctx context.Context, ingredientNames []string) ([]*Recipe, error)
	EnsureIndexes(ctx context.Context) error
}

type recipeRepository struct {
	col *mongo.Collection
}

func NewRecipeRepository(col *mongo.Collection) RecipeRepository {
	return &recipeRepository{col: col}
}

// Task 1.1 — indexes
func (r *recipeRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "tags", Value: "text"},
			},
			Options: options.Index().SetDefaultLanguage("none"),
		},
		{Keys: bson.D{{Key: "main_ingredient", Value: 1}}},
		{Keys: bson.D{{Key: "required_ingredients", Value: 1}}},
	})
	return err
}

// Task 1.2 — FindCandidatesByIngredients: recipes that share at least one required/main ingredient
func (r *recipeRepository) FindCandidatesByIngredients(ctx context.Context, ingredientNames []string) ([]*Recipe, error) {
	if len(ingredientNames) == 0 {
		return nil, nil
	}
	filter := bson.M{
		"$or": bson.A{
			bson.M{"required_ingredients": bson.M{"$in": ingredientNames}},
			bson.M{"main_ingredient": bson.M{"$in": ingredientNames}},
		},
	}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var recipes []*Recipe
	return recipes, cursor.All(ctx, &recipes)
}

// Task 1.3 — FindByID
func (r *recipeRepository) FindByID(ctx context.Context, id bson.ObjectID) (*Recipe, error) {
	var recipe Recipe
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&recipe)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &recipe, err
}

// Task 1.4 — Search (text search + category filter)
func (r *recipeRepository) Search(ctx context.Context, keyword, category string, limit int) ([]*Recipe, error) {
	if limit <= 0 {
		limit = 20
	}
	filter := bson.M{}
	if keyword != "" {
		filter["$text"] = bson.M{"$search": keyword}
	}
	if category != "" {
		filter["category"] = category
	}
	opt := options.Find().SetLimit(int64(limit))
	if keyword != "" {
		opt.SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	}
	cursor, err := r.col.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var recipes []*Recipe
	return recipes, cursor.All(ctx, &recipes)
}

// Task 1.5 — FindAll (category filter + limit)
func (r *recipeRepository) FindAll(ctx context.Context, filter RecommendationRequest) ([]*Recipe, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	f := bson.M{}
	if filter.Category != "" {
		f["category"] = filter.Category
	}
	cursor, err := r.col.Find(ctx, f, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var recipes []*Recipe
	return recipes, cursor.All(ctx, &recipes)
}
