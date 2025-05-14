package authen

import (
	"errors"
	"strings"

	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"github.com/savsgio/atreugo/v11"
)

type ctxKey struct{}

var EXCEPT_APIS = map[string]map[string]bool{
	"/api/v1/users": {
		"POST": true,
	},
	"/api/v1/sessions": {
		"POST": true,
	},
}

func FromContext(ctx *atreugo.RequestCtx) *domain.User {
	if ctx.Value(&ctxKey{}) == nil {
		return nil
	}
	user, ok := ctx.Value(&ctxKey{}).(*domain.User)
	if !ok {
		return nil
	}
	return user
}

func Handler(ctx *atreugo.RequestCtx, deps *di.Dependency) error {
	fullPath := string(ctx.Request.URI().Path())
	httpMethod := string(ctx.Request.Header.Method())
	if !isReqiredAuth(fullPath, string(httpMethod)) {
		return nil
	}
	authHeader := string(ctx.Request.Header.Peek("Authorization"))
	if authHeader == "" {
		return errors.New("invalid token")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	tokenStr = strings.TrimSpace(tokenStr)
	if tokenStr == "" {
		return errors.New("invalid token")
	}
	authService := di.Get[authservices.IAuthSerivce](deps)
	tokenInfo, err := authService.DecodeToken(tokenStr)
	if err != nil {
		return ctx.ErrorResponse(err)
	}

	userService := di.Get[userservices.IUserService](deps)
	user, err := userService.GetUserByID(tokenInfo.UserID)
	if err != nil {
		return ctx.ErrorResponse(err)
	}

	parts := strings.Split(fullPath, "/")
	if len(parts) < 4 {
		return errors.New("invalid path")
	}
	// "/api/v1/users" → ["", "api", "v1", "users"]
	rootPath := parts[3]
	var methodNumber int32 = 0
	if string(httpMethod) == "POST" {
		methodNumber = 1
	} else if string(httpMethod) == "GET" {
		methodNumber = 2
	} else if string(httpMethod) == "PATCH" {
		methodNumber = 3
	} else if string(httpMethod) == "DELETE" {
		methodNumber = 4
	} else if string(httpMethod) == "OPTIONS" {
		methodNumber = 0
	} else {
		return errors.New("invalid method")
	}

	userScope := ((*user.Scopes)[rootPath]).(int32)
	if methodNumber > 0 && methodNumber&userScope == 0 {
		return errors.New("user does not have permission to access this resource")
	}

	di.Provide(deps, domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Scopes:    user.Scopes,
	})

	return ctx.Next()
}

func isReqiredAuth(fullPath string, method string) bool {
	if _, ok := EXCEPT_APIS[fullPath]; ok {
		if _, ok := EXCEPT_APIS[fullPath][method]; ok {
			return false
		}
	}
	return true
}
