package user

import (
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Group) {
	app := e.Group("/user")
	app.GET("", findAllUser)
	app.GET("/me", getUser)
	app.PUT("/me", updateUser)
	app.PATCH("/resetpassword/:id", resetPassword)
}
