package nevergo

import (
	"nevergo/db"
	"nevergo/user"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// StartServices :: desc
func StartServices() {
	e := echo.New()

	// Start DB
	db.ConnDB()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Group prefix routes
	v1 := e.Group("/api/v1")

	// Routest user controller
	user.UserRoutes(v1)

	e.Logger.Fatal(e.Start(":9001"))
}
