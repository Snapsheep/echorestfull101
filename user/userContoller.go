package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nevergo/db"
	middleware "nevergo/middlewares"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type resetPass struct {
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

// @Summary Retrieves user based on query
// @Description Post User
// @Accept  json
// @Produce json
// @Success 200 {array} User
// @Router /api/v1/user/create [post]
func CreateUser(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}

	if len(u.Email) > 0 && !middleware.IsEmailValid(u.Email) {
		fmt.Println(u.Email + " is not a valid email")
		c.JSON(http.StatusBadRequest, u.Email+" is no a valid email")
		return
	}

	if !middleware.IsEmpty(u.Username) && !middleware.IsEmpty(u.Password) {
		fmt.Println("username, password is not empty")
		c.JSON(http.StatusBadRequest, "username, password require field")
		return
	}

	sqlStatement := `SELECT count(id) FROM users WHERE username=$1`
	row := db.SqliteHandler.Conn.QueryRow(sqlStatement, u.Username)
	var (
		numID int
	)
	err = row.Scan(&numID)
	if err != nil {
		fmt.Println(err)
	}

	if numID > 0 {
		c.JSON(http.StatusBadRequest, "username "+u.Username+" is duplicate!")
		return
	}

	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("INSERT INTO users (username, password, fname, lname, email, tel) VALUES ('%v','%v','%v','%v','%v','%v')", u.Username, middleware.HashAndSalt([]byte(u.Password)), u.Fname, u.Lname, u.Email, u.Telephone)
	execute(params)
	return c.JSON(http.StatusOK, "Create user success.")
}

// @Summary Retrieves user based on query
// @Description Get User
// @Accept  json
// @Produce json
// @Success 200 {array} User
// @Router /api/v1/user/me [get]
func getUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	userID := claims.ID

	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("SELECT id, username, fname, lname, email, tel FROM users WHERE id = %d", userID)
	u := query(params)
	return c.JSON(http.StatusOK, u)
}

// @Summary Retrieves user based on query
// @Description Put User
// @Accept  json
// @Produce json
// @Success 200 {array} User
// @Router /api/v1/user/me [put]
func updateUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	userID := claims.ID

	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("UPDATE users SET username='%v', fname='%v', lname='%v', email='%v', tel='%v' WHERE id = %d", u.Username, u.Fname, u.Lname, u.Email, u.Telephone, userID)
	execute(params)
	return c.JSON(http.StatusOK, u)
}

// @Summary Retrieves user based on query
// @Description Get User
// @Accept  json
// @Produce json
// @Param Authorization header string true "Bearer"
// @Success 200 {array} User
// @Router /api/v1/user [get]
func findAllUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("SELECT id, username, fname, lname, email, tel FROM users")
	u := query(params)
	return c.JSON(http.StatusOK, u)
}

// @Summary Retrieves user based on query
// @Description Post User
// @Accept  json
// @Produce json
// @Success 200 {array} User
// @Router /api/v1/login [post]
func Login(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}

	sqlStatement := `SELECT id, password FROM users WHERE username=$1`
	row := db.SqliteHandler.Conn.QueryRow(sqlStatement, u.Username)
	var (
		id       int
		password string
	)
	err = row.Scan(&id, &password)
	if err != nil {
		fmt.Println(err)
	}

	var user User
	user.ID = id
	user.Password = password

	if !middleware.ComparePasswords(user.Password, []byte(u.Password)) {
		return echo.ErrUnauthorized
	}

	// // Createa token
	token := jwt.New(jwt.SigningMethodHS256)

	// // Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = json.Number(strconv.FormatInt(int64(user.ID), 10))
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

// @Summary Reset password
// @Description Patch resetPass
// @Accept  json
// @Produce json
// @Success 200 {array} resetPass
// @Router /api/v1/user/resetpassword [patch]
func resetPassword(c echo.Context) (err error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)
	userID := claims.ID

	// log.Printf("Log reset password : %d | %v", id, c)
	u := new(resetPass)
	if err = c.Bind(u); err != nil {
		return
	}
	log.Printf("Log user password : %v | %v", u.Password, u.NewPassword)
	sqlStatement := `SELECT password from users where id=$1`
	row := db.SqliteHandler.Conn.QueryRow(sqlStatement, userID)
	var (
		password string
	)
	err = row.Scan(&password)
	if err != nil {
		fmt.Println(err)
	}

	var user User
	user.Password = password

	if !middleware.ComparePasswords(user.Password, []byte(u.Password)) {
		return echo.ErrNotFound
	}

	hashPwd := middleware.HashAndSalt([]byte(u.NewPassword))

	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("UPDATE users SET password='%v' where id = %d", hashPwd, userID)
	execute(params)
	return c.JSON(http.StatusOK, "Update password success.")
}
