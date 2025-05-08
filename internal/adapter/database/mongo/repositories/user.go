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

	res, err := u.db.Collection(MONGO_USER_COLLECTION).InsertOne(ctx, user)
	fmt.Println("Inserted user with ID:", res.InsertedID)
	return err
}

func (u *UserRepository) GetUserByID(id string) (*domain.User, error) {
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
	var user domain.User
	if err := cur.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) Search(user domain.User) ([]*domain.User, error) {
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
