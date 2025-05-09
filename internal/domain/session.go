package domain

type Session struct {
	ID        *string `json:"id" bson:"_id,omitempty"`
	UserID    *string `json:"user_id" bson:"user_id"`
	CreatedAt *string `json:"created_at" bson:"created_at"`
	ExpiredAt *string `json:"expired_at" bson:"expired_at"`
	DeletedAt *string `json:"deleted_at" bson:"deleted_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
