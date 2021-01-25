package user

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username" validate:"required,min=8,max=30"`
	Password  string `json:"password" validate:"required,min=8,max=30"`
	Fname     string `json:"fname" validate:"required,min=2,max=100"`
	Lname     string `json:"lname" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Telephone string `json:"tel"validate:"required,min=2,max=20"`
}

type JwtCustomClaims struct {
	ID    int  `json:"id"`
	Admin bool `json:"admin"`
	jwt.StandardClaims
}
