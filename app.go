package nevergo

import (
	"nevergo/db"
	"nevergo/user"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func StartServices() {
	e := echo.New()

	db.ConnDB()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ===== Initial resource from MongoDB
	v1 := e.Group("/api/v1")
	user.UserRoutes(v1)

	e.Logger.Fatal(e.Start(":9001"))
}
