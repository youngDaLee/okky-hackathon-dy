package notification

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	AlertTypeD3  = "D_3"
	AlertTypeD1  = "D_1"
	AlertTypeDay = "D_DAY"
)

type ExpiryAlert struct {
	ID             bson.ObjectID `bson:"_id,omitempty"    json:"id"`
	UserID         bson.ObjectID `bson:"user_id"          json:"userId"`
	IngredientID   bson.ObjectID `bson:"ingredient_id"    json:"ingredientId"`
	IngredientName string             `bson:"ingredient_name"  json:"ingredientName"`
	Type           string             `bson:"type"             json:"type"`
	ScheduledDate  time.Time          `bson:"scheduled_date"   json:"scheduledDate"`
	SentAt         *time.Time         `bson:"sent_at"          json:"sentAt"`
	IsRead         bool               `bson:"is_read"          json:"isRead"`
	CreatedAt      time.Time          `bson:"created_at"       json:"createdAt"`
}

type AlertResponse struct {
	ExpiryAlert
}

type AlertCountResponse struct {
	Count int64 `json:"count"`
}
