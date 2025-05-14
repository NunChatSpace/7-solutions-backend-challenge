package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	repo := UserRepository{db: db}
	if err := repo.initIndexes(); err != nil {
		log.Fatal("could not create indexes:", err)
	}
	return &repo
}

func (u *UserRepository) initIndexes() error {
	collection := u.db.Collection(MONGO_USER_COLLECTION)

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"deleted_at": bson.M{"$eq": nil},
			}),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return fmt.Errorf("could not create partial unique index: %w", err)
	}
	return nil
}

func (u *UserRepository) InsertUser(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Hash the password
	if user.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		hashedStr := string(hashed)
		user.Password = &hashedStr
	}

	// 2. Set placeholder scopes
	scope := map[string]interface{}{
		"users": 15,
	}
	user.Scopes = &scope

	// 3. Insert into MongoDB
	now := time.Now()
	user.CreatedAt = &now

	res, err := u.db.Collection(MONGO_USER_COLLECTION).InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("error: %s was used", *user.Email)
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// ✅ 4. Assign inserted ID back to user.ID
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		idStr := oid.Hex()
		user.ID = &idStr
	} else {
		return fmt.Errorf("inserted ID is not an ObjectID")
	}

	return nil
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

func (u *UserRepository) UpdateUser(id string, user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %w", err)
	}

	filter := bson.M{"_id": oid, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{
		"$set": bson.M{
			"email": user.Email,
			"name":  user.Name,
		},
	}
	now := time.Now()
	user.UpdatedAt = &now
	res, err := u.db.Collection(MONGO_USER_COLLECTION).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	var updatedUser domain.User
	err = u.db.Collection(MONGO_USER_COLLECTION).FindOne(ctx, filter).Decode(&updatedUser)
	if err != nil {
		return fmt.Errorf("failed to fetch updated user: %w", err)
	}
	user = &updatedUser

	return nil
}

func (u *UserRepository) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %w", err)
	}

	filter := bson.M{"_id": oid, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	// You might not even need this GetUserByID call unless you want to validate or log
	if _, err := u.db.Collection(MONGO_USER_COLLECTION).UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("failed to update deleted_at: %w", err)
	}

	return nil
}
