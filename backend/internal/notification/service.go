package notification

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationService interface {
	ScheduleAlerts(ctx context.Context, ingredientID bson.ObjectID, ingredientName string, expiryDate time.Time, userID bson.ObjectID) error
	RescheduleAlerts(ctx context.Context, ingredientID bson.ObjectID, ingredientName string, expiryDate time.Time, userID bson.ObjectID) error
	CancelAlerts(ctx context.Context, ingredientID bson.ObjectID) error

	ListAlerts(ctx context.Context, userID string, limit int) ([]*AlertResponse, error)
	CountUnread(ctx context.Context, userID string) (int64, error)
	ReadAlert(ctx context.Context, id, userID string) error
	ReadAllAlerts(ctx context.Context, userID string) error

	ProcessDueAlerts(ctx context.Context) error
}

type notificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) ScheduleAlerts(ctx context.Context, ingredientID bson.ObjectID, ingredientName string, expiryDate time.Time, userID bson.ObjectID) error {
	return nil
}

func (s *notificationService) RescheduleAlerts(ctx context.Context, ingredientID bson.ObjectID, ingredientName string, expiryDate time.Time, userID bson.ObjectID) error {
	return nil
}

func (s *notificationService) CancelAlerts(ctx context.Context, ingredientID bson.ObjectID) error {
	return nil
}

func (s *notificationService) ListAlerts(ctx context.Context, userID string, limit int) ([]*AlertResponse, error) {
	return nil, nil
}

func (s *notificationService) CountUnread(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (s *notificationService) ReadAlert(ctx context.Context, id, userID string) error {
	return nil
}

func (s *notificationService) ReadAllAlerts(ctx context.Context, userID string) error {
	return nil
}

func (s *notificationService) ProcessDueAlerts(ctx context.Context) error {
	return nil
}
