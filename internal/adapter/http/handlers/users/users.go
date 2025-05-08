package users

import "github.com/savsgio/atreugo/v11"

// InitRoutes initializes API routes
func InitRoutes(router *atreugo.Router) {
	router.GET("/users", func(rc *atreugo.RequestCtx) error {
		// implementation for searching users
		return nil
	})
	router.GET("/users/:id", func(rc *atreugo.RequestCtx) error {
		// implementation for getting a user by ID
		return nil
	})
	router.POST("/users", func(rc *atreugo.RequestCtx) error {
		// implementation for creating a user
		return nil
	})
	router.PATCH("/users/:id", func(rc *atreugo.RequestCtx) error {
		// implementation for updating a user by ID
		return nil
	})
	router.DELETE("/users/:id", func(rc *atreugo.RequestCtx) error {
		// implementation for deleting a user by ID
		return nil
	})
}
