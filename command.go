package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
	"went-framework/app/database"
	wentlog "went-framework/internal/logger"
	"went-framework/internal/swagger"
	"went-framework/router"
)

// StartServer starts the HTTP server
func StartServer() {
	start := time.Now()
	wentlog.Info("Starting WentFramework server...")

	// Setup routes using the router package
	r := router.SetupRoutes()

	// Get server configuration from environment
	port := getEnv("SERVER_PORT", "3000")
	host := getEnv("SERVER_HOST", "0.0.0.0")

	wentlog.Info("Server configuration loaded", map[string]interface{}{
		"host": host,
		"port": port,
		"env":  getEnv("APP_ENV", "development"),
	})

	// Print available routes (now automatically generated)
	router.PrintRoutes(r)

	// Override the port in PrintRoutes output
	fmt.Printf("🌐 Server will bind to %s:%s\n", host, port)

	// Start the server
	address := fmt.Sprintf("%s:%s", host, port)

	wentlog.Infof("Server startup completed in %v", time.Since(start))
	wentlog.Infof("Server listening on %s", address)

	// Log server shutdown on exit
	defer wentlog.Info("Server shutdown completed")

	// Start the server - this will block until the server shuts down
	if err := http.ListenAndServe(address, r); err != nil {
		wentlog.Errorf("Server failed: %v", err)
		fmt.Printf("❌ Server failed: %v\n", err)
		os.Exit(1)
	}
}

// TestDatabaseConnection tests the database connection
func TestDatabaseConnection() {
	wentlog.Info("Testing database connection...")

	fmt.Println("🔌 Testing database connection...")
	fmt.Printf("📊 Database Config:\n")
	fmt.Printf("   Host: %s\n", getEnv("DB_HOST", "localhost"))
	fmt.Printf("   Port: %s\n", getEnv("DB_PORT", "5432"))
	fmt.Printf("   User: %s\n", getEnv("DB_USER", "postgres"))
	fmt.Printf("   Database: %s\n", getEnv("DB_NAME", "testdb"))
	fmt.Printf("   SSL Mode: %s\n", getEnv("DB_SSLMODE", "disable"))

	start := time.Now()

	// Attempt to connect
	database.Connect()

	duration := time.Since(start)

	wentlog.Info("Database connection successful", map[string]interface{}{
		"host":         getEnv("DB_HOST", "localhost"),
		"database":     getEnv("DB_NAME", "testdb"),
		"connect_time": duration.Milliseconds(),
	})

	fmt.Println("✅ Database connection successful!")
}

// GenerateSwaggerDocs generates Swagger documentation
func GenerateSwaggerDocs() {
	wentlog.Info("Starting Swagger documentation generation...")

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
		wentlog.Errorf("Error generating swagger spec: %v", err)
		fmt.Printf("❌ Failed to generate Swagger documentation: %v\n", err)
		return
	}

	// Save to file
	filename := "docs/swagger.json"
	if err := swagger.SaveSwaggerSpec(spec, filename); err != nil {
		wentlog.Errorf("Error saving swagger spec: %v", err)
		fmt.Printf("❌ Failed to save Swagger documentation: %v\n", err)
		return
	}

	wentlog.Info("Swagger documentation generated successfully", map[string]interface{}{
		"filename": filename,
		"host":     host,
		"port":     port,
	})

	fmt.Printf("✅ Swagger documentation generated successfully!\n")
	fmt.Printf("📄 Saved to: %s\n", filename)
	fmt.Printf("🌐 When server is running, view at: http://%s:%s/swagger/\n", host, port)
}

// MakeModel creates model and controller files from templates
func MakeModel(modelName string) {
	wentlog.Info("Starting model generation", map[string]interface{}{
		"model_name": modelName,
	})

	modelPath := "app/models/" + modelName + ".go"
	controllerPath := "app/controllers/" + modelName + "Controller.go"

	createFileFromTemplate("internal/templates/model.tpl", modelPath, modelName)
	createFileFromTemplate("internal/templates/controller.tpl", controllerPath, modelName)

	wentlog.Info("Model generation completed", map[string]interface{}{
		"model_name":      modelName,
		"model_file":      modelPath,
		"controller_file": controllerPath,
	})

	fmt.Println("Files created successfully!")
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// createFileFromTemplate creates a file from a template
func createFileFromTemplate(templatePath, outputPath, modelName string) {
	wentlog.Debug("Creating file from template", map[string]interface{}{
		"template_path": templatePath,
		"output_path":   outputPath,
		"model_name":    modelName,
	})

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		wentlog.Error("Failed to parse template", map[string]interface{}{
			"template_path": templatePath,
			"error":         err.Error(),
		})
		fmt.Printf("❌ Error parsing template %s: %v\n", templatePath, err)
		return
	}

	os.MkdirAll(getDir(outputPath), os.ModePerm)

	if _, err := os.Stat(outputPath); err == nil {
		wentlog.Warn("File already exists, skipping", map[string]interface{}{
			"output_path": outputPath,
		})
		fmt.Printf("Skipped (already exists): %s\n", outputPath)
		return
	}

	f, err := os.Create(outputPath)
	if err != nil {
		wentlog.Error("Failed to create output file", map[string]interface{}{
			"output_path": outputPath,
			"error":       err.Error(),
		})
		fmt.Printf("❌ Error creating file %s: %v\n", outputPath, err)
		return
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
		wentlog.Error("Failed to execute template", map[string]interface{}{
			"template_path": templatePath,
			"output_path":   outputPath,
			"error":         err.Error(),
		})
		fmt.Printf("❌ Error executing template: %v\n", err)
		return
	}

	wentlog.Debug("File created successfully", map[string]interface{}{
		"output_path": outputPath,
	})
}

// getDir gets the directory part of a file path
func getDir(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx == -1 {
		return "."
	}
	return path[:idx]
}
