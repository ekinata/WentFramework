//
//    TODO: Make migration command to auto-inspect models
//

package commands

import (
	"fmt"
	"went-framework/app/database"
	"went-framework/app/models"
	"went-framework/internal/logger"
)

func Migrate() {
	database.Connect()

	err := database.DB.AutoMigrate(
		// System Tables
		&logger.LogEntry{}, // Add logs table

		&models.User{},
		// ... Add other models here as needed
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration completed.")
}

func MigrateFresh() {
	database.Connect()

	// Tüm tabloları sil
	err := database.DB.Migrator().DropTable(
		// System Tables
		&logger.LogEntry{}, // Add logs table

		&models.User{},
		// ... Add other models here as needed
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("All tables dropped. Recreating...")

	// Tabloları yeniden oluştur
	Migrate()
}

func MigrateRollback() {
	database.Connect()

	// Tüm tabloları sil
	err := database.DB.Migrator().DropTable(
		// System Tables
		&logger.LogEntry{}, // Add logs table

		&models.User{},
		// ... Add other models here as needed
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration rollback completed.")
}
