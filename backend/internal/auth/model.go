package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email"         json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Nickname     string             `bson:"nickname"      json:"nickname"`
	DietaryPrefs []string           `bson:"dietary_prefs" json:"dietaryPrefs"`
	Allergens    []string           `bson:"allergens"     json:"allergens"`
	CreatedAt    time.Time          `bson:"created_at"    json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at"    json:"updatedAt"`
}

type RefreshToken struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"user_id"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expires_at"`
	CreatedAt time.Time          `bson:"created_at"`
}

type SignupRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Nickname string `json:"nickname" binding:"required,min=2,max=20"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UpdateMeRequest struct {
	Nickname     *string  `json:"nickname"`
	DietaryPrefs []string `json:"dietaryPrefs"`
	Allergens    []string `json:"allergens"`
}

type UserResponse struct {
	ID           bson.ObjectID `json:"id"`
	Email        string             `json:"email"`
	Nickname     string             `json:"nickname"`
	DietaryPrefs []string           `json:"dietaryPrefs"`
	Allergens    []string           `json:"allergens"`
	CreatedAt    time.Time          `json:"createdAt"`
}
