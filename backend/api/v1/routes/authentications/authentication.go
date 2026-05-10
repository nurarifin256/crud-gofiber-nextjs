package authentications

import (
	handler "simple-crud/api/v1/handlers/authentications"
	users "simple-crud/internal/users"
	"simple-crud/pkg/appcontext"
	"sync"
)

type AuthRouteInitializer struct {
	userService users.Service
	once        sync.Once
}

func NewAuthRouteInitializer() *AuthRouteInitializer {
	return &AuthRouteInitializer{}
}

func (d *AuthRouteInitializer) Initialize() []appcontext.RouteDefinition {
	d.once.Do(func() {
		d.userService = users.InitializeService()
	})

	return AuthRoutes(d.userService)
}

func AuthRoutes(service users.Service) []appcontext.RouteDefinition {
	authHandlers := handler.NewAuthHandler(service)
	return []appcontext.RouteDefinition{
		{
			Method:       "GET",
			Path:         "/user/find-by-nik/:nik",
			Action:       authHandlers.FindUserByNik,
			AuthRequired: false,
		},
	}
}
