package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository interface {
	EnsureIndexes(ctx context.Context) error
	CreateUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id bson.ObjectID) (*User, error)
	UpdateUser(ctx context.Context, id bson.ObjectID, fields bson.M) error

	SaveRefreshToken(ctx context.Context, token *RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type userRepository struct {
	users         *mongo.Collection
	refreshTokens *mongo.Collection
}

func NewUserRepository(users, refreshTokens *mongo.Collection) UserRepository {
	return &userRepository{users: users, refreshTokens: refreshTokens}
}

func (r *userRepository) EnsureIndexes(ctx context.Context) error {
	// users: email unique index
	_, err := r.users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// refresh_tokens: token unique index
	_, err = r.refreshTokens.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "token", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// refresh_tokens: TTL index on expires_at
	_, err = r.refreshTokens.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	return err
}

func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	result, err := r.users.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserByID(ctx context.Context, id bson.ObjectID) (*User, error) {
	var user User
	err := r.users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id bson.ObjectID, fields bson.M) error {
	fields["updated_at"] = time.Now().UTC()
	_, err := r.users.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": fields},
	)
	return err
}

func (r *userRepository) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
	_, err := r.refreshTokens.InsertOne(ctx, token)
	return err
}

func (r *userRepository) FindRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := r.refreshTokens.FindOne(ctx, bson.M{"token": token}).Decode(&rt)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *userRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.refreshTokens.DeleteOne(ctx, bson.M{"token": token})
	return err
}
