package nevergo

import (
	"nevergo/db"
	"nevergo/user"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// StartServices :: desc
func StartServices() {
	e := echo.New()

	// Start DB
	db.ConnDB()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Group prefix routes
	routes := e.Group("/api/v1")

	// ===== Unauthenticate route
	routes.GET("/accessible", user.Accessible)
	routes.POST("/login", user.Login)
	routes.POST("/create", user.CreateUser)

	// ===== Route of documents
	routes.GET("/docs/*", echoSwagger.WrapHandler)

	config := middleware.JWTConfig{
		Claims:     &user.JwtCustomClaims{},
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}
	routes.Use(middleware.JWTWithConfig(config))

	// Routest user controller
	user.UserRoutes(routes)

	e.Logger.Fatal(e.Start(":9001"))
}
