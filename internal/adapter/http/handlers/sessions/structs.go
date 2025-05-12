package sessions

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
