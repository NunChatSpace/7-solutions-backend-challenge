package domain

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenInfo struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	ExpiredAt string `json:"expired_at"`
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"`
	// int of scope is CRUD translated as binary 4 bits
	// C = 1
	// R = 2
	// U = 4
	// D = 8
	Scopes []map[string]int `json:"scopes"`
}
