package migrations

import (
	"simple-crud/configs"
	models "simple-crud/database/models"
)

func UpMUserTable() {
	configs.DB.AutoMigrate(&models.User{})
}

func DownMUserTable() {
	configs.DB.Migrator().DropTable(&models.User{})
}
