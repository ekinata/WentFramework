# WentFramework Logging System

WentFramework includes a comprehensive, flexible logging system that supports multiple storage backends and output formats. The logging system is designed to be simple to use while providing powerful features for production applications.

## Features

- ✅ **Multiple Log Levels** - DEBUG, INFO, WARN, ERROR
- ✅ **Multiple Storage Options** - Database, File, Console (stdout)
- ✅ **Multiple Formats** - JSON and Text
- ✅ **Context Support** - Rich contextual information
- ✅ **Global Access** - Use anywhere in your application
- ✅ **Automatic Fallback** - Graceful degradation if primary storage fails
- ✅ **HTTP Request Logging** - Built-in request/response logging
- ✅ **Database Query Logging** - Track database performance

## Configuration

Configure logging through environment variables in your `.env` file:

```properties
# Log level (debug/info/warn/error)
LOG_LEVEL=info

# Log format (json/text)
LOG_FORMAT=json

# Log storage (db/file/stdout)
LOG_STORAGE=local
```

### Log Levels

| Level | Description | When to Use |
|-------|-------------|-------------|
| `DEBUG` | Detailed information for debugging | Development, troubleshooting |
| `INFO` | General application flow | Normal operation, audit trail |
| `WARN` | Potentially harmful situations | Recoverable errors, deprecations |
| `ERROR` | Error events that don't stop the application | Errors, exceptions, failures |

### Storage Options

#### 1. Console Output (`LOG_STORAGE=stdout`)
**Best for:** Development, Docker containers, cloud platforms

```properties
LOG_STORAGE=stdout
```

Logs are written directly to stdout/console. This is the default and recommended option for:
- Local development
- Docker containers 
- Cloud platforms (AWS CloudWatch, Google Cloud Logging, etc.)
- CI/CD pipelines

#### 2. File Storage (`LOG_STORAGE=file`)
**Best for:** Traditional server deployments, log rotation needs

```properties
LOG_STORAGE=file
```

Logs are written to daily rotating files in the `logs/` directory:
- `logs/wentframework-2025-08-12.log`
- `logs/wentframework-2025-08-13.log`

#### 3. Database Storage (`LOG_STORAGE=db`)
**Best for:** Applications requiring searchable logs, audit trails

```properties
LOG_STORAGE=db
```

Logs are stored in a `logs` table in your PostgreSQL database. This enables:
- Querying logs via API
- Searching and filtering
- Long-term log retention
- Integration with your application data

## Basic Usage

### Import and Initialize

The logger is automatically initialized when your application starts. Simply import the log package:

```go
import "went-framework/internal/logger"
```

**Note:** If you're also using Go's standard `log` package, you can use an alias to avoid naming conflicts:

```go
import (
    "log"  // Standard Go log package
    wentlog "went-framework/internal/logger"  // WentFramework log package
)

// Then use:
log.Printf("Standard Go logging")
wentlog.Info("WentFramework logging")
```

### Simple Logging

```go
// Basic log messages
log.Info("Application started successfully")
log.Warn("This feature is deprecated")
log.Error("Failed to connect to external service")
log.Debug("Processing user data")

// Formatted messages
log.Infof("User %s logged in from %s", username, ipAddress)
log.Errorf("Database query failed: %v", err)
```

### Contextual Logging

Add structured context to your logs for better debugging and monitoring:

```go
// Log with context
log.Info("User login attempt", map[string]interface{}{
    "user_id":    12345,
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0...",
    "success":    true,
})

log.Error("Payment processing failed", map[string]interface{}{
    "order_id":     "ORD-123456",
    "amount":       99.99,
    "currency":     "USD",
    "error_code":   "CARD_DECLINED",
    "gateway":      "stripe",
})
```

## Advanced Usage

### HTTP Request Logging

Built-in function for logging HTTP requests:

```go
// In your HTTP middleware or handlers
start := time.Now()
// ... handle request ...
duration := time.Since(start)

log.LogRequest(
    r.Method,           // GET, POST, etc.
    r.URL.Path,         // /api/users/123
    r.UserAgent(),      // Browser/client info
    responseStatusCode, // 200, 404, 500, etc.
    duration,           // Request processing time
)
```

### Database Query Logging

Track database performance and errors:

```go
start := time.Now()
err := db.Find(&users).Error
duration := time.Since(start)

log.LogDatabaseQuery(
    "SELECT * FROM users WHERE active = true",
    duration,
    err,
)
```

### Retrieving Logs (Database Storage Only)

When using database storage, you can retrieve logs programmatically:

```go
// Get last 100 logs
logs, err := log.GetLogs(100, "")

// Get only error logs
errorLogs, err := log.GetLogs(50, log.ERROR)

// Get all logs (no limit)
allLogs, err := log.GetLogs(0, "")
```

## Log Formats

### JSON Format (`LOG_FORMAT=json`)
**Best for:** Machine processing, structured logging platforms

```json
{
  "timestamp": "2025-08-12T10:30:45Z",
  "level": "info",
  "message": "User login attempt",
  "context": {
    "user_id": 12345,
    "ip_address": "192.168.1.100",
    "success": true
  }
}
```

### Text Format (`LOG_FORMAT=text`)
**Best for:** Human reading, simple deployments

```
[2025-08-12 10:30:45] INFO: User login attempt | Context: {"user_id":12345,"ip_address":"192.168.1.100","success":true}
```

## Integration Examples

### In Controllers

```go
package controllers

import (
    "net/http"
    "time"
    "went-framework/internal/logger"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    log.Debug("GetUsers endpoint called", map[string]interface{}{
        "user_agent": r.UserAgent(),
        "ip":         r.RemoteAddr,
    })

    // Your business logic here
    users, err := getUsersFromDatabase()
    if err != nil {
        log.Error("Failed to fetch users", map[string]interface{}{
            "error": err.Error(),
            "query_time": time.Since(start).Milliseconds(),
        })
        http.Error(w, "Internal server error", 500)
        return
    }

    log.LogRequest(r.Method, r.URL.Path, r.UserAgent(), 200, time.Since(start))
    
    // Send response...
}
```

### In Models

```go
package models

import (
    "time"
    "went-framework/internal/logger"
    "gorm.io/gorm"
)

func (u *User) Create(db *gorm.DB) error {
    start := time.Now()
    
    err := db.Create(u).Error
    duration := time.Since(start)
    
    if err != nil {
        log.LogDatabaseQuery("INSERT INTO users", duration, err)
        log.Error("User creation failed", map[string]interface{}{
            "email": u.Email,
            "error": err.Error(),
        })
        return err
    }
    
    log.LogDatabaseQuery("INSERT INTO users", duration, nil)
    log.Info("User created successfully", map[string]interface{}{
        "user_id": u.ID,
        "email":   u.Email,
    })
    
    return nil
}
```

### In Middleware

```go
package middleware

import (
    "net/http"
    "time"
    "went-framework/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Log incoming request
        log.Debug("Incoming request", map[string]interface{}{
            "method":     r.Method,
            "path":       r.URL.Path,
            "user_agent": r.UserAgent(),
            "ip":         r.RemoteAddr,
        })
        
        // Create a response writer wrapper to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        // Process request
        next.ServeHTTP(wrapped, r)
        
        // Log request completion
        duration := time.Since(start)
        log.LogRequest(r.Method, r.URL.Path, r.UserAgent(), wrapped.statusCode, duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

## Production Considerations

### Log Rotation (File Storage)

For file storage, implement log rotation:

```bash
# Add to crontab for daily rotation
0 0 * * * /usr/sbin/logrotate /path/to/logrotate.conf
```

Create `/etc/logrotate.d/wentframework`:
```
/path/to/your/app/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 app app
}
```

### Database Storage Considerations

When using database storage:

1. **Performance**: Logging is synchronous and may impact performance
2. **Storage**: Logs can grow large quickly
3. **Cleanup**: Implement log retention policies

```sql
-- Example: Delete logs older than 30 days
DELETE FROM logs WHERE created_at < NOW() - INTERVAL '30 days';
```

### Cloud Integration

#### AWS CloudWatch (with stdout)
```properties
LOG_STORAGE=stdout
LOG_FORMAT=json
```

#### Google Cloud Logging (with stdout)
```properties
LOG_STORAGE=stdout
LOG_FORMAT=json
```

#### Azure Monitor (with stdout)
```properties
LOG_STORAGE=stdout
LOG_FORMAT=json
```

## Best Practices

### 1. Choose Appropriate Log Levels
```go
// ✅ Good
handlers.Debug("Processing user input", context)  // Development info
handlers.Info("User created", context)            // Business events
handlers.Warn("Rate limit exceeded", context)     // Potential issues
handlers.Error("Database connection failed", context) // Actual errors

// ❌ Avoid
handlers.Info("Variable x = 5")                   // Too verbose for INFO
handlers.Error("User not found")                  // Not really an error
```

### 2. Include Relevant Context
```go
// ✅ Good - Rich context
handlers.Error("Payment failed", map[string]interface{}{
    "order_id":   order.ID,
    "user_id":    user.ID,
    "amount":     order.Total,
    "gateway":    "stripe",
    "error_code": result.ErrorCode,
})

// ❌ Avoid - No context
handlers.Error("Payment failed")
```

### 3. Don't Log Sensitive Information
```go
// ✅ Good
handlers.Info("User login", map[string]interface{}{
    "user_id": user.ID,
    "email":   user.Email,
})

// ❌ Avoid - Sensitive data
handlers.Info("User login", map[string]interface{}{
    "password":    user.Password,    // Never log passwords
    "credit_card": user.CreditCard,  // Never log PII
})
```

### 4. Use Structured Logging
```go
// ✅ Good - Structured
handlers.Info("Order processed", map[string]interface{}{
    "order_id":     order.ID,
    "total_amount": order.Total,
    "item_count":   len(order.Items),
})

// ❌ Avoid - Unstructured
handlers.Infof("Order %d processed with %d items for $%.2f", 
    order.ID, len(order.Items), order.Total)
```

## Troubleshooting

### Logger Not Working

1. **Check initialization**:
   ```go
   log.Init() // Must be called after loading env vars
   ```

2. **Check log level**:
   ```properties
   LOG_LEVEL=debug  # Make sure level allows your messages
   ```

3. **Check file permissions** (file storage):
   ```bash
   ls -la logs/
   ```

### Database Logging Issues

1. **Database connection**: Ensure database is connected before logging
2. **Table creation**: The `logs` table is auto-created, but check for migration errors
3. **Fallback**: Logger automatically falls back to stdout if database fails

### Performance Issues

1. **Database logging**: Consider switching to file or stdout for high-traffic applications
2. **Log level**: Use higher log levels (INFO/WARN/ERROR) in production
3. **Context size**: Keep context objects reasonably sized

## Migration from Other Logging Libraries

### From standard `log` package:
```go
// Before
log.Printf("User %s logged in", username)

// After
wentlog.Infof("User %s logged in", username)
```

### From `logrus`:
```go
// Before
logrus.WithFields(logrus.Fields{
    "user_id": 123,
    "action":  "login",
}).Info("User action")

// After
log.Info("User action", map[string]interface{}{
    "user_id": 123,
    "action":  "login",
})
```

---

## Quick Reference

### Environment Variables
```properties
LOG_LEVEL=info      # debug/info/warn/error
LOG_FORMAT=json     # json/text
LOG_STORAGE=stdout  # db/file/stdout
```

### Functions
```go
// Basic logging
log.Debug(message)
log.Info(message)
log.Warn(message)
log.Error(message)

// Formatted logging
log.Debugf(format, args...)
log.Infof(format, args...)
log.Warnf(format, args...)
log.Errorf(format, args...)

// Contextual logging
log.Info(message, map[string]interface{}{...})

// Specialized logging
log.LogRequest(method, path, userAgent, statusCode, duration)
log.LogDatabaseQuery(query, duration, error)

// Retrieve logs (database only)
log.GetLogs(limit, level)
```

**Happy logging! 📝**
