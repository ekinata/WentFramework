//
//    TODO: Generation should auto-detect and include all models for default crud operations
//

package swagger

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"went-framework/app/models"

	"github.com/gorilla/mux"
)

// SwaggerInfo holds the basic API information
type SwaggerInfo struct {
	Version     string `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Host        string `json:"host"`
	BasePath    string `json:"basePath"`
}

// SwaggerSpec represents the complete OpenAPI 3.0 specification
type SwaggerSpec struct {
	OpenAPI    string              `json:"openapi"`
	Info       SwaggerInfo         `json:"info"`
	Servers    []Server            `json:"servers"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

// Server represents an API server
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// PathItem represents a path and its operations
type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
}

// Operation represents an API operation
type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary"`
	Description string              `json:"description,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

// Parameter represents an operation parameter
type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Schema      Schema `json:"schema"`
}

// RequestBody represents a request body
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content"`
}

// Response represents an API response
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// MediaType represents a media type object
type MediaType struct {
	Schema Schema `json:"schema"`
}

// Components holds reusable objects
type Components struct {
	Schemas map[string]Schema `json:"schemas"`
}

// Schema represents a JSON schema
type Schema struct {
	Type                 string            `json:"type,omitempty"`
	Format               string            `json:"format,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty"`
	Required             []string          `json:"required,omitempty"`
	Items                *Schema           `json:"items,omitempty"`
	Ref                  string            `json:"$ref,omitempty"`
	AdditionalProperties interface{}       `json:"additionalProperties,omitempty"`
	Example              interface{}       `json:"example,omitempty"`
}

// RouteInfo holds information about a route
type RouteInfo struct {
	Method      string
	Path        string
	HandlerName string
	Tags        []string
}

// GenerateSwagger generates OpenAPI/Swagger documentation
func GenerateSwagger(router *mux.Router, info SwaggerInfo) (*SwaggerSpec, error) {
	spec := &SwaggerSpec{
		OpenAPI: "3.0.0",
		Info:    info,
		Servers: []Server{
			{
				URL:         fmt.Sprintf("http://%s", info.Host),
				Description: "Development server",
			},
		},
		Paths: make(map[string]PathItem),
		Components: Components{
			Schemas: make(map[string]Schema),
		},
	}

	// Extract routes from router
	routes := extractRoutes(router)

	// Generate schemas from models
	if err := generateSchemas(spec); err != nil {
		return nil, fmt.Errorf("error generating schemas: %v", err)
	}

	// Generate paths from routes
	if err := generatePaths(spec, routes); err != nil {
		return nil, fmt.Errorf("error generating paths: %v", err)
	}

	return spec, nil
}

// extractRoutes extracts route information from the mux router
func extractRoutes(router *mux.Router) []RouteInfo {
	var routes []RouteInfo

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		for _, method := range methods {
			// Extract handler name and tags
			handlerName := extractHandlerName(pathTemplate, method)
			tags := extractTags(pathTemplate)

			routes = append(routes, RouteInfo{
				Method:      method,
				Path:        pathTemplate,
				HandlerName: handlerName,
				Tags:        tags,
			})
		}

		return nil
	})

	return routes
}

// generateSchemas generates schema definitions from models
func generateSchemas(spec *SwaggerSpec) error {
	// Generate User schema
	userSchema := generateModelSchema(reflect.TypeOf(models.User{}))
	spec.Components.Schemas["User"] = userSchema

	// Generate common response schemas
	spec.Components.Schemas["Response"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"status": {
				Type:    "string",
				Example: "success",
			},
			"message": {
				Type:    "string",
				Example: "Operation completed successfully",
			},
			"data": {
				AdditionalProperties: true,
			},
		},
		Required: []string{"status", "message"},
	}

	spec.Components.Schemas["ErrorResponse"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"status": {
				Type:    "string",
				Example: "error",
			},
			"message": {
				Type:    "string",
				Example: "Error description",
			},
		},
		Required: []string{"status", "message"},
	}

	return nil
}

// generateModelSchema generates a schema from a Go struct type
func generateModelSchema(t reflect.Type) Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	schema := Schema{
		Type:       "object",
		Properties: make(map[string]Schema),
		Required:   []string{},
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse JSON tag
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			jsonName = strings.ToLower(field.Name)
		}

		// Generate field schema
		fieldSchema := generateFieldSchema(field.Type)

		// Add example based on field name
		if example := generateExample(jsonName, field.Type); example != nil {
			fieldSchema.Example = example
		}

		schema.Properties[jsonName] = fieldSchema

		// Check if field is required (not omitempty)
		if !strings.Contains(jsonTag, "omitempty") {
			schema.Required = append(schema.Required, jsonName)
		}
	}

	return schema
}

// generateFieldSchema generates schema for a struct field
func generateFieldSchema(t reflect.Type) Schema {
	switch t.Kind() {
	case reflect.String:
		return Schema{Type: "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Schema{Type: "integer"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Schema{Type: "integer", Format: "int64"}
	case reflect.Float32, reflect.Float64:
		return Schema{Type: "number"}
	case reflect.Bool:
		return Schema{Type: "boolean"}
	case reflect.Slice:
		return Schema{
			Type:  "array",
			Items: &Schema{Type: "string"}, // Default to string, can be improved
		}
	case reflect.Struct:
		// Handle time.Time specifically
		if t.String() == "time.Time" {
			return Schema{Type: "string", Format: "date-time"}
		}
		return generateModelSchema(t)
	default:
		return Schema{Type: "string"}
	}
}

// generateExample generates example values based on field name and type
func generateExample(fieldName string, t reflect.Type) interface{} {
	switch strings.ToLower(fieldName) {
	case "id":
		return 1
	case "name":
		return "John Doe"
	case "email":
		return "john@example.com"
	case "created_at", "updated_at":
		return "2025-07-31T15:42:18.792477+03:00"
	default:
		switch t.Kind() {
		case reflect.String:
			return "example string"
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
			return 1
		case reflect.Bool:
			return true
		}
	}
	return nil
}

// generatePaths generates API paths documentation
func generatePaths(spec *SwaggerSpec, routes []RouteInfo) error {
	for _, route := range routes {
		pathItem, exists := spec.Paths[route.Path]
		if !exists {
			pathItem = PathItem{}
		}

		operation := generateOperation(route)

		switch route.Method {
		case "GET":
			pathItem.Get = operation
		case "POST":
			pathItem.Post = operation
		case "PUT":
			pathItem.Put = operation
		case "DELETE":
			pathItem.Delete = operation
		}

		spec.Paths[route.Path] = pathItem
	}

	return nil
}

// generateOperation generates operation documentation
func generateOperation(route RouteInfo) *Operation {
	operation := &Operation{
		Tags:      route.Tags,
		Summary:   generateSummary(route.Method, route.Path),
		Responses: generateResponses(route.Method, route.Path),
	}

	// Add path parameters
	if strings.Contains(route.Path, "{id}") {
		operation.Parameters = []Parameter{
			{
				Name:        "id",
				In:          "path",
				Description: "Resource ID",
				Required:    true,
				Schema:      Schema{Type: "integer"},
			},
		}
	}

	// Add request body for POST and PUT
	if route.Method == "POST" || route.Method == "PUT" {
		operation.RequestBody = generateRequestBody(route.Path)
	}

	return operation
}

// generateSummary generates operation summary
func generateSummary(method, path string) string {
	cleanPath := strings.TrimPrefix(path, "/api")

	switch {
	case path == "/api/health":
		return "Health check"
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
		return fmt.Sprintf("%s %s", method, cleanPath)
	}
}

// generateResponses generates response documentation
func generateResponses(method, path string) map[string]Response {
	responses := make(map[string]Response)

	if path == "/api/health" {
		responses["200"] = Response{
			Description: "Health check successful",
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{
						Type: "object",
						Properties: map[string]Schema{
							"status":  {Type: "string", Example: "healthy"},
							"message": {Type: "string", Example: "Server is running"},
						},
					},
				},
			},
		}
		return responses
	}

	// Standard responses
	switch method {
	case "GET":
		if strings.Contains(path, "{id}") {
			responses["200"] = Response{
				Description: "Resource retrieved successfully",
				Content: map[string]MediaType{
					"application/json": {
						Schema: Schema{Ref: "#/components/schemas/Response"},
					},
				},
			}
			responses["404"] = Response{
				Description: "Resource not found",
				Content: map[string]MediaType{
					"application/json": {
						Schema: Schema{Ref: "#/components/schemas/ErrorResponse"},
					},
				},
			}
		} else {
			responses["200"] = Response{
				Description: "Resources retrieved successfully",
				Content: map[string]MediaType{
					"application/json": {
						Schema: Schema{Ref: "#/components/schemas/Response"},
					},
				},
			}
		}
	case "POST":
		responses["201"] = Response{
			Description: "Resource created successfully",
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{Ref: "#/components/schemas/Response"},
				},
			},
		}
	case "PUT":
		responses["200"] = Response{
			Description: "Resource updated successfully",
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{Ref: "#/components/schemas/Response"},
				},
			},
		}
	case "DELETE":
		responses["200"] = Response{
			Description: "Resource deleted successfully",
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{Ref: "#/components/schemas/Response"},
				},
			},
		}
	}

	// Common error responses
	responses["400"] = Response{
		Description: "Bad request",
		Content: map[string]MediaType{
			"application/json": {
				Schema: Schema{Ref: "#/components/schemas/ErrorResponse"},
			},
		},
	}

	responses["500"] = Response{
		Description: "Internal server error",
		Content: map[string]MediaType{
			"application/json": {
				Schema: Schema{Ref: "#/components/schemas/ErrorResponse"},
			},
		},
	}

	return responses
}

// generateRequestBody generates request body documentation
func generateRequestBody(path string) *RequestBody {
	if strings.Contains(path, "/users") {
		return &RequestBody{
			Description: "User data",
			Required:    true,
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{
						Type: "object",
						Properties: map[string]Schema{
							"name":  {Type: "string", Example: "John Doe"},
							"email": {Type: "string", Example: "john@example.com"},
						},
						Required: []string{"name", "email"},
					},
				},
			},
		}
	}

	return nil
}

// extractHandlerName extracts handler name from path and method
func extractHandlerName(path, method string) string {
	return fmt.Sprintf("%s%s", method, strings.ReplaceAll(path, "/", "_"))
}

// extractTags extracts tags from path
func extractTags(path string) []string {
	if strings.Contains(path, "/users") {
		return []string{"Users"}
	}
	if strings.Contains(path, "/health") {
		return []string{"Health"}
	}
	return []string{"API"}
}

// SaveSwaggerSpec saves the swagger specification to a file
func SaveSwaggerSpec(spec *SwaggerSpec, filename string) error {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// ParseControllerComments parses comments from controller files for additional documentation
func ParseControllerComments(controllerDir string) (map[string]string, error) {
	comments := make(map[string]string)

	err := filepath.Walk(controllerDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Doc != nil {
					funcName := fn.Name.Name
					comment := fn.Doc.Text()
					comments[funcName] = strings.TrimSpace(comment)
				}
			}
		}

		return nil
	})

	return comments, err
}
