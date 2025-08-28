package bank

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BankError represents an error from the Bank API
type BankError struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

// NewBankError creates a new BankError
func NewBankError(statusCode int, errorCode, message string) *BankError {
	return &BankError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

// Error implements the error interface
func (e *BankError) Error() string {
	return fmt.Sprintf("bank API error %d (%s): %s", e.StatusCode, e.ErrorCode, e.Message)
}

// IsInsufficientFunds checks if the error is due to insufficient funds
func (e *BankError) IsInsufficientFunds() bool {
	return e.StatusCode == http.StatusBadRequest &&
		(e.ErrorCode == "insufficient_funds" ||
			e.ErrorCode == "INSUFFICIENT_FUNDS" ||
			e.ErrorCode == "insufficient_balance")
}

// IsDuplicateTransaction checks if the error is due to a duplicate transaction
func (e *BankError) IsDuplicateTransaction() bool {
	return e.StatusCode == http.StatusConflict ||
		e.ErrorCode == "duplicate_transaction" ||
		e.ErrorCode == "DUPLICATE_TRANSACTION"
}

// IsUnauthorized checks if the error is due to authentication failure
func (e *BankError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// ToGRPCError converts a BankError to a gRPC error
func (e *BankError) ToGRPCError() error {
	switch {
	case e.IsUnauthorized():
		return status.Error(codes.Unauthenticated, "authentication with bank failed")
	case e.IsInsufficientFunds():
		return status.Error(codes.FailedPrecondition, "insufficient funds in account")
	case e.IsDuplicateTransaction():
		return status.Error(codes.AlreadyExists, "duplicate transaction")
	case e.StatusCode >= 400 && e.StatusCode < 500:
		return status.Error(codes.InvalidArgument, e.Message)
	case e.StatusCode >= 500:
		return status.Error(codes.Internal, "bank service error")
	default:
		return status.Error(codes.Unknown, e.Message)
	}
}

// HandleBankError converts a bank error to an appropriate gRPC error
func HandleBankError(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's a BankError
	if bankErr, ok := err.(*BankError); ok {
		return bankErr.ToGRPCError()
	}

	// Generic error handling
	return status.Error(codes.Internal, fmt.Sprintf("bank communication error: %v", err))
}
