package sessions

import (
	"strings"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/middlewares/logger"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/common"
	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	sessionservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/session_services"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"github.com/savsgio/atreugo/v11"
)

func InitRoutes(router *atreugo.Router, deps *di.Dependency) {
	router.POST("/sessions", func(rc *atreugo.RequestCtx) error {
		body, err := common.BindBodyToStruct[loginRequest](rc)
		if err != nil {
			return rc.ErrorResponse(err)
		}

		userService := di.Get[userservices.IUserService](deps)
		sessionService := di.Get[sessionservices.ISessionService](deps)
		user := domain.User{
			Email:    &body.Email,
			Password: &body.Password,
		}
		if err := userService.Authenticate(&user); err != nil {
			return rc.ErrorResponse(err)
		}
		session, err := sessionService.CreateSession(*user.ID)
		if err != nil {
			return rc.ErrorResponse(err)
		}

		return rc.JSONResponse(session)
	})

	router.POST("/sessions/refresh", func(rc *atreugo.RequestCtx) error {
		body, err := common.BindBodyToStruct[refreshTokenRequest](rc)
		if err != nil {
			return rc.ErrorResponse(err)
		}

		authService := di.Get[authservices.IAuthSerivce](deps)
		token, err := authService.DecodeToken(body.RefreshToken)
		if err != nil {
			return rc.ErrorResponse(err)
		}
		tokenInfo := domain.TokenInfo{
			UserID:    token.UserID,
			SessionID: token.SessionID,
			Type:      "access_token",
		}
		accessToken, refreshToken, err := authService.GenerateTokens(tokenInfo)
		if err != nil {
			return rc.ErrorResponse(err)
		}
		return rc.JSONResponse(domain.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	})

	router.DELETE("/sessions", func(rc *atreugo.RequestCtx) error {
		logger := logger.FromContext(rc, "DELETE /sessions")
		authHeader := string(rc.Request.Header.Peek("Authorization"))
		if authHeader == "" {
			logger.Warn("Authorization header is missing")
			return rc.RawResponse("session terminated", 204)
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Warn("Authorization header format is invalid")
			return rc.RawResponse("session terminated", 204)
		}

		token := parts[1]
		sessionService := di.Get[sessionservices.ISessionService](deps)
		authService := di.Get[authservices.IAuthSerivce](deps)
		tokenInfo, _ := authService.DecodeToken(token)
		sessionService.TerminateSession(tokenInfo.SessionID)

		return rc.RawResponse("session terminated", 204)
	})
}
