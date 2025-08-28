package bank

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check path
		if r.URL.Path != "/transactions" {
			t.Errorf("Expected path /transactions, got %s", r.URL.Path)
		}

		// Check method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || authHeader[:7] != "Bearer " {
			t.Errorf("Missing or invalid Authorization header: %s", authHeader)
		}

		// Parse request body
		var req TransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Validate request fields
		if req.FromAccountNum == "" || req.ToAccountNum == "" {
			t.Errorf("Missing account numbers")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusCreated)
		resp := TransactionResponse{
			TransactionID:  12345,
			FromAccountNum: req.FromAccountNum,
			ToAccountNum:   req.ToAccountNum,
			Amount:         req.Amount,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client with mock authenticator
	client := NewClient(server.URL, &mockAuthenticator{})

	// Test request
	req := &TransactionRequest{
		FromAccountNum: "1234567890",
		FromRoutingNum: "123456789",
		ToAccountNum:   "9999999999",
		ToRoutingNum:   "123456789",
		Amount:         10000,
		UUID:           "test-uuid-123",
	}

	resp, err := client.CreateTransaction(req)
	if err != nil {
		t.Fatalf("CreateTransaction failed: %v", err)
	}

	if resp.TransactionID != 12345 {
		t.Errorf("Expected transaction ID 12345, got %d", resp.TransactionID)
	}
}

func TestCreateTransactionInsufficientFunds(t *testing.T) {
	// Create mock server that returns insufficient funds error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		errResp := ErrorResponse{
			Error:   "insufficient_funds",
			Message: "Account has insufficient funds",
		}
		json.NewEncoder(w).Encode(errResp)
	}))
	defer server.Close()

	client := NewClient(server.URL, &mockAuthenticator{})

	req := &TransactionRequest{
		FromAccountNum: "1234567890",
		FromRoutingNum: "123456789",
		ToAccountNum:   "9999999999",
		ToRoutingNum:   "123456789",
		Amount:         1000000, // Large amount
		UUID:           "test-uuid-456",
	}

	_, err := client.CreateTransaction(req)
	if err == nil {
		t.Fatal("Expected error for insufficient funds, got nil")
	}

	bankErr, ok := err.(*BankError)
	if !ok {
		t.Fatalf("Expected BankError, got %T", err)
	}

	if !bankErr.IsInsufficientFunds() {
		t.Errorf("Expected insufficient funds error, got %s", bankErr.ErrorCode)
	}
}

func TestCreateTransactionDuplicate(t *testing.T) {
	// Create mock server that returns duplicate transaction error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		errResp := ErrorResponse{
			Error:   "duplicate_transaction",
			Message: "Transaction with UUID already exists",
		}
		json.NewEncoder(w).Encode(errResp)
	}))
	defer server.Close()

	client := NewClient(server.URL, &mockAuthenticator{})

	req := &TransactionRequest{
		FromAccountNum: "1234567890",
		FromRoutingNum: "123456789",
		ToAccountNum:   "9999999999",
		ToRoutingNum:   "123456789",
		Amount:         5000,
		UUID:           "duplicate-uuid",
	}

	_, err := client.CreateTransaction(req)
	if err == nil {
		t.Fatal("Expected error for duplicate transaction, got nil")
	}

	bankErr, ok := err.(*BankError)
	if !ok {
		t.Fatalf("Expected BankError, got %T", err)
	}

	if !bankErr.IsDuplicateTransaction() {
		t.Errorf("Expected duplicate transaction error, got %s", bankErr.ErrorCode)
	}
}

func TestHealthCheck(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ready" {
			t.Errorf("Expected path /ready, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)

	if err := client.HealthCheck(); err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

func TestHealthCheckFailed(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)

	if err := client.HealthCheck(); err == nil {
		t.Error("Expected error for failed health check, got nil")
	}
}

// mockAuthenticator for testing
type mockAuthenticator struct{}

func (m *mockAuthenticator) GetAuthHeader(accountNumber string) (string, error) {
	return "Bearer mock-jwt-token", nil
}
