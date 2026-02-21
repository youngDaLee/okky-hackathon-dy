package recommendation

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Recipe struct {
	ID                  bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title               string             `bson:"title"                json:"title"`
	Description         string             `bson:"description"          json:"description"`
	RequiredIngredients []string           `bson:"required_ingredients" json:"requiredIngredients"`
	OptionalIngredients []string           `bson:"optional_ingredients" json:"optionalIngredients"`
	RawIngredients      []string           `bson:"raw_ingredients"      json:"rawIngredients"`
	MainIngredient      string             `bson:"main_ingredient"      json:"mainIngredient"`
	SourceType          string             `bson:"source_type"          json:"sourceType"`
	SourceURL           string             `bson:"source_url"           json:"sourceUrl"`
	ThumbnailURL        string             `bson:"thumbnail_url"        json:"thumbnailUrl"`
	Category            string             `bson:"category"             json:"category"`
	Tags                []string           `bson:"tags"                 json:"tags"`
	CookingTimeMin      int                `bson:"cooking_time_min"     json:"cookingTimeMin"`
	Difficulty          string             `bson:"difficulty"           json:"difficulty"`
	CreatedAt           time.Time          `bson:"created_at"           json:"createdAt"`
}

type RecommendationResult struct {
	Recipe             Recipe   `json:"recipe"`
	Tier               int      `json:"tier"`
	MatchScore         float64  `json:"matchScore"`
	MatchedIngredients []string `json:"matchedIngredients"`
	MissingIngredients []string `json:"missingIngredients"`
	UrgencyBonus       bool     `json:"urgencyBonus"`
}

type RecommendationRequest struct {
	Tier       *int   `form:"tier"`
	Category   string `form:"category"`
	MaxMissing *int   `form:"max_missing"`
	Limit      int    `form:"limit"`
}

type ExternalRecipe struct {
	Title        string `json:"title"`
	SourceURL    string `json:"sourceUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
	SourceType   string `json:"sourceType"`
}
