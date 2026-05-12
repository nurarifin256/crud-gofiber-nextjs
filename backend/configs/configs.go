package configs

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	Name            string
	Env             string
	Host            string
	Port            int
	SecretKey       string
	CorsAllowOrigin string
	CorsHeaders     string
}

type DBConf struct {
	Driver  string
	Host    string
	Port    int
	User    string
	Pass    string
	Name    string
	Params  string
	Migrate bool
}

type FtpConfig struct {
	Host     string
	Username string
	Password string
	BaseDir  string
	UseTLS   bool
	Timeout  time.Duration
}

type Config struct {
	App App
	DB  DBConf
	Ftp FtpConfig
}

func Load() Config {
	_ = godotenv.Load()

	v := viper.New()
	v.AutomaticEnv()

	// default
	v.SetDefault("APP_NAME", "crud-gofiber-nextjs")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_HOST", "0.0.0.0")
	v.SetDefault("APP_PORT", 8001)
	v.SetDefault("DB_DRIVER", "postgres")
	v.SetDefault("DB_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local")
	v.SetDefault("AUTO_MIGRATE", "1")

	cfg := Config{
		App: App{
			Name:            v.GetString("APP_NAME"),
			Env:             v.GetString("APP_ENV"),
			Host:            v.GetString("APP_HOST"),
			Port:            v.GetInt("APP_PORT"),
			SecretKey:       v.GetString("SECRET_KEY"),
			CorsAllowOrigin: v.GetString("CORS_ALLOW_ORIGIN"),
			CorsHeaders:     v.GetString("CORS_HEADERS"),
		},
		DB: DBConf{
			Driver:  v.GetString("DB_DRIVER"),
			Host:    v.GetString("DB_HOST"),
			Port:    v.GetInt("DB_PORT"),
			User:    v.GetString("DB_USER"),
			Pass:    v.GetString("DB_PASS"),
			Name:    v.GetString("DB_NAME"),
			Params:  v.GetString("DB_PARAMS"),
			Migrate: v.GetBool("AUTO_MIGRATE"),
		},
	}

	log.Printf("loaded config app=%s env=%s driver=%s", cfg.App.Name, cfg.App.Env, cfg.DB.Driver)
	return cfg
}

func (a App) Addr() string { return fmt.Sprintf("%s:%d", a.Host, a.Port) }
