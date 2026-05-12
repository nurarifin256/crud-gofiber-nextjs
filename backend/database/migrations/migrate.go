package migrations

import "fmt"

func MigrateAll() {
	fmt.Println("Migrating all tables...")
	UpMUserTable()

	fmt.Println("All tables migrated successfully.")
}

func RollbackAll() {
	fmt.Println("Rolling back all tables...")
	DownMUserTable()
	fmt.Println("All tables rolled back successfully.")
}
