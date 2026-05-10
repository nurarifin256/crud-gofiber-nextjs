package configs

import (
	"fmt"
	"time"

	dberrors "simple-crud/pkg/dberror"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// Open opens database connection and assigns to DB global variable
func Open(cfg DBConf) *gorm.DB {
	var dial gorm.Dialector
	switch cfg.Driver {
	case "mysql":

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Asia%%2FJakarta&parseTime=true&%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.Params)
		dial = mysql.Open(dsn)
	case "postgres":

		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s timezone=Asia/Jakarta %s", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name, cfg.Params)
		dial = postgres.Open(dsn)
	case "sqlserver":

		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.Params)
		dial = sqlserver.Open(dsn)
	default:
		err := fmt.Errorf("unsupported driver: %s", cfg.Driver)
		fmt.Printf("Database configuration error: %v\n", err)
		return nil
	}

	gdb, err := gorm.Open(dial, &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Warn),
		NowFunc:                func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		fmt.Printf("Failed to connect to %s database: %v\n", cfg.Driver, err)
		return nil
	}

	_ = gdb.Use(dberrors.NewFriendlyErrors())

	fmt.Printf(" (Driver: %s) %s database connected via GORM\n", cfg.Driver, cfg.Name)
	return gdb
}
