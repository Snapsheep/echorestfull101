package user

import (
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Group) {
	app := e.Group("/user")
	app.GET("", findAllUser)
	app.GET("/:id", getUser)
	app.PATCH("/:id", updateUser)
	app.DELETE("/:id", deleteUser)
}
