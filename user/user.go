package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Group) {
	app := e.Group("/user")
	app.GET("", findAllUser)
	app.GET("/:id", getUser)
	app.PUT("/:id", updateUser)
	app.DELETE("/:id", deleteUser)
}

func CreateUser(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("INSERT INTO users (username, password, fname, lname, email, tel) VALUES ('%v','%v','%v','%v','%v','%v')", u.Username, u.Password, u.Fname, u.Lname, u.Email, u.Telephone)
	Execute(params)
	return c.JSON(http.StatusOK, "Create user success.")
}

func findAllUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("select id, username, fname, lname, email, tel from users")
	u := Query(params)
	return c.JSON(http.StatusOK, u)
}

func getUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	id, _ := strconv.Atoi(c.Param("id"))
	params["sql"] = fmt.Sprintf("select id, username, fname, lname, email, tel from users where id = %d", id)
	u := Query(params)
	log.Printf("Log get user by id : %v", u)
	return c.JSON(http.StatusOK, u)
}

func updateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("UPDATE users SET username='%v', fname='%v', lname='%v', email='%v', tel='%v' where id = %d", u.Username, u.Fname, u.Lname, u.Email, u.Telephone, id)
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

func Login(c echo.Context) (err error) {
	u := new(UserLogin)
	if err = c.Bind(u); err != nil {
		return
	}

	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("select id, username, password from users where username = '%s' and password = '%s'", u.Username, u.Password)
	result := LoginFunc(params)

	log.Printf("log result login : %v", result)

	if u.Username != result[0].Username {
		return echo.ErrUnauthorized
	}

	// // Createa token
	token := jwt.New(jwt.SigningMethodHS256)

	// // Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = json.Number(strconv.FormatInt(100, 10))
	claims["name"] = "Wachararkon"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}
