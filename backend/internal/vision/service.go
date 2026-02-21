package vision

import (
	"context"
	"mime/multipart"
)

type VisionService interface {
	CreateJob(ctx context.Context, userID string, jobType string, file multipart.File, filename string) (*JobResponse, error)
	GetJob(ctx context.Context, id, userID string) (*JobResponse, error)
	ConfirmJob(ctx context.Context, id, userID string, req ConfirmRequest) error
}

type visionService struct {
	repo VisionRepository
}

func NewVisionService(repo VisionRepository) VisionService {
	return &visionService{repo: repo}
}

func (s *visionService) CreateJob(ctx context.Context, userID string, jobType string, file multipart.File, filename string) (*JobResponse, error) {
	return nil, nil
}

func (s *visionService) GetJob(ctx context.Context, id, userID string) (*JobResponse, error) {
	return nil, nil
}

func (s *visionService) ConfirmJob(ctx context.Context, id, userID string, req ConfirmRequest) error {
	return nil
}
