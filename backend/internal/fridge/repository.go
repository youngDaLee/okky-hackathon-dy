package fridge

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ExpiryStatusCount struct {
	Urgent   int
	Soon     int
	Normal   int
	NoExpiry int
	Total    int
}

type IngredientRepository interface {
	FindAllByUserID(ctx context.Context, userID bson.ObjectID, filter ListFilter) ([]*Ingredient, error)
	FindByID(ctx context.Context, id, userID bson.ObjectID) (*Ingredient, error)
	FindByNameAndUserID(ctx context.Context, name string, userID bson.ObjectID) (*Ingredient, error)
	Create(ctx context.Context, ingredient *Ingredient) error
	Update(ctx context.Context, id, userID bson.ObjectID, fields bson.M) (*Ingredient, error)
	Delete(ctx context.Context, id, userID bson.ObjectID) error
	BulkDelete(ctx context.Context, ids []bson.ObjectID, userID bson.ObjectID) (int64, error)
	CountByUserID(ctx context.Context, userID bson.ObjectID) (int64, error)
	EnsureIndexes(ctx context.Context) error
}

type ingredientRepository struct {
	col *mongo.Collection
}

func NewIngredientRepository(col *mongo.Collection) IngredientRepository {
	return &ingredientRepository{col: col}
}

// EnsureIndexes creates required MongoDB indexes.
func (r *ingredientRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "expiry_date", Value: 1},
			},
		},
		{
			Keys:    bson.D{{Key: "name", Value: "text"}},
			Options: options.Index().SetDefaultLanguage("none"),
		},
	}
	_, err := r.col.Indexes().CreateMany(ctx, indexes)
	return err
}

// Task 1.2 — FindAllByUserID with filters
func (r *ingredientRepository) FindAllByUserID(ctx context.Context, userID bson.ObjectID, filter ListFilter) ([]*Ingredient, error) {
	f := bson.M{"user_id": userID}

	if filter.Category != "" {
		f["category"] = filter.Category
	}
	if filter.Search != "" {
		f["name"] = bson.M{"$regex": "^" + strings.ToLower(filter.Search), "$options": "i"}
	}

	cursor, err := r.col.Find(ctx, f, options.Find().SetSort(bson.D{{Key: "added_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*Ingredient
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// Task 1.3 — FindByID (ownership enforced by including userID)
func (r *ingredientRepository) FindByID(ctx context.Context, id, userID bson.ObjectID) (*Ingredient, error) {
	var item Ingredient
	err := r.col.FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &item, err
}

// Task 1.9 — FindByNameAndUserID
func (r *ingredientRepository) FindByNameAndUserID(ctx context.Context, name string, userID bson.ObjectID) (*Ingredient, error) {
	var item Ingredient
	err := r.col.FindOne(ctx, bson.M{
		"name":    bson.M{"$regex": "^" + name + "$", "$options": "i"},
		"user_id": userID,
	}).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &item, err
}

// Task 1.4 — Create
func (r *ingredientRepository) Create(ctx context.Context, ingredient *Ingredient) error {
	res, err := r.col.InsertOne(ctx, ingredient)
	if err != nil {
		return err
	}
	ingredient.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

// Task 1.5 — Update (partial $set)
func (r *ingredientRepository) Update(ctx context.Context, id, userID bson.ObjectID, fields bson.M) (*Ingredient, error) {
	after := options.After
	opt := options.FindOneAndUpdate().SetReturnDocument(after)
	var updated Ingredient
	err := r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id, "user_id": userID},
		bson.M{"$set": fields},
		opt,
	).Decode(&updated)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &updated, err
}

// Task 1.6 — DeleteByID
func (r *ingredientRepository) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	return err
}

// Task 1.7 — BulkDelete (silent ignore — DeleteMany)
func (r *ingredientRepository) BulkDelete(ctx context.Context, ids []bson.ObjectID, userID bson.ObjectID) (int64, error) {
	res, err := r.col.DeleteMany(ctx, bson.M{
		"_id":     bson.M{"$in": ids},
		"user_id": userID,
	})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// Task 1.8 — CountByUserID
func (r *ingredientRepository) CountByUserID(ctx context.Context, userID bson.ObjectID) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{"user_id": userID})
}
