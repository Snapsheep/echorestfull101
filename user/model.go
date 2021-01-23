package user

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Fname     string `json:"fname" validate:"required"`
	Lname     string `json:"lname" validate:"required"`
	Email     string `json:"email"`
	Telephone string `json:"tel"`
}

type JwtCustomClaims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

type UserLogin struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
