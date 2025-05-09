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

const MONGO_SESSION_COLLECTION = "sessions"

func (r *RepositoryImpl) Session() database.ISessionRepository {
	return r.sessionRepo
}

type SessionRepository struct {
	db *mongo.Database
}

func NewSessionRepository(db *mongo.Database) database.ISessionRepository {
	return &SessionRepository{db}
}

func (s *SessionRepository) InsertSession(session *domain.Session) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.db.Collection(MONGO_SESSION_COLLECTION).InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}
	_h := insertedID.Hex()
	session.ID = &_h

	return session, nil
}

func (s *SessionRepository) GetSessionByID(id string) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the string ID to an ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	filter := bson.M{"_id": oid, "deleted_at": bson.M{"$eq": nil}}

	cur := s.db.Collection(MONGO_SESSION_COLLECTION).FindOne(ctx, filter)
	if err := cur.Err(); err != nil {
		return nil, err
	}
	var session domain.Session
	if err := cur.Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionRepository) TerminateSession(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the string ID to an ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id")
	}

	filter := bson.M{"_id": oid}
	_, err = s.db.Collection(MONGO_SESSION_COLLECTION).DeleteOne(ctx, filter)
	return err
}
