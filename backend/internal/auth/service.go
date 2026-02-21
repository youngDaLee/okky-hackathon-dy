package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

// Sentinel errors for handler mapping
var (
	ErrDuplicateEmail      = errors.New("DUPLICATE_EMAIL")
	ErrInvalidCredentials  = errors.New("INVALID_CREDENTIALS")
	ErrInvalidRefreshToken = errors.New("INVALID_REFRESH_TOKEN")
	ErrNotFound            = errors.New("NOT_FOUND")
	ErrValidation          = errors.New("VALIDATION_ERROR")
)

var passwordRegex = regexp.MustCompile(`[a-zA-Z]`)
var passwordDigitRegex = regexp.MustCompile(`[0-9]`)

type AuthService interface {
	Signup(ctx context.Context, req SignupRequest) (*TokenResponse, error)
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetMe(ctx context.Context, userID string) (*UserResponse, error)
	UpdateMe(ctx context.Context, userID string, req UpdateMeRequest) (*UserResponse, error)
	ValidateAccessToken(tokenStr string) (string, error)
}

type authService struct {
	repo          UserRepository
	jwtSecret     []byte
	jwtAccessTTL  time.Duration
	jwtRefreshTTL time.Duration
}

func NewAuthService(repo UserRepository, jwtSecret string, accessTTL, refreshTTL int) AuthService {
	return &authService{
		repo:          repo,
		jwtSecret:     []byte(jwtSecret),
		jwtAccessTTL:  time.Duration(accessTTL) * time.Second,
		jwtRefreshTTL: time.Duration(refreshTTL) * time.Second,
	}
}

func (s *authService) generateAccessToken(userID bson.ObjectID) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.Hex(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtAccessTTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("%w: 비밀번호는 최소 8자여야 합니다", ErrValidation)
	}
	if !passwordRegex.MatchString(password) || !passwordDigitRegex.MatchString(password) {
		return fmt.Errorf("%w: 비밀번호는 영문+숫자 조합이어야 합니다", ErrValidation)
	}
	return nil
}

func (s *authService) issueTokenPair(ctx context.Context, user *User) (*TokenResponse, error) {
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}
	rt := &RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.jwtRefreshTTL),
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, err
	}
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.jwtAccessTTL.Seconds()),
		User:         toUserResponse(user),
	}, nil
}

func (s *authService) Signup(ctx context.Context, req SignupRequest) (*TokenResponse, error) {
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateEmail
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Nickname:     req.Nickname,
		DietaryPrefs: []string{},
		Allergens:    []string{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return s.issueTokenPair(ctx, user)
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokenPair(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	rt, err := s.repo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if rt == nil || time.Now().After(rt.ExpiresAt) {
		return nil, ErrInvalidRefreshToken
	}

	// Delete old token (rotation)
	if err := s.repo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	// Issue new tokens
	accessToken, err := s.generateAccessToken(rt.UserID)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}
	newRT := &RefreshToken{
		UserID:    rt.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(s.jwtRefreshTTL),
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveRefreshToken(ctx, newRT); err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.jwtAccessTTL.Seconds()),
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.DeleteRefreshToken(ctx, refreshToken)
}

func (s *authService) GetMe(ctx context.Context, userID string) (*UserResponse, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrNotFound
	}
	user, err := s.repo.FindUserByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	resp := toUserResponse(user)
	return &resp, nil
}

func (s *authService) UpdateMe(ctx context.Context, userID string, req UpdateMeRequest) (*UserResponse, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrNotFound
	}

	fields := bson.M{}

	if req.Nickname != nil {
		if len(*req.Nickname) < 2 || len(*req.Nickname) > 20 {
			return nil, fmt.Errorf("%w: 닉네임은 2~20자여야 합니다", ErrValidation)
		}
		fields["nickname"] = *req.Nickname
	}
	if req.DietaryPrefs != nil {
		if invalid := validateDietaryPrefs(req.DietaryPrefs); invalid != "" {
			return nil, fmt.Errorf("%w: 허용되지 않는 dietary_prefs 값입니다: %s", ErrValidation, invalid)
		}
		fields["dietary_prefs"] = req.DietaryPrefs
	}
	if req.Allergens != nil {
		if invalid := validateAllergens(req.Allergens); invalid != "" {
			return nil, fmt.Errorf("%w: 허용되지 않는 allergens 값입니다: %s", ErrValidation, invalid)
		}
		fields["allergens"] = req.Allergens
	}

	if len(fields) > 0 {
		if err := s.repo.UpdateUser(ctx, oid, fields); err != nil {
			return nil, err
		}
	}

	user, err := s.repo.FindUserByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	resp := toUserResponse(user)
	return &resp, nil
}

func (s *authService) ValidateAccessToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", ErrInvalidCredentials
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", ErrInvalidCredentials
	}
	return claims.Subject, nil
}
