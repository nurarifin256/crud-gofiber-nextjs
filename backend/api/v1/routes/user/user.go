package user

import (
	handler "simple-crud/api/v1/handlers/user"
	users "simple-crud/internal/users"
	"simple-crud/pkg/appcontext"
	"sync"
)

type UserRouteInitializer struct {
	userService users.Service
	once        sync.Once
}

func NewUserRouteInitializer() *UserRouteInitializer {
	return &UserRouteInitializer{}
}

func (d *UserRouteInitializer) Initialize() []appcontext.RouteDefinition {
	d.once.Do(func() {
		d.userService = users.InitializeService()
	})

	return UserRoutes(d.userService)
}

func UserRoutes(service users.Service) []appcontext.RouteDefinition {
	userHandlers := handler.NewUserHandler(service)
	return []appcontext.RouteDefinition{
		{
			Method:       "POST",
			Path:         "user/submit",
			Action:       userHandlers.SubmitUser,
			AuthRequired: false,
		},
	}
}
