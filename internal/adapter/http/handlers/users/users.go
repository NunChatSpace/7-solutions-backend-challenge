package users

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/common"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

// InitRoutes initializes API routes
func InitRoutes(router *atreugo.Router, deps *di.Dependency) {
	router.GET("/users", func(rc *atreugo.RequestCtx) error {
		queryParams, err := common.BindQueryToStruct[QueryUsersRequest](rc)
		if err != nil {
			return rc.ErrorResponse(err)
		}
		userService := di.Get[userservices.IUserService](deps)
		users, err := userService.SearchUsers(domain.User{
			Name:  queryParams.Name,
			Email: queryParams.Email,
		})
		if err != nil {
			return rc.ErrorResponse(err)
		}

		return rc.JSONResponse(users)
	})

	router.GET("/users/{id}", func(rc *atreugo.RequestCtx) error {
		params := common.GetParams(rc, []string{"id"})
		userService := di.Get[userservices.IUserService](deps)
		user, err := userService.GetUserByID(params["id"])
		if err != nil {
			return rc.ErrorResponse(err)
		}

		return rc.JSONResponse(user)

	})

	router.POST("/users", func(rc *atreugo.RequestCtx) error {
		reqBody, err := common.BindBodyToStruct[CreateUserRequest](rc)
		if err != nil {
			return rc.ErrorResponse(err)
		}

		validate := validator.New()
		if err := validate.Struct(reqBody); err != nil {
			return rc.ErrorResponse(err)
		}
		userService := di.Get[userservices.IUserService](deps)
		user := domain.User{
			Name:     &reqBody.Name,
			Email:    &reqBody.Email,
			Password: &reqBody.Password,
		}
		if err := userService.CreateUser(&user); err != nil {
			return rc.ErrorResponse(err)
		}

		return rc.JSONResponse(user.ToUserResponse(), fasthttp.StatusCreated)
	})
	router.PATCH("/users/{id}", func(rc *atreugo.RequestCtx) error {
		body, err := common.BindBodyToStruct[UpdateUserRequest](rc)
		if err != nil {
			return rc.ErrorResponse(err)
		}
		params := common.GetParams(rc, []string{"id"})
		userService := di.Get[userservices.IUserService](deps)
		user := domain.User{
			Name:  &body.Name,
			Email: &body.Email,
		}
		if err := userService.UpdateUser(params["id"], &user); err != nil {
			return rc.ErrorResponse(err)
		}

		return rc.JSONResponse(user.ToUserResponse())
	})

	router.DELETE("/users/{id}", func(rc *atreugo.RequestCtx) error {
		params := common.GetParams(rc, []string{"id"})

		userService := di.Get[userservices.IUserService](deps)
		if err := userService.DeleteUser(params["id"]); err != nil {
			return rc.ErrorResponse(err)
		}
		return rc.RawResponse("deleted", 200)
	})
}
