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
- 🔍 **CRUD Operations** with JSON API responses
- 📚 **Auto-Swagger Generation** with interactive API documentation

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [API Usage](#api-usage)
- [API Documentation](#api-documentation)
- [Commands](#commands)
- [Project Structure](#project-structure)
- [Development](#development)

## Installation

### Prerequisites

- Go 1.24.4 or higher
- Docker and Docker Compose (for PostgreSQL)
- Git

### Clone the Repository

```bash
git clone <your-repo-url>
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

This will start a PostgreSQL container with the following configuration:
- Host: `localhost`
- Port: `5432`
- Database: `went_test`
- User: `went_user`
- Password: `went_password`

### 2. Set Up Environment

Copy the example environment file:

```bash
cp .env.example .env
```

The default `.env` file is already configured to work with the Docker PostgreSQL setup.

### 3. Run Database Migrations

```bash
go run main.go migrate
```

### 4. Start the Server

```bash
go run main.go serve
```

The server will start on `http://localhost:3003` (or the port specified in your `.env` file).

### 5. Generate API Documentation

```bash
go run main.go swagger:generate
```

This will generate interactive Swagger documentation for your API.

### 6. Test the API

```bash
# Health check
curl http://localhost:3003/api/health

# Get all users
curl http://localhost:3003/api/users

# Create a user
curl -X POST http://localhost:3003/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

### 7. View API Documentation

Once the server is running, you can access the interactive API documentation:

- **Swagger UI**: `http://localhost:3003/swagger/`
- **OpenAPI JSON**: `http://localhost:3003/swagger.json`

## Configuration

### Environment Variables

The project uses environment variables for configuration. Copy `.env.example` to `.env` and modify as needed:

```properties
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=went_user
DB_PASSWORD=went_password
DB_NAME=went_test
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=3003
SERVER_HOST=0.0.0.0

# Application Configuration
APP_ENV=development
APP_NAME=WentFramework
APP_VERSION=1.0.0

# JWT Configuration (for future use)
JWT_SECRET=your-super-secure-jwt-secret-key-change-this-in-production
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## Database Setup

### Using Docker (Recommended)

```bash
# Start PostgreSQL container
docker-compose up -d

# Check if container is running
docker-compose ps

# Stop the container
docker-compose down
```

### Manual PostgreSQL Setup

If you prefer to use a local PostgreSQL installation:

1. Install PostgreSQL
2. Create a database and user
3. Update the `.env` file with your database credentials

### Testing Database Connection

```bash
go run main.go db:test
```

This command will test the database connection and display the configuration being used.

## API Usage

### Base URL

```
http://localhost:3003/api
```

### Authentication

Currently, the API doesn't require authentication. JWT support is planned for future releases.

### Endpoints

#### Health Check

```http
GET /api/health
```

**Response:**
```json
{
  "status": "healthy",
  "message": "Server is running"
}
```

#### Users

##### Get All Users

```http
GET /api/users
```

**Response:**
```json
{
  "status": "success",
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2025-07-31T15:42:18.792477+03:00",
      "updated_at": "2025-07-31T15:42:18.792477+03:00"
    }
  ]
}
```

##### Get User by ID

```http
GET /api/users/{id}
```

**Response:**
```json
{
  "status": "success",
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-07-31T15:42:18.792477+03:00",
    "updated_at": "2025-07-31T15:42:18.792477+03:00"
  }
}
```

##### Create User

```http
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-07-31T15:42:18.792477+03:00",
    "updated_at": "2025-07-31T15:42:18.792477+03:00"
  }
}
```

##### Update User

```http
PUT /api/users/{id}
Content-Type: application/json

{
  "name": "John Updated",
  "email": "john.updated@example.com"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "name": "John Updated",
    "email": "john.updated@example.com",
    "created_at": "2025-07-31T15:42:18.792477+03:00",
    "updated_at": "2025-07-31T15:43:10.469621+03:00"
  }
}
```

##### Delete User

```http
DELETE /api/users/{id}
```

**Response:**
```json
{
  "status": "success",
  "message": "User 1 deleted successfully"
}
```

### Error Responses

All error responses follow this format:

```json
{
  "status": "error",
  "message": "Error description"
}
```

Common HTTP status codes:
- `400 Bad Request` - Invalid input data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## API Documentation

WentFramework provides **automatic API documentation** through Swagger/OpenAPI integration. The documentation is generated automatically from your models, controllers, and routes.

### Interactive Swagger UI

Access the interactive API documentation at:
```
http://localhost:3003/swagger/
```

The Swagger UI provides:
- 📖 **Complete API Reference** - All endpoints with detailed descriptions
- 🧪 **Interactive Testing** - Test APIs directly from the browser
- 📝 **Request/Response Examples** - Sample data for all endpoints
- 🔍 **Schema Documentation** - Complete data model definitions
- 📋 **Parameter Details** - Required and optional parameters for each endpoint

### OpenAPI Specification

Get the machine-readable OpenAPI 3.0 specification:
```
http://localhost:3003/swagger.json
```

### Auto-Generation Features

The Swagger documentation is automatically generated using reflection and includes:

- **Model Schemas** - Automatically extracted from Go structs in `models/`
- **API Endpoints** - Dynamically discovered from your router configuration  
- **Request/Response Types** - Inferred from controller function signatures
- **Parameter Validation** - Automatically documented from route parameters
- **Example Data** - Generated examples for all data types

### Generating Documentation

Update the API documentation after making changes:

```bash
# Generate/update Swagger documentation
go run main.go swagger:generate

# Start server with updated docs
go run main.go serve
```

The documentation will be automatically served when you start the server.

### Customizing Documentation

The auto-generated documentation can be customized by:

1. **Adding struct tags** to your models for better examples
2. **Using descriptive function names** in controllers
3. **Following RESTful naming conventions** for automatic categorization

Example model with documentation tags:
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
    Name      string    `json:"name" gorm:"not null" example:"John Doe"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null" example:"john@example.com"`
    CreatedAt time.Time `json:"created_at" example:"2025-07-31T15:42:18.792477+03:00"`
    UpdatedAt time.Time `json:"updated_at" example:"2025-07-31T15:42:18.792477+03:00"`
}
```

## Commands

### Server Commands

```bash
# Start the HTTP server
go run main.go serve

# Test database connection
go run main.go db:test
```

### Documentation Commands

```bash
# Generate/update Swagger API documentation
go run main.go swagger:generate
```

### Migration Commands

```bash
# Run migrations (create/update tables)
go run main.go migrate

# Fresh migration (drop and recreate tables)
go run main.go migrate:fresh
```

### Code Generation Commands

```bash
# Generate model, controller, and migration files
go run main.go make:model ModelName

# Example: Create a Post model
go run main.go make:model Post
```

This will create:
- `models/Post.go` - Model file
- `controllers/PostController.go` - Controller file
- `migrations/timestamp_create_posts_table.go` - Migration file

### Help

```bash
# Show available commands
go run main.go
```

## Project Structure

```
WentFramework/
├── .env                     # Environment variables (create from .env.example)
├── .env.example            # Environment variables template
├── .gitignore              # Git ignore rules
├── docker-compose.yml      # PostgreSQL container configuration
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies
├── main.go                 # Main application entry point
├── commands/               # Command implementations
│   └── migrate.go         # Migration commands
├── controllers/            # HTTP request handlers
│   └── UserController.go  # User CRUD operations
├── database/               # Database connection and configuration
│   └── connection.go      # Database connection setup
├── docs/                   # Generated documentation
│   └── swagger.json       # Auto-generated OpenAPI specification
├── models/                 # Data models and database operations
│   └── User.go           # User model
├── router/                 # HTTP routing configuration
│   └── router.go         # Route definitions and setup
├── swagger/                # API documentation generation
│   └── generator.go       # Auto-swagger generation system
└── templates/              # Code generation templates
    ├── controller.tpl     # Controller template
    ├── migration.tpl      # Migration template
    └── model.tpl          # Model template
```

## Development

### Adding New Models

1. Generate the model:
   ```bash
   go run main.go make:model Product
   ```

2. Update the model file (`models/Product.go`) with your fields
3. Update the controller file (`controllers/ProductController.go`) with your logic
4. Add routes in `router/router.go`
5. Run migrations:
   ```bash
   go run main.go migrate
   ```

### Adding New Routes

Edit `router/router.go` and add your routes to the appropriate setup function:

```go
func setupProductRoutes(api *mux.Router) {
    api.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
    api.HandleFunc("/products/{id}", controllers.GetProduct).Methods("GET")
    // ... more routes
}
```

### Environment-Specific Configuration

For different environments (development, staging, production), create separate `.env` files:

- `.env.development`
- `.env.staging`
- `.env.production`

Load the appropriate file based on your deployment environment.

### Database Migrations

The framework uses GORM's auto-migration feature. For production environments, you might want to implement more sophisticated migration handling.

## Docker Support

### Development with Docker

```bash
# Start PostgreSQL only
docker-compose up -d

# Stop PostgreSQL
docker-compose down
```

### Full Docker Setup (Future Enhancement)

You can extend the `docker-compose.yml` to include the Go application:

```yaml
services:
  app:
    build: .
    ports:
      - "3003:3003"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please:

1. Check the [existing issues](https://github.com/your-username/WentFramework/issues)
2. Create a new issue if your problem isn't already reported
3. Provide detailed information about your environment and the issue

## Roadmap & Planned Features

WentFramework is actively being developed. Here are the upcoming features and improvements planned for future releases:

### ✅ **Completed Features**
- **Auto-Swagger Generation** - Automatic OpenAPI/Swagger documentation generation from models, controllers, and routes
- **Interactive API Documentation** - Live Swagger UI with testing capabilities
- **Reflection-based Schema Generation** - Automatic model schema extraction

### 🔒 **Middleware System**
- **Authentication Middleware** - JWT-based authentication for protected routes
- **Authorization Middleware** - Role-based access control (RBAC)
- **CORS Middleware** - Cross-origin resource sharing configuration
- **Rate Limiting** - API rate limiting to prevent abuse
- **Custom Middleware** - Easy-to-implement custom middleware interface

### 📊 **Logging & Observability**
- **Structured Logging** - JSON-formatted logs with configurable levels
- **Request/Response Logging** - Automatic HTTP request and response logging
- **OpenTelemetry Integration** - Distributed tracing and metrics collection
- **Performance Monitoring** - Request duration and performance metrics
- **Error Tracking** - Enhanced error reporting and stack traces

### ✅ **Request Validation**
- **Input Validation** - Automatic request body validation with custom rules
- **Schema Validation** - JSON schema-based validation
- **Type Safety** - Strong typing for request/response data
- **Custom Validators** - Extensible validation system
- **Sanitization** - Input sanitization to prevent XSS and injection attacks

### 📦 **Response Management**
- **Resource Classes** - Structured response formatting for single resources
- **Collection Classes** - Standardized collection responses with metadata
- **Response Transformers** - Data transformation before sending responses
- **Conditional Responses** - ETags and conditional request handling
- **Content Negotiation** - Support for multiple response formats (JSON, XML, etc.)

### 📄 **Pagination System**
- **Cursor-based Pagination** - Efficient pagination for large datasets
- **Offset-based Pagination** - Traditional page-based pagination
- **Configurable Page Sizes** - Customizable pagination parameters
- **Pagination Metadata** - Rich pagination information in responses
- **Search Integration** - Pagination combined with search and filtering

### 🚀 **Additional Planned Features**
- **API Versioning** - Support for multiple API versions
- **Background Jobs** - Queue system for asynchronous tasks
- **Caching Layer** - Redis integration for caching
- **File Upload Handling** - Multipart file upload support
- **WebSocket Support** - Real-time communication capabilities
- **Health Checks** - Advanced health monitoring endpoints
- **Database Seeders** - Automated database seeding for development
- **Testing Utilities** - Built-in testing helpers and assertions

### 📅 **Release Timeline**

| Feature | Priority | Status |
|---------|----------|--------|
| Auto-Swagger Generation | High | ✅ **Completed** |
| Middleware System | High | v0.2.0 |
| Request Validation | High | v0.2.0 |
| Logging & OpenTelemetry | High | v0.3.0 |
| Response Resources | Medium | v0.3.0 |
| Pagination Methods | Medium | v0.4.0 |
| Background Jobs | Low | v0.5.0 |
| WebSocket Support | Low | v0.6.0 |

### 🤝 **Contributing to Planned Features**

We welcome contributions to any of these planned features! If you're interested in implementing any of these features:

1. **Check the Issues** - Look for existing issues related to the feature
2. **Create a Discussion** - Start a discussion about your implementation approach
3. **Fork & Implement** - Create a fork and implement the feature
4. **Submit PR** - Submit a pull request with your implementation

### 💡 **Feature Requests**

Have an idea for a feature not listed here? We'd love to hear about it:

1. Open a [Feature Request Issue](https://github.com/ekinata/WentFramework/issues/new)
2. Describe the feature and its use case
3. Participate in the discussion with the community

---

**Happy coding with WentFramework! 🚀**
