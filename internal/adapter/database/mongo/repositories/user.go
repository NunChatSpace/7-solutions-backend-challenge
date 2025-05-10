package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const MONGO_USER_COLLECTION = "users"

func (r *RepositoryImpl) User() database.IUserRepository {
	return r.userRepo
}

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) database.IUserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) InsertUser(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 1. Check if the email already exists
	filter := bson.M{"email": user.Email}
	count, err := u.db.Collection(MONGO_USER_COLLECTION).CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to check for existing email: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("email already exists")
	}

	// 2. Hash the password
	if user.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		hashedStr := string(hashed)
		user.Password = &hashedStr
	}

	// Just for authorization testing purposes
	// In a real application, you would not set scopes like this
	// and would instead use a proper authorization system
	// and set scopes based on user roles or permissions
	// This is just a placeholder for the sake of example
	scope := map[string]interface{}{
		"users": 15,
	}
	user.Scopes = &scope
	res, err := u.db.Collection(MONGO_USER_COLLECTION).InsertOne(ctx, user)
	fmt.Println("Inserted user with ID:", res.InsertedID)
	return err
}

func (u *UserRepository) GetUserByID(id string) (*domain.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}
	filter := bson.M{"_id": oid, "deleted_at": bson.M{"$eq": nil}}

	cur := u.db.Collection(MONGO_USER_COLLECTION).FindOne(ctx, filter)
	if err := cur.Err(); err != nil {
		return nil, err
	}
	var user domain.UserResponse
	if err := cur.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) Search(user domain.User) ([]*domain.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"deleted_at": bson.M{"$eq": nil},
	}

	if user.Email != nil {
		filter["email"] = *user.Email
	}
	if user.Name != nil {
		filter["status"] = *user.Name
	}

	cur, err := u.db.Collection(MONGO_USER_COLLECTION).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []*domain.UserResponse
	for cur.Next(ctx) {
		var usr domain.UserResponse
		if err := cur.Decode(&usr); err != nil {
			return nil, err
		}
		users = append(users, &usr)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepository) SearchForAuth(user domain.User) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"deleted_at": bson.M{"$eq": nil},
	}

	if user.Email != nil {
		filter["email"] = *user.Email
	}
	if user.Name != nil {
		filter["status"] = *user.Name
	}

	cur, err := u.db.Collection(MONGO_USER_COLLECTION).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []*domain.User
	for cur.Next(ctx) {
		var usr domain.User
		if err := cur.Decode(&usr); err != nil {
			return nil, err
		}
		users = append(users, &usr)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepository) UpdateUser(user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": user.ID, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{
		"$set": bson.M{
			"email":    user.Email,
			"password": user.Password,
		},
	}

	res, err := u.db.Collection(MONGO_USER_COLLECTION).UpdateOne(ctx, filter, update)
	if err != nil {
		return user, err
	}

	if res.MatchedCount == 0 {
		return user, fmt.Errorf("user not found")
	}

	return user, nil
}

func (u *UserRepository) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("delete fail: user not found")
	}

	if _, err := u.db.Collection(MONGO_USER_COLLECTION).UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}
