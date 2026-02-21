package cookbook

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RecipeSnap struct {
	Title          string   `bson:"title"            json:"title"`
	SourceURL      string   `bson:"source_url"       json:"sourceUrl"`
	ThumbnailURL   string   `bson:"thumbnail_url"    json:"thumbnailUrl"`
	SourceType     string   `bson:"source_type"      json:"sourceType"`
	MainIngredient string   `bson:"main_ingredient"  json:"mainIngredient"`
	Category       string   `bson:"category"         json:"category"`
	RawIngredients []string `bson:"raw_ingredients"  json:"rawIngredients"`
}

type SavedRecipe struct {
	ID          bson.ObjectID  `bson:"_id,omitempty"     json:"id"`
	UserID      bson.ObjectID  `bson:"user_id"           json:"userId"`
	RecipeID    *bson.ObjectID `bson:"recipe_id"         json:"recipeId"`
	RecipeSnap  RecipeSnap          `bson:"recipe_snapshot"   json:"recipeSnapshot"`
	Label       string              `bson:"label"             json:"label"`
	Note        string              `bson:"note"              json:"note"`
	Rating      *int                `bson:"rating"            json:"rating"`
	SavedAt     time.Time           `bson:"saved_at"          json:"savedAt"`
	UpdatedAt   time.Time           `bson:"updated_at"        json:"updatedAt"`
}

type SaveRecipeRequest struct {
	RecipeID    *bson.ObjectID `json:"recipeId"`
	RecipeSnap  RecipeSnap          `json:"recipeSnapshot" binding:"required"`
	Label       string              `json:"label"`
	Note        string              `json:"note"  binding:"max=500"`
}

type UpdateRecipeRequest struct {
	Label  *string `json:"label"  binding:"omitempty,max=20"`
	Note   *string `json:"note"   binding:"omitempty,max=500"`
	Rating *int    `json:"rating" binding:"omitempty,min=1,max=5"`
}

type SavedRecipeResponse struct {
	SavedRecipe
}

type LabelSummary struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type ListFilter struct {
	Label  string
	Search string
}
