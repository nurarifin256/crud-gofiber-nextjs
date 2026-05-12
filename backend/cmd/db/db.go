package main

import (
	"flag"
	"fmt"
	"simple-crud/configs"
	"simple-crud/database/migrations"
)

func main() {
	cfg := configs.Load()
	configs.DB = configs.Open(cfg.DB)

	migrate := flag.Bool("db:migrate", false, "Run all migrations")
	rollback := flag.Bool("db:rollback", false, "Rollback the last migration")

	flag.Parse()
	if *migrate {
		fmt.Println("Running migrations...")
		migrations.MigrateAll()
		fmt.Println("Migrations completed.")
	} else if *rollback {
		fmt.Println("Rolling back the last migration...")
		migrations.RollbackAll()
		fmt.Println("Rollback completed.")
	} else {
		fmt.Println("No database operation specified. Use -db:migrate to run migrations or -db:rollback to rollback.")
	}
}
