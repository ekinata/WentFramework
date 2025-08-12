package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"went-framework/app/database"
)

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	DEBUG LogLevel = "debug"
	INFO  LogLevel = "info"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
)

// LogEntry represents a single log entry
type LogEntry struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Level     string    `json:"level" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	Context   string    `json:"context" gorm:"type:text"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for GORM
func (LogEntry) TableName() string {
	return "logs"
}

// Logger is the main logging struct
type Logger struct {
	level   LogLevel
	format  string
	storage string
	writer  io.Writer
}

var (
	// Global logger instance
	GlobalLogger *Logger
)

// Init initializes the global logger with environment configuration
func Init() {
	level := strings.ToLower(getEnv("LOG_LEVEL", "info"))
	format := strings.ToLower(getEnv("LOG_FORMAT", "json"))
	storage := strings.ToLower(getEnv("LOG_STORAGE", "stdout"))

	GlobalLogger = NewLogger(LogLevel(level), format, storage)
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, format, storage string) *Logger {
	logger := &Logger{
		level:   level,
		format:  format,
		storage: storage,
	}

	switch storage {
	case "file":
		logger.initFileWriter()
	case "db":
		logger.initDatabaseWriter()
	default: // stdout
		logger.writer = os.Stdout
	}

	return logger
}

// initFileWriter sets up file-based logging
func (l *Logger) initFileWriter() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		l.writer = os.Stdout
		return
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("wentframework-%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		l.writer = os.Stdout
		return
	}

	l.writer = file
}

// initDatabaseWriter sets up database logging
func (l *Logger) initDatabaseWriter() {
	// Ensure the logs table exists
	if database.DB != nil {
		database.DB.AutoMigrate(&LogEntry{})
	}
}

// shouldLog determines if a message should be logged based on level
func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
	}

	return levels[level] >= levels[l.level]
}

// log writes a log entry
func (l *Logger) log(level LogLevel, message string, context map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Level:     string(level),
		Message:   message,
		Timestamp: time.Now(),
	}

	// Add context if provided
	if context != nil {
		contextBytes, _ := json.Marshal(context)
		entry.Context = string(contextBytes)
	}

	switch l.storage {
	case "db":
		l.logToDatabase(entry)
	default: // file or stdout
		l.logToWriter(entry)
	}
}

// logToDatabase saves log entry to database
func (l *Logger) logToDatabase(entry LogEntry) {
	if database.DB == nil {
		// Fallback to stdout if database is not available
		l.writer = os.Stdout
		l.logToWriter(entry)
		return
	}

	if err := database.DB.Create(&entry).Error; err != nil {
		// Fallback to stdout if database write fails
		fmt.Fprintf(os.Stderr, "Failed to write log to database: %v\n", err)
		l.writer = os.Stdout
		l.logToWriter(entry)
	}
}

// logToWriter writes log entry to the configured writer
func (l *Logger) logToWriter(entry LogEntry) {
	var output string

	if l.format == "json" {
		jsonData := map[string]interface{}{
			"timestamp": entry.Timestamp.Format(time.RFC3339),
			"level":     entry.Level,
			"message":   entry.Message,
		}

		if entry.Context != "" {
			var contextData map[string]interface{}
			if err := json.Unmarshal([]byte(entry.Context), &contextData); err == nil {
				jsonData["context"] = contextData
			}
		}

		jsonBytes, _ := json.Marshal(jsonData)
		output = string(jsonBytes) + "\n"
	} else {
		// Text format
		output = fmt.Sprintf("[%s] %s: %s",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			strings.ToUpper(entry.Level),
			entry.Message)

		if entry.Context != "" {
			output += fmt.Sprintf(" | Context: %s", entry.Context)
		}
		output += "\n"
	}

	if l.writer != nil {
		l.writer.Write([]byte(output))
	}
}

// Public logging functions

// Debug logs a debug message
func Debug(message string, context ...map[string]interface{}) {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	if GlobalLogger != nil {
		GlobalLogger.log(DEBUG, message, ctx)
	}
}

// Info logs an info message
func Info(message string, context ...map[string]interface{}) {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	if GlobalLogger != nil {
		GlobalLogger.log(INFO, message, ctx)
	}
}

// Warn logs a warning message
func Warn(message string, context ...map[string]interface{}) {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	if GlobalLogger != nil {
		GlobalLogger.log(WARN, message, ctx)
	}
}

// Error logs an error message
func Error(message string, context ...map[string]interface{}) {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	if GlobalLogger != nil {
		GlobalLogger.log(ERROR, message, ctx)
	}
}

// Convenience functions with formatting

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	Debug(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}

// LogRequest logs HTTP request information
func LogRequest(method, path, userAgent string, statusCode int, duration time.Duration) {
	Info("HTTP Request", map[string]interface{}{
		"method":      method,
		"path":        path,
		"user_agent":  userAgent,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	})
}

// LogDatabaseQuery logs database query information
func LogDatabaseQuery(query string, duration time.Duration, err error) {
	context := map[string]interface{}{
		"query":       query,
		"duration_ms": duration.Milliseconds(),
	}

	if err != nil {
		context["error"] = err.Error()
		Error("Database Query Failed", context)
	} else {
		Debug("Database Query", context)
	}
}

// GetLogs retrieves logs from database (only works when LOG_STORAGE=db)
func GetLogs(limit int, level LogLevel) ([]LogEntry, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database not available")
	}

	var logs []LogEntry
	query := database.DB.Order("timestamp DESC")

	if level != "" {
		query = query.Where("level = ?", string(level))
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

// getEnv helper function
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
