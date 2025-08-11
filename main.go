package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"went-framework/commands"
	"went-framework/database"
	"went-framework/router"
	"went-framework/swagger"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

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
		createFileFromTemplate("templates/model.tpl", "models/"+modelName+".go", modelName)
		createFileFromTemplate("templates/controller.tpl", "controllers/"+modelName+"Controller.go", modelName)
		fmt.Println("Files created successfully!")

	case "migrate":
		commands.Migrate()

	case "migrate:fresh":
		commands.MigrateFresh()

	case "migrate:rollback":
		commands.MigrateRollback()

	case "serve":
		startServer()

	case "db:test":
		testDatabaseConnection()

	case "swagger:generate":
		generateSwaggerDocs()

	default:
		fmt.Println("Unknown command:", command)
	}
}

func startServer() {
	// Setup routes using the router package
	r := router.SetupRoutes()

	// Get server configuration from environment
	port := getEnv("SERVER_PORT", "3000")
	host := getEnv("SERVER_HOST", "0.0.0.0")

	// Print available routes (now automatically generated)
	router.PrintRoutes(r)

	// Override the port in PrintRoutes output
	fmt.Printf("🌐 Server will bind to %s:%s\n", host, port)

	// Start the server
	address := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting server on %s", address)
	log.Fatal(http.ListenAndServe(address, r))
}

func testDatabaseConnection() {
	fmt.Println("🔌 Testing database connection...")
	fmt.Printf("📊 Database Config:\n")
	fmt.Printf("   Host: %s\n", getEnv("DB_HOST", "localhost"))
	fmt.Printf("   Port: %s\n", getEnv("DB_PORT", "5432"))
	fmt.Printf("   User: %s\n", getEnv("DB_USER", "postgres"))
	fmt.Printf("   Database: %s\n", getEnv("DB_NAME", "testdb"))
	fmt.Printf("   SSL Mode: %s\n", getEnv("DB_SSLMODE", "disable"))

	// Attempt to connect
	database.Connect()

	fmt.Println("✅ Database connection successful!")
}

func generateSwaggerDocs() {
	fmt.Println("📚 Generating Swagger documentation...")

	// Setup routes to analyze
	r := router.SetupRoutes()

	// Generate swagger specification
	host := getEnv("SERVER_HOST", "localhost")
	port := getEnv("SERVER_PORT", "3000")

	if host == "0.0.0.0" {
		host = "localhost"
	}

	info := swagger.SwaggerInfo{
		Version:     getEnv("APP_VERSION", "1.0.0"),
		Title:       getEnv("APP_NAME", "WentFramework API"),
		Description: "Auto-generated API documentation for WentFramework - A lightweight Go framework for building RESTful APIs",
		Host:        fmt.Sprintf("%s:%s", host, port),
		BasePath:    "/api",
	}

	spec, err := swagger.GenerateSwagger(r, info)
	if err != nil {
		log.Fatalf("Error generating swagger spec: %v", err)
	}

	// Save to file
	filename := "docs/swagger.json"
	if err := swagger.SaveSwaggerSpec(spec, filename); err != nil {
		log.Fatalf("Error saving swagger spec: %v", err)
	}

	fmt.Printf("✅ Swagger documentation generated successfully!\n")
	fmt.Printf("📄 Saved to: %s\n", filename)
	fmt.Printf("🌐 When server is running, view at: http://%s:%s/swagger/\n", host, port)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func createFileFromTemplate(templatePath, outputPath, modelName string) {
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}

	os.MkdirAll(getDir(outputPath), os.ModePerm)

	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("Skipped (already exists): %s\n", outputPath)
		return
	}

	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := struct {
		ModelName string
		TableName string
	}{
		ModelName: modelName,
		TableName: strings.ToLower(modelName) + "s",
	}

	err = tpl.Execute(f, data)
	if err != nil {
		panic(err)
	}
}

func getDir(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx == -1 {
		return "."
	}
	return path[:idx]
}
