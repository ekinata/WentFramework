package main

import (
	"fmt"
	"log"
	"os"
	"went-framework/internal/commands"
	wentlog "went-framework/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	// Initialize logger
	wentlog.Init()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [make:model <ModelName>|migrate|migrate:fresh|migrate:rollback|serve|db:test|swagger:generate]")
		return
	}

	command := os.Args[1]

	switch command {
	case "make:model":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go make:model ModelName")
			return
		}
		modelName := os.Args[2]
		commands.MakeModel(modelName)

	case "migrate":
		commands.Migrate()

	case "migrate:fresh":
		commands.MigrateFresh()

	case "migrate:rollback":
		commands.MigrateRollback()

	case "serve":
		commands.StartServer()

	case "db:test":
		commands.TestDatabaseConnection()

	case "swagger:generate":
		commands.GenerateSwaggerDocs()

	default:
		fmt.Println("Unknown command:", command)
	}
}
