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

func CreateUser(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("INSERT INTO users (username, password, fname, lname, email, tel) VALUES ('%v','%v','%v','%v','%v','%v')", u.Username, middleware.HashAndSalt([]byte(u.Password)), u.Fname, u.Lname, u.Email, u.Telephone)
	execute(params)
	return c.JSON(http.StatusOK, "Create user success.")
}

func getUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	id, _ := strconv.Atoi(c.Param("id"))
	params["sql"] = fmt.Sprintf("select id, username, fname, lname, email, tel from users where id = %d", id)
	u := query(params)
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
	execute(params)
	return c.JSON(http.StatusOK, u)
}

func findAllUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	params["sql"] = fmt.Sprintf("select id, username, fname, lname, email, tel from users")
	u := query(params)
	return c.JSON(http.StatusOK, u)
}

func deleteUser(c echo.Context) error {
	params := make(map[interface{}]interface{})
	id, _ := strconv.Atoi(c.Param("id"))
	params["sql"] = fmt.Sprintf("delete from users where id = %d", id)
	execute(params)
	return c.NoContent(http.StatusNoContent)
}

func Login(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}

	sqlStatement := `SELECT id, password from users where username=$1`
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

	// log.Printf("Log user password : %v | %v", user.Password, u.Password)

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

func resetPassword(c echo.Context) (err error) {
	type resetPass struct {
		Password    string `json:"password"`
		NewPassword string `json:"newpassword"`
	}
	id, _ := strconv.Atoi(c.Param("id"))
	log.Printf("Log reset password : %d | %v", id, c)
	u := new(resetPass)
	if err = c.Bind(u); err != nil {
		return
	}
	log.Printf("Log user password : %v | %v", u.Password, u.NewPassword)
	sqlStatement := `SELECT password from users where id=$1`
	row := db.SqliteHandler.Conn.QueryRow(sqlStatement, id)
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
	params["sql"] = fmt.Sprintf("UPDATE users SET password='%v' where id = %d", hashPwd, id)
	execute(params)
	return c.JSON(http.StatusOK, "Update password success.")
}

func Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}
