package authservices

import (
	"errors"
	"fmt"
	"time"

	dbrepo "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type Port interface {
	GenerateTokens(tokenInfo domain.TokenInfo) (string, string, error)
	GenerateRefreshToken(tokenInfo domain.TokenInfo) (string, error)
	ValidateRefreshToken(domain.TokenInfo) error
	GenerateAccessToken(tokenInfo domain.TokenInfo) (string, error)
	ValidateAccessToken(domain.TokenInfo) error

	DecodeToken(token string) (*domain.TokenInfo, error)
	EncodeToken(tokenInfo domain.TokenInfo) (string, error)

	IsPermit(tokenStr string) (*domain.User, error)
}

type authService struct {
	Repository dbrepo.Repository
	Config     *config.Config
}

func NewAuthService(repo dbrepo.Repository, cfg *config.Config) Port {
	return &authService{
		Repository: repo,
		Config:     cfg,
	}
}
func (s *authService) GenerateTokens(tokenInfo domain.TokenInfo) (string, string, error) {
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

func (s *authService) GenerateRefreshToken(tokenInfo domain.TokenInfo) (string, error) {
	return s.EncodeToken(tokenInfo)
}
func (s *authService) ValidateRefreshToken(tokenInfo domain.TokenInfo) error {
	return s.validateToken(tokenInfo, "refresh_token")
}
func (s *authService) GenerateAccessToken(tokenInfo domain.TokenInfo) (string, error) {
	token, err := s.EncodeToken(tokenInfo)
	return token, err
}
func (s *authService) ValidateAccessToken(tokenInfo domain.TokenInfo) error {
	return s.validateToken(tokenInfo, "access_token")
}

func (s *authService) DecodeToken(tokenStr string) (*domain.TokenInfo, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.Config.JWT.SecretKey), nil
	})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return nil, err
	}

	// Extract and return claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	fmt.Printf("Claims: %v\n", claims)

	tokenInfo := domain.TokenInfo{}
	tokenInfo.FromJWTClaims(claims)
	if err := s.ValidateAccessToken(tokenInfo); err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

func (s *authService) EncodeToken(tokenInfo domain.TokenInfo) (string, error) {
	tokenInfo.Expired = time.Now().Add(24 * time.Hour)
	claims := tokenInfo.ToJWTClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with secret key
	signedToken, err := token.SignedString([]byte(s.Config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *authService) validateToken(tokenInfo domain.TokenInfo, kind string) error {
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

func (s *authService) IsPermit(tokenStr string) (*domain.User, error) {

	// tokenInfo, err := dep.Services.Auth().DecodeToken(tokenStr)
	// if err != nil {
	// 	return ctx.ErrorResponse(err)
	// }

	// user, err := dep.Services.User().GetUserByID(tokenInfo.UserID)
	// if err != nil {
	// 	return ctx.ErrorResponse(err)
	// }

	// parts := strings.Split(fullPath, "/")
	// if len(parts) < 4 {
	// 	return errors.New("invalid path")
	// }
	// // "/api/v1/users" → ["", "api", "v1", "users"]
	// rootPath := parts[3]
	// methodNumber := 0
	// if string(ctx.Request.Header.Method()) == "POST" {
	// 	methodNumber = 1
	// } else if string(ctx.Request.Header.Method()) == "GET" {
	// 	methodNumber = 2
	// } else if string(ctx.Request.Header.Method()) == "PATCH" {
	// 	methodNumber = 3
	// } else if string(ctx.Request.Header.Method()) == "DELETE" {
	// 	methodNumber = 4
	// } else if string(ctx.Request.Header.Method()) == "OPTIONS" {
	// 	methodNumber = 0
	// } else {
	// 	return errors.New("invalid method")
	// }

	// userScope := (*user.Scopes)[rootPath]
	// if methodNumber > 0 && methodNumber&userScope == 0 {
	// 	return errors.New("user does not have permission to access this resource")
	// }
	return nil, nil
}
