package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

// Logger provides structured logging capabilities
type Logger struct {
	serviceName string
	output      *log.Logger
	minLevel    LogLevel
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp      string                 `json:"timestamp"`
	Level          LogLevel               `json:"level"`
	Service        string                 `json:"service"`
	TransactionID  string                 `json:"transaction_id,omitempty"`
	CorrelationID  string                 `json:"correlation_id,omitempty"`
	AccountNumber  string                 `json:"account_number,omitempty"`
	Message        string                 `json:"message"`
	Error          string                 `json:"error,omitempty"`
	Amount         int64                  `json:"amount,omitempty"`
	Currency       string                 `json:"currency,omitempty"`
	FromAccount    string                 `json:"from_account,omitempty"`
	ToAccount      string                 `json:"to_account,omitempty"`
	Duration       float64                `json:"duration_ms,omitempty"`
	AdditionalData map[string]interface{} `json:"data,omitempty"`
}

// NewLogger creates a new structured logger
func NewLogger(serviceName string) *Logger {
	minLevel := INFO
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		minLevel = DEBUG
	}

	return &Logger{
		serviceName: serviceName,
		output:      log.New(os.Stdout, "", 0),
		minLevel:    minLevel,
	}
}

// shouldLog checks if the message should be logged based on level
func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
		FATAL: 4,
	}
	return levels[level] >= levels[l.minLevel]
}

// log writes a log entry
func (l *Logger) log(level LogLevel, entry LogEntry) {
	if !l.shouldLog(level) {
		return
	}

	entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	entry.Level = level
	entry.Service = l.serviceName

	jsonData, err := json.Marshal(entry)
	if err != nil {
		l.output.Printf(`{"timestamp":"%s","level":"ERROR","service":"%s","message":"Failed to marshal log entry","error":"%s"}`,
			time.Now().UTC().Format(time.RFC3339), l.serviceName, err.Error())
		return
	}

	l.output.Println(string(jsonData))

	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, data map[string]interface{}) {
	l.log(DEBUG, LogEntry{Message: message, AdditionalData: data})
}

// Info logs an info message
func (l *Logger) Info(message string, data map[string]interface{}) {
	l.log(INFO, LogEntry{Message: message, AdditionalData: data})
}

// Warn logs a warning message
func (l *Logger) Warn(message string, data map[string]interface{}) {
	l.log(WARN, LogEntry{Message: message, AdditionalData: data})
}

// Error logs an error message
func (l *Logger) Error(message string, err error, data map[string]interface{}) {
	entry := LogEntry{Message: message, AdditionalData: data}
	if err != nil {
		entry.Error = err.Error()
	}
	l.log(ERROR, entry)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, err error) {
	entry := LogEntry{Message: message}
	if err != nil {
		entry.Error = err.Error()
	}
	l.log(FATAL, entry)
}

// LogTransaction logs a payment transaction
func (l *Logger) LogTransaction(txID string, fromAcct string, toAcct string, amount int64, currency string, message string) {
	l.log(INFO, LogEntry{
		TransactionID: txID,
		FromAccount:   fromAcct,
		ToAccount:     toAcct,
		Amount:        amount,
		Currency:      currency,
		Message:       message,
	})
}

// LogPaymentRequest logs an incoming payment request
func (l *Logger) LogPaymentRequest(ctx context.Context, txID string, amount int64, currency string, cardLast4 string) {
	l.log(INFO, LogEntry{
		TransactionID: txID,
		Amount:        amount,
		Currency:      currency,
		Message:       fmt.Sprintf("Payment request received for card ****%s", cardLast4),
	})
}

// LogPaymentResponse logs a payment response
func (l *Logger) LogPaymentResponse(ctx context.Context, txID string, success bool, duration time.Duration, err error) {
	entry := LogEntry{
		TransactionID: txID,
		Duration:      duration.Seconds() * 1000, // Convert to milliseconds
	}

	if success {
		entry.Message = "Payment processed successfully"
		l.log(INFO, entry)
	} else {
		entry.Message = "Payment failed"
		if err != nil {
			entry.Error = err.Error()
		}
		l.log(ERROR, entry)
	}
}

// LogBankAPICall logs calls to the Bank API
func (l *Logger) LogBankAPICall(txID string, accountNum string, amount int64, duration time.Duration, err error) {
	entry := LogEntry{
		TransactionID: txID,
		AccountNumber: accountNum,
		Amount:        amount,
		Duration:      duration.Seconds() * 1000,
	}

	if err == nil {
		entry.Message = "Bank API call successful"
		l.log(INFO, entry)
	} else {
		entry.Message = "Bank API call failed"
		entry.Error = err.Error()
		l.log(ERROR, entry)
	}
}

// ExtractCorrelationID extracts correlation ID from context
func ExtractCorrelationID(ctx context.Context) string {
	if val := ctx.Value("correlation_id"); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}
