# WentFramework

*They said us to go, so we went.*

A lightweight Go framework for building RESTful APIs with PostgreSQL database integration, migrations, and environment-based configuration.

## Features

- 🚀 **HTTP Server** with organized routing
- 🗄️ **PostgreSQL Integration** with GORM
- 🔄 **Database Migrations** with auto-migration support
- ⚙️ **Environment Configuration** with .env file support
- 🏗️ **Code Generation** for models and controllers
- 📦 **Docker Support** with docker-compose
- ☸️ **Kubernetes Ready** with complete manifests and deployment scripts
- 🔍 **CRUD Operations** with JSON API responses
- 📚 **Auto-Swagger Generation** with interactive API documentation
- 📝 **Flexible Logging System** with multiple storage backends (database, file, console)
- 🎯 **Clean Architecture** with separated command functions
- 🔄 **Global Request/Response Logging** with automatic HTTP middleware
- 🛡️ **Security Features** with CORS support and sensitive data filtering

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Commands](#commands)
- [Project Structure](#project-structure)
- [API Usage](#api-usage)
- [Logging](#logging)
- [Deployment](#deployment)
- [Development](#development)

## Installation

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for PostgreSQL)
- Git

### Clone the Repository

```bash
git clone https://github.com/ekinata/WentFramework.git
cd WentFramework
```

### Install Dependencies

```bash
go mod tidy
```

## Quick Start

### 1. Start PostgreSQL Database

```bash
docker-compose up -d
```

This will start a PostgreSQL container with the default configuration.

### 2. Set Up Environment

Copy the example environment file:

```bash
cp .env.example .env
```

The default `.env` file is configured to work with the Docker PostgreSQL setup.

### 3. Run Database Migrations

```bash
go run . migrate
```

### 4. Start the Server

```bash
go run . serve
```

The server will start on the configured port (default: `http://localhost:3000`).

### 5. Generate API Documentation

```bash
go run . swagger:generate
```

### 6. Test the API

```bash
# Health check (if available)
curl http://localhost:3000/api/health

# Get all users
curl http://localhost:3000/api/users

# Create a user
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

## Configuration

### Environment Variables

The project uses environment variables for configuration. Key variables include:

```properties
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=3000
SERVER_HOST=0.0.0.0

# Application Configuration
APP_ENV=development
APP_NAME=WentFramework
APP_VERSION=1.0.0
```

## Commands

WentFramework provides a clean command-line interface for various operations:

### Server Commands

```bash
# Start the HTTP server
go run . serve

# Test database connection
go run . db:test
```

### Migration Commands

```bash
# Run migrations (create/update tables)
go run . migrate

# Fresh migration (drop and recreate all tables)
go run . migrate:fresh

# Rollback migrations (drop all tables)
go run . migrate:rollback
```

### Code Generation Commands

```bash
# Generate model and controller files
go run . make:model ModelName

# Example: Create a Post model
go run . make:model Post
```

This will create:
- `app/models/Post.go` - Model file with GORM integration
- `app/controllers/PostController.go` - Controller file with CRUD operations

### Documentation Commands

```bash
# Generate/update Swagger API documentation
go run . swagger:generate
```

### Help

```bash
# Show available commands
go run .
```

## Project Structure

```
WentFramework/
├── .env                     # Environment variables (create from .env.example)
├── .env.example            # Environment variables template
├── .gitignore              # Git ignore rules
├── docker-compose.yml      # PostgreSQL container configuration
├── Dockerfile              # Multi-stage Docker build for the application
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies
├── main.go                 # Main application entry point (CLI routing only)
├── command.go              # Command implementations (refactored)
├── app/                    # Application core
│   ├── controllers/        # HTTP request handlers
│   │   └── UserController.go
│   ├── database/           # Database connection and configuration
│   │   └── connection.go
│   └── models/             # Data models and database operations
│       └── User.go
├── docs/                   # Generated documentation
│   ├── swagger.json       # Auto-generated OpenAPI specification
│   └── LOG.md             # Logging system documentation
├── internal/               # Internal packages
│   ├── commands/           # Database command implementations
│   │   └── migrate.go
│   ├── log/                # Logging system
│   │   └── log.go          # Logging system implementation
│   ├── middleware/         # HTTP middleware
│   │   └── logging.go      # Request/response logging middleware
│   ├── swagger/            # API documentation generation
│   │   └── generator.go
│   └── templates/          # Code generation templates
│       ├── controller.tpl
│       └── model.tpl
├── logs/                   # Log files (created when LOG_STORAGE=file)
├── k8s/                    # Kubernetes deployment manifests
│   ├── namespace.yaml     # Kubernetes namespace
│   ├── configmap.yaml     # Configuration management
│   ├── secrets.yaml       # Sensitive data management
│   ├── postgres-*.yaml    # PostgreSQL deployment and services
│   ├── wentframework-*.yaml # Application deployment and services
│   ├── ingress.yaml       # Ingress configuration
│   ├── hpa.yaml           # Horizontal Pod Autoscaler
│   ├── deploy.sh          # Automated deployment script
│   ├── cleanup.sh         # Cleanup script
│   └── README.md          # Kubernetes deployment guide
├── internal/               # Internal packages
│   ├── commands/           # Database command implementations
│   │   └── migrate.go
│   ├── swagger/            # API documentation generation
│   │   └── generator.go
│   └── templates/          # Code generation templates
│       ├── controller.tpl
│       └── model.tpl
├── router/                 # HTTP routing configuration
│   └── router.go
└── templates/              # Additional templates (legacy)
    ├── controller.tpl
    └── model.tpl
```

## API Usage

### Base URL

```
http://localhost:3000/api
```

### User Endpoints

The framework includes a complete User model with CRUD operations:

#### Get All Users

```http
GET /api/users
```

#### Get User by ID

```http
GET /api/users/{id}
```

#### Create User

```http
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

#### Update User

```http
PUT /api/users/{id}
Content-Type: application/json

{
  "name": "John Updated",
  "email": "john.updated@example.com"
}
```

#### Delete User

```http
DELETE /api/users/{id}
```

## Logging

WentFramework includes a comprehensive logging system that supports multiple storage backends and formats, plus automatic HTTP request/response logging middleware.

### Quick Start

The logging system is automatically initialized and ready to use:

```go
import "went-framework/internal/logger"

// Basic logging
log.Info("Application started")
log.Error("Something went wrong")

// Contextual logging
log.Info("User login", map[string]interface{}{
    "user_id": 123,
    "ip":      "192.168.1.100",
})
```

### Global HTTP Request/Response Logging

All HTTP requests and responses are automatically logged with detailed information:

- **Request Details**: Method, URL, headers, body, client IP, user agent
- **Response Details**: Status code, headers, body, content type
- **Performance Metrics**: Request duration, timestamp
- **Security**: Sensitive headers (Authorization, Cookie) are automatically redacted
- **Smart Filtering**: Health checks and static files are excluded to reduce noise

#### Example Log Output

```json
{
  "level": "info",
  "message": "HTTP Request",
  "context": {
    "type": "http_request",
    "request": {
      "method": "POST",
      "url": "/api/users",
      "path": "/api/users",
      "remote_addr": "192.168.1.100",
      "user_agent": "curl/7.64.1",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": "[REDACTED]"
      },
      "body": {
        "name": "John Doe",
        "email": "john@example.com"
      }
    },
    "response": {
      "status_code": 201,
      "status_text": "Created",
      "headers": {
        "Content-Type": "application/json"
      },
      "body": {
        "status": "success",
        "message": "User created successfully"
      }
    },
    "duration_ms": 45,
    "duration": "45.123ms"
  }
}
```

### Configuration

Configure logging via environment variables:

```properties
# Log level (debug/info/warn/error)
LOG_LEVEL=info

# Log format (json/text)  
LOG_FORMAT=json

# Log storage (db/file/stdout)
LOG_STORAGE=stdout
```

### Storage Options

1. **Console Output** (`LOG_STORAGE=stdout`) - Best for development and containers
2. **File Storage** (`LOG_STORAGE=file`) - Daily rotating files in `logs/` directory  
3. **Database Storage** (`LOG_STORAGE=db`) - Searchable logs in PostgreSQL

### API Endpoints

#### View Recent Logs

```http
GET /api/logs
```

Query parameters:
- `limit` - Number of logs to retrieve (default: 100)
- `level` - Filter by log level (debug/info/warn/error)

Example:
```bash
# Get last 50 error logs
curl "http://localhost:3003/api/logs?limit=50&level=error"
```

### Features

- ✅ Multiple log levels (DEBUG, INFO, WARN, ERROR)
- ✅ Structured logging with context
- ✅ **Automatic HTTP request/response logging**
- ✅ **Global middleware integration**
- ✅ **Sensitive data redaction**
- ✅ **Performance tracking**
- ✅ **Smart filtering for noise reduction**
- ✅ Database query performance tracking
- ✅ Automatic fallback if primary storage fails
- ✅ JSON and text output formats
- ✅ RESTful log retrieval API

For detailed documentation, examples, and best practices, see [docs/LOG.md](docs/LOG.md).

## Deployment

WentFramework supports multiple deployment options for different environments and use cases.

### Docker Deployment

#### Using Docker Compose (Recommended for Development)

```bash
# Start PostgreSQL and optionally the app
docker-compose up -d

# Build and run the application
docker build -t wentframework:latest .
docker run -p 3000:3000 --env-file .env wentframework:latest serve
```

#### Standalone Docker

```bash
# Build the image
docker build -t wentframework:latest .

# Run with environment variables
docker run -p 3000:3000 \
  -e DB_HOST=your-db-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  wentframework:latest serve
```

### Kubernetes Deployment

WentFramework includes complete Kubernetes manifests for production deployment.

#### Quick Kubernetes Deployment

```bash
# Deploy everything with one command
./k8s/deploy.sh

# Access the application
kubectl port-forward service/wentframework-service 8080:80 -n wentframework
```

#### What's Included

- **Namespace isolation** - Dedicated `wentframework` namespace
- **PostgreSQL database** - Persistent storage with automatic backups
- **Auto-scaling** - Horizontal Pod Autoscaler based on CPU/Memory
- **Load balancing** - LoadBalancer service for external access
- **Ingress support** - Domain-based routing with TLS ready
- **Health checks** - Liveness and readiness probes
- **Auto-migration** - Database migrations run automatically on startup

#### Production Features

- **3 replica minimum** with auto-scaling up to 10 pods
- **Resource limits** and requests for optimal scheduling
- **Persistent storage** for PostgreSQL data
- **ConfigMap** and **Secrets** for configuration management
- **Rolling updates** with zero-downtime deployments

See the [Kubernetes README](k8s/README.md) for detailed deployment instructions.

### Cloud Deployment

#### Google Cloud Platform (GKE)
```bash
# Create GKE cluster
gcloud container clusters create wentframework-cluster

# Deploy
./k8s/deploy.sh production
```

#### Amazon Web Services (EKS)
```bash
# Create EKS cluster
eksctl create cluster --name wentframework-cluster

# Deploy
./k8s/deploy.sh production
```

#### Microsoft Azure (AKS)
```bash
# Create AKS cluster
az aks create --resource-group myRG --name wentframework-cluster

# Deploy
./k8s/deploy.sh production
```

## Development

### Architecture Overview

WentFramework follows a clean architecture pattern:

1. **`main.go`** - Entry point that handles CLI argument parsing and routes to appropriate commands
2. **`command.go`** - Contains all command implementations (server, database, swagger, code generation)
3. **`app/`** - Core application logic (models, controllers, database)
4. **`internal/`** - Internal packages for commands, swagger generation, templates, and middleware
5. **`router/`** - HTTP routing configuration with global middleware integration

### Middleware System

WentFramework includes a comprehensive middleware system:

#### Available Middleware

- **`LoggingMiddleware`** - Automatic HTTP request/response logging
- **`CORSMiddleware`** - Cross-Origin Resource Sharing support
- **`RequestIDMiddleware`** - Unique request identifier tracking

#### Global Middleware Integration

All middleware is automatically applied to all routes through the router configuration:

```go
// In router/router.go
router := mux.NewRouter()

// Apply global middleware
router.Use(middleware.CORSMiddleware)
router.Use(middleware.RequestIDMiddleware)
router.Use(middleware.LoggingMiddleware)
```

#### Custom Middleware

To add custom middleware, create a new function in `internal/middleware/`:

```go
func CustomMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Custom logic before request
        next.ServeHTTP(w, r)
        // Custom logic after request
    })
}
```

Then add it to the router configuration.

### Adding New Models

1. Generate the model using the template:
   ```bash
   go run . make:model Product
   ```

2. The generated files will be:
   - `app/models/Product.go` - With GORM integration and CRUD methods
   - `app/controllers/ProductController.go` - With HTTP handlers

3. Add routes in `router/router.go`

4. Run migrations to create the database table:
   ```bash
   go run . migrate
   ```

### Code Generation Templates

The framework uses Go templates for code generation:

- **`internal/templates/model.tpl`** - GORM-based model template
- **`internal/templates/controller.tpl`** - HTTP controller template

These templates generate code that follows the same patterns as the User model and controller.

### Database Operations

All models include standard GORM operations:

- `Create(db *gorm.DB) error` - Create a new record
- `Update(db *gorm.DB) error` - Update an existing record  
- `Delete(db *gorm.DB) error` - Delete a record
- `GetAll[Model]s(db *gorm.DB) ([]Model, error)` - Get all records
- `Get[Model]ByID(db *gorm.DB, id uint) (*Model, error)` - Get by ID
- `TableName() string` - Custom table naming

### Environment-Specific Configuration

For different environments, create separate `.env` files:

- `.env.development`
- `.env.staging`
- `.env.production`

### Testing

Test your database connection:

```bash
go run . db:test
```

This will validate your database configuration and connection.

## Recent Updates

### v0.2.1 - Logging & Middleware Improvements (August 2025)
- 📝 Logging system refactored for better performance and reliability
- 🔄 Global request/response logging middleware now wraps all routes
- 🛡️ Sensitive data redaction and smart filtering improved
- 🗄️ Database logging table migration issues fixed
- 🧩 Middleware chain supports CORS, RequestID, and custom middlewares
- 🛠️ All built-in log usages replaced with wentlog package
- 📖 README and LOG.md updated with new features and usage examples
- 🐞 Bug fixes and codebase cleanup

### v0.2.0 - Global Request/Response Logging
- ✅ Global HTTP Middleware - Automatic request/response logging for all endpoints
- ✅ Comprehensive Request Tracking - Method, URL, headers, body, client IP, user agent
- ✅ Detailed Response Logging - Status codes, headers, body, content type
- ✅ Performance Monitoring - Request duration and timing metrics
- ✅ Security Features - Automatic redaction of sensitive headers (Authorization, Cookie)
- ✅ Smart Filtering - Excludes health checks and static files to reduce log noise
- ✅ Log Retrieval API - RESTful endpoint to query and filter logs
- ✅ CORS Support - Built-in CORS middleware for frontend integration
- ✅ Request ID Tracking - Unique request identifiers for tracing

### v0.1.1 - Command Refactoring
- ✅ Separated command functions from `main.go` into `command.go`
- ✅ Cleaner main.go - Now only handles CLI routing
- ✅ Improved maintainability - Command implementations are organized
- ✅ Fixed model template - Updated to use GORM instead of raw SQL

### Model Template Improvements

The model template (`internal/templates/model.tpl`) has been updated to:
- Use GORM instead of raw SQL queries
- Include proper struct tags for JSON and GORM
- Follow the same patterns as the User model
- Generate CRUD methods that work with GORM

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

1. Check the [existing issues](https://github.com/ekinata/WentFramework/issues)
2. Create a new issue with detailed information
3. Provide your environment details and the specific issue

---

## CHANGES

- Gorilla mux to go-chi/chi/v5

### Reasons:
| Kriter                            | Gorilla Mux                                                                        | Go-Chi v5                                                                        |
| --------------------------------- | ---------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| **Durum**                         | Archived (güncellemeler durdu) ([LogRocket Blog][1], [GitHub][2])                  | Aktif geliştirme süreci devam ediyor ([GitHub][3])                               |
| **Özellik Esnekliği**             | Header, regex, route reversing, subrouter                                          | Middleware grupları, inline middleware, subrouter                                |
| **Performans & Kaynak Kullanımı** | Yüksek rota sayısında daha fazla yük ve gecikme ([GitHub][4], [Aprenda Golang][5]) | Daha düşük bellek/doğrudan kullanım, daha fazla throughput ([Aprenda Golang][5]) |
| **Kod Yazım Kolaylığı**           | Klasik, daha manuel bir yapı                                                       | Modüler ve okunaklı yapı ([GitHub][2], [Medium][6])                              |
| **Topluluk & Gelecek**            | Durmuş bir projeye yatırım yapmak… şüpheli olabilir.                               | Güncel ve geleceğe yön veren bir topluluk.                                       |

[1]: https://blog.logrocket.com/routing-go-gorilla-mux/?utm_source=chatgpt.com "An intro to routing in Go with Gorilla Mux"
[2]: https://github.com/go-chi/chi?utm_source=chatgpt.com "go-chi/chi: lightweight, idiomatic and composable router for ..."
[3]: https://github.com/go-chi/chi/blob/master/CHANGELOG.md?utm_source=chatgpt.com "chi/CHANGELOG.md at master · go-chi/chi"
[4]: https://github.com/cypriss/golang-mux-benchmark?utm_source=chatgpt.com "cypriss/golang-mux-benchmark: Performance shootout of ..."
[5]: https://aprendagolang.com.br/benchmark-dos-routers-http-chi-vs-gorilla-mux/?utm_source=chatgpt.com "Benchmark dos routers http: chi vs gorilla mux"
[6]: https://medium.com/%40hasanshahjahan/a-deep-dive-into-gin-chi-and-mux-in-go-33b9ad861e1b?utm_source=chatgpt.com "A Deep Dive into Gin, Chi, and Mux in Go"

---


**Happy coding with WentFramework! 🚀**

*Built with ❤️ for developers who want to go fast and build great APIs.*
