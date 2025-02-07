package configs

import (
	"fmt"
	"log"
	"os"

	// "simple-crud/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance
var errEnv error

func ConnectDb() {
	errEnv = godotenv.Load(".env")

	if errEnv != nil {
		log.Fatal("Error loading .env file", errEnv)
	}

	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		dbUser   = os.Getenv("DB_USER")
		dbName   = os.Getenv("DB_NAME")
		password = os.Getenv("DB_PASSWORD")
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, dbUser, password, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Error connecting to database. \n", err)
	}

	log.Println("Database connected")
	db.Logger = db.Logger.LogMode(logger.Info)

	// migration table
	// db.AutoMigrate(
	// 	&models.User{},
	// 	&models.UserToken{},
	// )

	DB = Dbinstance{
		Db: db,
	}
}
