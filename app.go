package nevergo

import (
	"fmt"
	"nevergo/db"
	"nevergo/user"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// StartServices :: desc
func StartServices() {
	e := echo.New()
	v := validator.New()

	a := user.User{
		Email:    "something",
		Username: "A girl has no name",
		Password: "1234",
	}

	err := v.Struct(a)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e)
	}

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
	routes.POST("/user/create", user.CreateUser)

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
