package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func UserRoutes(e *echo.Group) {
	app := e.Group("/user")
	app.GET("", findAllUser)
	app.GET("/:id", getUser)
	app.POST("/create", createUser)
	app.PUT("/:id", updateUser)
	app.DELETE("/:id", deleteUser)
}

func createUser(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("INSERT INTO users (name,email) VALUES ('%v','%v')", u.Name, u.Email)
	Execute(params)
	return c.JSON(http.StatusOK, u)
}

func findAllUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("select id, name, email from users")
	u := Query(params)
	return c.JSON(http.StatusOK, u)
}

func getUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	id, _ := strconv.Atoi(c.Param("id"))
	params["sql"] = fmt.Sprintf("select id, name, email from users where id = %d", id)
	u := Query(params)
	return c.JSON(http.StatusOK, u)
}

func updateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("UPDATE users SET name='%v', email='%v' where id = %d", u.Name, u.Email, id)
	Execute(params)
	return c.JSON(http.StatusOK, u)
}

func deleteUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	id, _ := strconv.Atoi(c.Param("id"))
	params["sql"] = fmt.Sprintf("delete from users where id = %d", id)
	Execute(params)
	return c.NoContent(http.StatusNoContent)
}
