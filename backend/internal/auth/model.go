package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// --- Entities ---

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string        `bson:"email"         json:"email"`
	PasswordHash string        `bson:"password_hash" json:"-"`
	Nickname     string        `bson:"nickname"      json:"nickname"`
	DietaryPrefs []string      `bson:"dietary_prefs" json:"dietary_prefs"`
	Allergens    []string      `bson:"allergens"     json:"allergens"`
	CreatedAt    time.Time     `bson:"created_at"    json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"    json:"updated_at"`
}

type RefreshToken struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"user_id"`
	Token     string        `bson:"token"`
	ExpiresAt time.Time     `bson:"expires_at"`
	CreatedAt time.Time     `bson:"created_at"`
}

// --- Requests ---

type SignupRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Nickname string `json:"nickname" binding:"required,min=2,max=20"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateMeRequest struct {
	Nickname     *string  `json:"nickname"`
	DietaryPrefs []string `json:"dietary_prefs"`
	Allergens    []string `json:"allergens"`
}

// --- Responses ---

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserResponse struct {
	ID           bson.ObjectID `json:"id"`
	Email        string        `json:"email"`
	Nickname     string        `json:"nickname"`
	DietaryPrefs []string      `json:"dietary_prefs"`
	Allergens    []string      `json:"allergens"`
	CreatedAt    time.Time     `json:"created_at"`
}

func toUserResponse(u *User) UserResponse {
	prefs := u.DietaryPrefs
	if prefs == nil {
		prefs = []string{}
	}
	allergens := u.Allergens
	if allergens == nil {
		allergens = []string{}
	}
	return UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		Nickname:     u.Nickname,
		DietaryPrefs: prefs,
		Allergens:    allergens,
		CreatedAt:    u.CreatedAt,
	}
}

// --- Enum validation ---

var allowedDietaryPrefs = map[string]bool{
	"vegetarian": true,
	"vegan":      true,
	"halal":      true,
	"kosher":     true,
	"gluten-free": true,
	"dairy-free": true,
	"low-sodium": true,
	"low-carb":   true,
}

var allowedAllergens = map[string]bool{
	"gluten":    true,
	"peanut":    true,
	"tree-nut":  true,
	"milk":      true,
	"egg":       true,
	"soy":       true,
	"fish":      true,
	"shellfish": true,
	"sesame":    true,
	"wheat":     true,
}

func validateDietaryPrefs(prefs []string) string {
	for _, p := range prefs {
		if !allowedDietaryPrefs[p] {
			return p
		}
	}
	return ""
}

func validateAllergens(allergens []string) string {
	for _, a := range allergens {
		if !allowedAllergens[a] {
			return a
		}
	}
	return ""
}
