package domain

import "github.com/golang-jwt/jwt/v4"

type Credentials struct {
	Id       string `json:"-" swaggerignore:"true"`
	Login    string `json:"login"`
	Password string `json:"password"`
	PassHash string `json:"-" swaggerignore:"true"`
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}
