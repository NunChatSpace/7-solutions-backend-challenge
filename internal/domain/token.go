package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenInfo struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	Expired   time.Time `json:"expired"`
	CreatedAt time.Time `json:"created_at"`
	Type      string    `json:"type"`
	// int of scope is CRUD translated as binary 4 bits
	// C = 1
	// R = 2
	// U = 4
	// D = 8
	Scopes map[string]interface{} `json:"scopes"`
}

func (t *TokenInfo) ToJWTClaims() jwt.MapClaims {
	claims := jwt.MapClaims{}
	claims["user_id"] = t.UserID
	claims["session_id"] = t.SessionID
	claims["type"] = t.Type
	claims["scopes"] = t.Scopes
	claims["exp"] = jwt.NewNumericDate(t.Expired)
	claims["iat"] = jwt.NewNumericDate(time.Now())

	return claims
}

func (t *TokenInfo) FromJWTClaims(claims jwt.MapClaims) {
	t.UserID = claims["user_id"].(string)
	t.SessionID = claims["session_id"].(string)
	t.Type = claims["type"].(string)
	t.Scopes = claims["scopes"].(map[string]interface{})
	iatRaw, ok := claims["iat"].(float64)
	if !ok {
		iatRaw = 0
	}
	t.CreatedAt = time.Unix(int64(iatRaw), 0)

	expRaw, ok := claims["exp"].(float64)
	if !ok {
		expRaw = 0
	}
	t.Expired = time.Unix(int64(expRaw), 0)
}
