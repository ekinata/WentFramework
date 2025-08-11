package commands

import (
	"fmt"
	"went-framework/database"
	"went-framework/models"
)

func Migrate() {
	database.Connect()

	err := database.DB.AutoMigrate(
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
		&models.User{},
		// ... Add other models here as needed
	)
	if err != nil {
		panic(err)
	}

	// Tabloları yeniden oluştur
	Migrate()
}

func MigrateRollback() {
	database.Connect()

	// Tüm tabloları sil
	err := database.DB.Migrator().DropTable(
		&models.User{},
		// ... Add other models here as needed
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration rollback completed.")
}
