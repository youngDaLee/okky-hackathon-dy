package vision

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	TypeReceipt = "RECEIPT"
	TypeFridge  = "FRIDGE"
)

const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusDone       = "DONE"
	StatusFailed     = "FAILED"
)

type ExtractedItem struct {
	Name           string   `bson:"name"            json:"name"`
	NormalizedName string   `bson:"normalized_name" json:"normalizedName"`
	Quantity       *float64 `bson:"quantity"        json:"quantity"`
	Unit           *string  `bson:"unit"            json:"unit"`
	Confidence     float64  `bson:"confidence"      json:"confidence"`
	Selected       bool     `bson:"selected"        json:"selected"`
}

type VisionJob struct {
	ID           bson.ObjectID `bson:"_id,omitempty"   json:"id"`
	UserID       bson.ObjectID `bson:"user_id"         json:"userId"`
	Type         string             `bson:"type"            json:"type"`
	ImageURL     string             `bson:"image_url"       json:"imageUrl"`
	Status       string             `bson:"status"          json:"status"`
	RawResult    interface{}        `bson:"raw_result"      json:"-"`
	Extracted    []ExtractedItem    `bson:"extracted"       json:"extracted"`
	ErrorMessage *string            `bson:"error_message"   json:"errorMessage"`
	CreatedAt    time.Time          `bson:"created_at"      json:"createdAt"`
	CompletedAt  *time.Time         `bson:"completed_at"    json:"completedAt"`
}

type JobResponse struct {
	VisionJob
}

type ConfirmRequest struct {
	Items []ExtractedItem `json:"items" binding:"required"`
}
