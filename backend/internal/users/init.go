package user

import "simple-crud/configs"

func InitializeService() Service {
	if configs.DB == nil {
		panic("Database is not initialized, Please ensure configs initdb")
	}

	userRepo := NewUserRepository(configs.DB)
	userService := NewUserService(userRepo)
	return userService
}
