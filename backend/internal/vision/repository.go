package vision

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type VisionRepository interface {
	Create(ctx context.Context, job *VisionJob) error
	FindByID(ctx context.Context, id, userID bson.ObjectID) (*VisionJob, error)
	UpdateStatus(ctx context.Context, id bson.ObjectID, status string, result *VisionJob) error
	DeleteExpired(ctx context.Context, before time.Time) error
}

type visionRepository struct {
	col *mongo.Collection
}

func NewVisionRepository(col *mongo.Collection) VisionRepository {
	return &visionRepository{col: col}
}

func (r *visionRepository) Create(ctx context.Context, job *VisionJob) error {
	return nil
}

func (r *visionRepository) FindByID(ctx context.Context, id, userID bson.ObjectID) (*VisionJob, error) {
	return nil, nil
}

func (r *visionRepository) UpdateStatus(ctx context.Context, id bson.ObjectID, status string, result *VisionJob) error {
	return nil
}

func (r *visionRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return nil
}
