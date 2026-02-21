package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id bson.ObjectID) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id bson.ObjectID, update *User) error

	SaveRefreshToken(ctx context.Context, token *RefreshToken) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteAllRefreshTokens(ctx context.Context, userID bson.ObjectID) error
}

type userRepository struct {
	users         *mongo.Collection
	refreshTokens *mongo.Collection
}

func NewUserRepository(users, refreshTokens *mongo.Collection) UserRepository {
	return &userRepository{users: users, refreshTokens: refreshTokens}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (r *userRepository) FindByID(ctx context.Context, id bson.ObjectID) (*User, error) {
	return nil, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	return nil
}

func (r *userRepository) Update(ctx context.Context, id bson.ObjectID, update *User) error {
	return nil
}

func (r *userRepository) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
	return nil
}

func (r *userRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	return nil, nil
}

func (r *userRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	return nil
}

func (r *userRepository) DeleteAllRefreshTokens(ctx context.Context, userID bson.ObjectID) error {
	return nil
}
