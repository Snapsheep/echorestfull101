package nevergo

import (
	"nevergo/db"
	"nevergo/user"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// StartServices :: desc
func StartServices() {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}

	// Start DB
	db.ConnDB()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Group prefix routes
	routes := e.Group("/api/v1")

	// ===== Unauthenticate route
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
