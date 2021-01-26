package user

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username" validate:"required,min=8,max=30"`
	Password  string `json:"password" validate:"required,min=8,max=30"`
	Fname     string `json:"fname" validate:",min=2,max=100"`
	Lname     string `json:"lname" validate:",min=2,max=100"`
	Email     string `json:"email" validate:",email"`
	Telephone string `json:"tel"validate:",min=2,max=20"`
}

type JwtCustomClaims struct {
	ID    int  `json:"id"`
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

type resetPass struct {
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
