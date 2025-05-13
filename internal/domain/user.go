package domain

import "time"

type User struct {
	ID        *string                 `json:"id" bson:"_id,omitempty"`
	Name      *string                 `json:"name" bson:"name"`
	Email     *string                 `json:"email"  bson:"email"`
	Password  *string                 `bson:"password"`
	CreatedAt *time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time              `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time              `json:"deleted_at" bson:"deleted_at"`
	Scopes    *map[string]interface{} `json:"scopes" bson:"scopes"`
}

func (user User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		Scopes:    user.Scopes,
	}
}

type UserResponse struct {
	ID        *string                 `json:"id" bson:"_id,omitempty"`
	Name      *string                 `json:"name" bson:"name"`
	Email     *string                 `json:"email"  bson:"email"`
	CreatedAt *time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time              `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time              `json:"deleted_at" bson:"deleted_at"`
	Scopes    *map[string]interface{} `json:"scopes" bson:"scopes"`
}
