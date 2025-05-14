package authservices

import (
	"errors"
	"fmt"
	"time"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type IAuthSerivce interface {
	GenerateTokens(tokenInfo domain.TokenInfo) (string, string, error)
	GenerateRefreshToken(tokenInfo domain.TokenInfo) (string, error)
	ValidateRefreshToken(domain.TokenInfo) error
	GenerateAccessToken(tokenInfo domain.TokenInfo) (string, error)
	ValidateAccessToken(domain.TokenInfo) error

	DecodeToken(token string) (*domain.TokenInfo, error)
	EncodeToken(tokenInfo domain.TokenInfo) (string, error)
}

type authService struct {
	Dependencies *di.Dependency

	cfg *config.Config
}

func NewAuthService(deps *di.Dependency) IAuthSerivce {
	return authService{
		Dependencies: deps,
		cfg:          di.Get[*config.Config](deps),
	}
}
func (s authService) GenerateTokens(tokenInfo domain.TokenInfo) (string, string, error) {
	// Implementation for generating access and refresh tokens
	// This is a placeholder implementation
	tokenInfo.Type = "access_token"
	accessToken, err := s.GenerateAccessToken(tokenInfo)
	if err != nil {
		return "", "", err
	}
	tokenInfo.Type = "refresh_token"
	refreshToken, err := s.GenerateRefreshToken(tokenInfo)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s authService) GenerateRefreshToken(tokenInfo domain.TokenInfo) (string, error) {
	return s.EncodeToken(tokenInfo)
}
func (s authService) ValidateRefreshToken(tokenInfo domain.TokenInfo) error {
	return s.validateToken(tokenInfo, "refresh_token")
}
func (s authService) GenerateAccessToken(tokenInfo domain.TokenInfo) (string, error) {
	token, err := s.EncodeToken(tokenInfo)
	return token, err
}
func (s authService) ValidateAccessToken(tokenInfo domain.TokenInfo) error {
	return s.validateToken(tokenInfo, "access_token")
}

func (s authService) DecodeToken(tokenStr string) (*domain.TokenInfo, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.cfg.JWT.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Extract and return claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	tokenInfo := domain.TokenInfo{}
	tokenInfo.FromJWTClaims(claims)
	if tokenInfo.Type == "access_token" {
		if err := s.ValidateAccessToken(tokenInfo); err != nil {
			return nil, err
		}
	} else if tokenInfo.Type == "refresh_token" {
		if err := s.ValidateRefreshToken(tokenInfo); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unexpected token type: %v", tokenInfo.Type)
	}

	return &tokenInfo, nil
}

func (s authService) EncodeToken(tokenInfo domain.TokenInfo) (string, error) {
	if tokenInfo.Expired.IsZero() {
		tokenInfo.Expired = time.Now().Add(24 * time.Hour)
	}
	claims := tokenInfo.ToJWTClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with secret key
	signedToken, err := token.SignedString([]byte(s.cfg.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s authService) validateToken(tokenInfo domain.TokenInfo, kind string) error {
	// Check type
	if tokenInfo.Type != kind {
		return fmt.Errorf("unexpected token type: %v", tokenInfo.Type)
	}

	// Check expiry
	if tokenInfo.Expired.Before(time.Now()) {
		return errors.New("token is expired")
	}

	return nil
}
