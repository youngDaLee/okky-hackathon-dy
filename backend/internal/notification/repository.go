package notification

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationRepository interface {
	BulkCreate(ctx context.Context, alerts []*ExpiryAlert) error
	FindByUserID(ctx context.Context, userID bson.ObjectID, limit int) ([]*ExpiryAlert, error)
	FindUnsentByDate(ctx context.Context, date time.Time) ([]*ExpiryAlert, error)
	CountUnread(ctx context.Context, userID bson.ObjectID) (int64, error)
	MarkSent(ctx context.Context, ids []bson.ObjectID) error
	MarkRead(ctx context.Context, id, userID bson.ObjectID) error
	MarkAllRead(ctx context.Context, userID bson.ObjectID) error
	DeleteByIngredientID(ctx context.Context, ingredientID bson.ObjectID) error
	DeleteUnsentByIngredientID(ctx context.Context, ingredientID bson.ObjectID) error
}

type notificationRepository struct {
	col *mongo.Collection
}

func NewNotificationRepository(col *mongo.Collection) NotificationRepository {
	return &notificationRepository{col: col}
}

func (r *notificationRepository) BulkCreate(ctx context.Context, alerts []*ExpiryAlert) error {
	return nil
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID bson.ObjectID, limit int) ([]*ExpiryAlert, error) {
	return nil, nil
}

func (r *notificationRepository) FindUnsentByDate(ctx context.Context, date time.Time) ([]*ExpiryAlert, error) {
	return nil, nil
}

func (r *notificationRepository) CountUnread(ctx context.Context, userID bson.ObjectID) (int64, error) {
	return 0, nil
}

func (r *notificationRepository) MarkSent(ctx context.Context, ids []bson.ObjectID) error {
	return nil
}

func (r *notificationRepository) MarkRead(ctx context.Context, id, userID bson.ObjectID) error {
	return nil
}

func (r *notificationRepository) MarkAllRead(ctx context.Context, userID bson.ObjectID) error {
	return nil
}

func (r *notificationRepository) DeleteByIngredientID(ctx context.Context, ingredientID bson.ObjectID) error {
	return nil
}

func (r *notificationRepository) DeleteUnsentByIngredientID(ctx context.Context, ingredientID bson.ObjectID) error {
	return nil
}
