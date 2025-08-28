package bank

import "time"

// TransactionRequest represents a request to create a bank transaction
type TransactionRequest struct {
	FromAccountNum string `json:"fromAccountNum"`
	FromRoutingNum string `json:"fromRoutingNum"`
	ToAccountNum   string `json:"toAccountNum"`
	ToRoutingNum   string `json:"toRoutingNum"`
	Amount         int64  `json:"amount"`
	UUID           string `json:"uuid"`
}

// TransactionResponse represents the response from creating a bank transaction
type TransactionResponse struct {
	TransactionID  int64     `json:"transactionId"`
	FromAccountNum string    `json:"fromAccountNum"`
	FromRoutingNum string    `json:"fromRoutingNum"`
	ToAccountNum   string    `json:"toAccountNum"`
	ToRoutingNum   string    `json:"toRoutingNum"`
	Amount         int64     `json:"amount"`
	Timestamp      time.Time `json:"timestamp"`
}

// BalanceResponse represents an account balance response
type BalanceResponse struct {
	AccountNum string `json:"accountNum"`
	RoutingNum string `json:"routingNum"`
	Balance    int64  `json:"balance"`
}

// ErrorResponse represents an error from the Bank API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
