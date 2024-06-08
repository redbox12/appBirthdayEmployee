package models

import "github.com/dgrijalva/jwt-go"

type Token struct {
	UserId int
	jwt.StandardClaims
}

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Token      string `json:"token"`
}

type Birthday struct {
}

type Notification struct {
}

type UsersNotifications struct {
}
