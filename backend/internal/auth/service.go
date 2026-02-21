package auth

import "context"

type AuthService interface {
	Signup(ctx context.Context, req SignupRequest) (*TokenResponse, error)
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetMe(ctx context.Context, userID string) (*UserResponse, error)
	UpdateMe(ctx context.Context, userID string, req UpdateMeRequest) (*UserResponse, error)
}

type authService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Signup(ctx context.Context, req SignupRequest) (*TokenResponse, error) {
	return nil, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	return nil, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	return nil, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return nil
}

func (s *authService) GetMe(ctx context.Context, userID string) (*UserResponse, error) {
	return nil, nil
}

func (s *authService) UpdateMe(ctx context.Context, userID string, req UpdateMeRequest) (*UserResponse, error) {
	return nil, nil
}
