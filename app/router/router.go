package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"went-framework/internal/middleware"
	"went-framework/internal/swagger"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configures and returns the main router with all routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.MiddlewareChain(
		middleware.RequestIDMiddleware,
		middleware.CORSMiddleware,
		middleware.LoggingMiddleware,
	))

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Setup different route groups
	setupUserRoutes(api)
	setupHealthRoutes(api)
	setupSwaggerRoutes(router)

	// Placeholder for other route groups
	setupOtherRoutes(api)

	return router
}

func setupOtherRoutes(api *mux.Router) {
	// Placeholder for other routes
	// Add more route groups as needed
}

// setupHealthRoutes configures health check routes
func setupHealthRoutes(api *mux.Router) {
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "healthy", "message": "Server is running"}`)
	}).Methods("GET")
}

// setupSwaggerRoutes configures Swagger documentation routes
func setupSwaggerRoutes(router *mux.Router) {
	// Swagger JSON endpoint
	router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		// Generate swagger spec dynamically
		spec, err := generateSwaggerSpec(router)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating swagger spec: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(spec)
	}).Methods("GET")

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))
}

// generateSwaggerSpec generates the Swagger specification for the current router
func generateSwaggerSpec(router *mux.Router) (*swagger.SwaggerSpec, error) {
	host := getEnv("SERVER_HOST", "localhost")
	port := getEnv("SERVER_PORT", "3000")

	if host == "0.0.0.0" {
		host = "localhost"
	}

	info := swagger.SwaggerInfo{
		Version:     getEnv("APP_VERSION", "1.0.0"),
		Title:       getEnv("APP_NAME", "WentFramework API"),
		Description: "Auto-generated API documentation for WentFramework",
		Host:        fmt.Sprintf("%s:%s", host, port),
		BasePath:    "/api",
	}

	return swagger.GenerateSwagger(router, info)
}

// PrintRoutes displays all available routes in a formatted way
func PrintRoutes(router *mux.Router) {
	port := getEnv("SERVER_PORT", "3000")
	host := getEnv("SERVER_HOST", "0.0.0.0")

	if host == "0.0.0.0" {
		host = "localhost"
	}

	fmt.Printf("🚀 Server starting on :%s\n", port)
	fmt.Printf("📡 API available at http://%s:%s/api\n", host, port)
	fmt.Printf("📚 Swagger UI available at http://%s:%s/swagger/\n", host, port)
	fmt.Printf("📄 Swagger JSON available at http://%s:%s/swagger.json\n", host, port)
	fmt.Println("👥 Available endpoints:")

	// Extract routes from the router
	routes := extractRoutes(router)

	// Sort routes for consistent display
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	// Display routes in a formatted way
	for _, route := range routes {
		description := getRouteDescription(route.Method, route.Path)
		fmt.Printf("   %-6s %s - %s\n", route.Method, route.Path, description)
	}
}

// RouteInfo holds information about a route
type RouteInfo struct {
	Method string
	Path   string
}

// extractRoutes extracts all routes from the mux router
func extractRoutes(router *mux.Router) []RouteInfo {
	var routes []RouteInfo

	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		methods, err := route.GetMethods()
		if err != nil {
			// If no methods specified, skip this route
			return nil
		}

		for _, method := range methods {
			routes = append(routes, RouteInfo{
				Method: method,
				Path:   pathTemplate,
			})
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking routes: %v\n", err)
	}

	return routes
}

// getRouteDescription returns a human-readable description for a route
func getRouteDescription(method, path string) string {
	// Clean up the path for matching
	cleanPath := strings.TrimPrefix(path, "/api")

	switch {
	case path == "/api/health":
		return "Health check"
	case path == "/swagger.json":
		return "Swagger JSON specification"
	case strings.HasPrefix(path, "/swagger/"):
		return "Swagger UI documentation"
	case method == "GET" && cleanPath == "/users":
		return "Get all users"
	case method == "GET" && strings.Contains(cleanPath, "/users/{id}"):
		return "Get user by ID"
	case method == "POST" && cleanPath == "/users":
		return "Create new user"
	case method == "PUT" && strings.Contains(cleanPath, "/users/{id}"):
		return "Update user"
	case method == "DELETE" && strings.Contains(cleanPath, "/users/{id}"):
		return "Delete user"
	default:
		// Generate a generic description based on method and path
		pathParts := strings.Split(strings.Trim(cleanPath, "/"), "/")
		if len(pathParts) > 0 {
			resource := pathParts[0]
			if strings.Contains(cleanPath, "{id}") {
				switch method {
				case "GET":
					return fmt.Sprintf("Get %s by ID", strings.TrimSuffix(resource, "s"))
				case "PUT", "PATCH":
					return fmt.Sprintf("Update %s", strings.TrimSuffix(resource, "s"))
				case "DELETE":
					return fmt.Sprintf("Delete %s", strings.TrimSuffix(resource, "s"))
				}
			} else {
				switch method {
				case "GET":
					return fmt.Sprintf("Get all %s", resource)
				case "POST":
					return fmt.Sprintf("Create new %s", strings.TrimSuffix(resource, "s"))
				}
			}
		}
		return fmt.Sprintf("%s %s", method, cleanPath)
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
