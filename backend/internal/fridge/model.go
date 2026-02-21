package fridge

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	CategoryVegetable = "VEGETABLE"
	CategoryFruit     = "FRUIT"
	CategoryMeat      = "MEAT"
	CategorySeafood   = "SEAFOOD"
	CategoryDairy     = "DAIRY"
	CategoryGrain     = "GRAIN"
	CategoryCondiment = "CONDIMENT"
	CategoryFrozen    = "FROZEN"
	CategoryOther     = "OTHER"
)

const (
	ExpiryStatusUrgent   = "URGENT"
	ExpiryStatusSoon     = "SOON"
	ExpiryStatusNormal   = "NORMAL"
	ExpiryStatusNoExpiry = "NO_EXPIRY"
)

const (
	SourceManual  = "manual"
	SourceReceipt = "receipt"
	SourceVision  = "vision"
)

type Ingredient struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     bson.ObjectID `bson:"user_id"       json:"userId"`
	Name       string             `bson:"name"          json:"name"`
	Category   string             `bson:"category"      json:"category"`
	Quantity   float64            `bson:"quantity"      json:"quantity"`
	Unit       string             `bson:"unit"          json:"unit"`
	ExpiryDate *time.Time         `bson:"expiry_date"   json:"expiryDate"`
	Source     string             `bson:"source"        json:"source"`
	AddedAt    time.Time          `bson:"added_at"      json:"addedAt"`
	UpdatedAt  time.Time          `bson:"updated_at"    json:"updatedAt"`
}

type CreateIngredientReq struct {
	Name       string     `json:"name"       binding:"required,max=50"`
	Category   string     `json:"category"   binding:"required"`
	Quantity   float64    `json:"quantity"   binding:"required,gt=0"`
	Unit       string     `json:"unit"       binding:"required"`
	ExpiryDate *time.Time `json:"expiryDate"`
}

type UpdateIngredientReq struct {
	Name       *string    `json:"name"`
	Category   *string    `json:"category"`
	Quantity   *float64   `json:"quantity"`
	Unit       *string    `json:"unit"`
	ExpiryDate *time.Time `json:"expiryDate"`
}

type IngredientResponse struct {
	Ingredient
	ExpiryStatus string `json:"expiryStatus"`
}

type FridgeSummary struct {
	Total  int `json:"total"`
	Urgent int `json:"urgent"`
	Soon   int `json:"soon"`
}

type BulkDeleteReq struct {
	IDs []bson.ObjectID `json:"ids" binding:"required"`
}

type ListFilter struct {
	Category     string
	ExpiryStatus string
	Search       string
}
