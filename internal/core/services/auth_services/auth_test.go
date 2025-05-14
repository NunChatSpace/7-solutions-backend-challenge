package authservices_test

import (
	"os"
	"testing"
	"time"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	testutils "github.com/NunChatSpace/7-solutions-backend-challenge/internal/test_utils"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestAuthService_TokenFlow(t *testing.T) {
	deps := testutils.NewTestDependency(cfg)
	service := authservices.NewAuthService(deps)
	t.Run("generate access and refresh tokens", func(t *testing.T) {
		info := domain.TokenInfo{
			UserID:  "u-123",
			Scopes:  map[string]interface{}{"users": 1},
			Expired: time.Now().Add(1 * time.Hour),
		}
		accessToken, refreshToken, err := service.GenerateTokens(info)
		if err != nil {
			t.Fatalf("unexpected error generating tokens: %v", err)
		}
		if accessToken == "" || refreshToken == "" {
			t.Fatal("expected non-empty access and refresh tokens")
		}

		accessDecoded, err := service.DecodeToken(accessToken)
		if err != nil {
			t.Fatalf("failed to decode access token: %v", err)
		}
		if accessDecoded.Type != "access_token" {
			t.Errorf("expected type 'access_token', got: %s", accessDecoded.Type)
		}
		if accessDecoded.UserID != info.UserID {
			t.Errorf("expected userID %s, got %s", info.UserID, accessDecoded.UserID)
		}

		refreshDecoded, err := service.DecodeToken(refreshToken)
		if err != nil {
			t.Fatalf("failed to decode refresh token: %v", err)
		}
		if refreshDecoded.Type != "refresh_token" {
			t.Errorf("expected type 'refresh_token', got: %s", refreshDecoded.Type)
		}
		if refreshDecoded.UserID != info.UserID {
			t.Errorf("expected userID %s, got %s", info.UserID, refreshDecoded.UserID)
		}
	})
	t.Run("generate and validate access token", func(t *testing.T) {
		info := domain.TokenInfo{
			UserID:  "u-123",
			Type:    "access_token",
			Scopes:  map[string]interface{}{"users": 1},
			Expired: time.Now().Add(1 * time.Hour),
		}
		accessToken, err := service.GenerateAccessToken(info)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		decoded, err := service.DecodeToken(accessToken)
		if err != nil {
			t.Fatalf("failed to decode token: %v", err)
		}
		if decoded.UserID != info.UserID {
			t.Errorf("expected userID %s, got %s", info.UserID, decoded.UserID)
		}
	})

	t.Run("generate and validate refresh token", func(t *testing.T) {
		info := domain.TokenInfo{
			UserID:  "u-123",
			Type:    "refresh_token",
			Scopes:  map[string]interface{}{"users": 1},
			Expired: time.Now().Add(1 * time.Hour),
		}
		refreshToken, err := service.GenerateRefreshToken(info)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		decoded, err := service.DecodeToken(refreshToken)
		if err != nil {
			t.Fatalf("failed to decode token: %v", err)
		}
		if decoded.Type != "refresh_token" {
			t.Errorf("expected type 'refresh_token', got: %s", decoded.Type)
		}
	})

	t.Run("token is expired", func(t *testing.T) {
		info := domain.TokenInfo{
			UserID:  "u-123",
			Type:    "access_token",
			Scopes:  map[string]interface{}{},
			Expired: time.Now().Add(-1 * time.Minute),
		}

		token, err := service.EncodeToken(info)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		decoded, err := service.DecodeToken(token)
		if err == nil || decoded != nil {
			t.Fatal("expected token validation to fail due to expiration")
		}
	})

	t.Run("token type mismatch", func(t *testing.T) {
		info := domain.TokenInfo{
			UserID:  "u-123",
			Type:    "wrong_type",
			Scopes:  map[string]interface{}{},
			Expired: time.Now().Add(1 * time.Hour),
		}
		token, err := service.EncodeToken(info)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = service.DecodeToken(token)
		if err == nil {
			t.Fatal("expected token validation to fail due to type mismatch")
		}
	})

	t.Run("EncodeToken sets default expiry if zero", func(t *testing.T) {
		info := domain.TokenInfo{
			UserID: "u-123",
			Type:   "access_token",
			Scopes: map[string]interface{}{"users": 1},
			// Expired is zero
		}

		token, err := service.EncodeToken(info)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		decoded, err := service.DecodeToken(token)
		if err != nil {
			t.Fatalf("failed to decode token: %v", err)
		}

		// Expect expiry to be around 24 hours from now
		if decoded.Expired.IsZero() {
			t.Fatal("expected non-zero expiry time")
		}
	})
}
