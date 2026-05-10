package middleware

import (
	"simple-crud/configs"
	users "simple-crud/internal/users"
)

type Middleware struct {
	AuthMiddleware *AuthMiddlewareHandler
}

func InitMiddleWare() *Middleware {
	secret := configs.Load().App.SecretKey
	// initialization service & middleware
	userService := users.InitializeService()
	return &Middleware{
		AuthMiddleware: NewAuthMiddleware(userService, secret),
	}
}
