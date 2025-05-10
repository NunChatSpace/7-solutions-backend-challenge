package domain

type User struct {
	ID        *string                 `json:"id" bson:"_id,omitempty"`
	Name      *string                 `json:"name" bson:"name"`
	Email     *string                 `json:"email"  bson:"email"`
	Password  *string                 `bson:"password"`
	CreatedAt *string                 `json:"created_at" bson:"created_at"`
	UpdatedAt *string                 `json:"updated_at" bson:"updated_at"`
	DeletedAt *string                 `json:"deleted_at" bson:"deleted_at"`
	Scopes    *map[string]interface{} `json:"scopes" bson:"scopes"`
}

type UserResponse struct {
	ID        *string                 `json:"id" bson:"_id,omitempty"`
	Name      *string                 `json:"name" bson:"name"`
	Email     *string                 `json:"email"  bson:"email"`
	CreatedAt *string                 `json:"created_at" bson:"created_at"`
	UpdatedAt *string                 `json:"updated_at" bson:"updated_at"`
	DeletedAt *string                 `json:"deleted_at" bson:"deleted_at"`
	Scopes    *map[string]interface{} `json:"scopes" bson:"scopes"`
}
