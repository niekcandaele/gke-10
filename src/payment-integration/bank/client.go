package bank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Authenticator interface for JWT token generation
type Authenticator interface {
	GetAuthHeader(accountNumber string) (string, error)
}

// Client represents a Bank of Anthos API client
type Client struct {
	baseURL       string
	httpClient    *http.Client
	authenticator Authenticator
}

// NewClient creates a new Bank API client
func NewClient(baseURL string, authenticator Authenticator) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authenticator: authenticator,
	}
}

// CreateTransaction creates a new bank transaction
func (c *Client) CreateTransaction(req *TransactionRequest) (*TransactionResponse, error) {
	url := fmt.Sprintf("%s/transactions", c.baseURL)

	// Marshal request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Creating bank transaction: from=%s to=%s amount=%d uuid=%s",
		req.FromAccountNum, req.ToAccountNum, req.Amount, req.UUID)

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Add JWT authentication if available
	if c.authenticator != nil {
		// Use the sender's account number for authentication
		authHeader, err := c.authenticator.GetAuthHeader(req.FromAccountNum)
		if err != nil {
			return nil, fmt.Errorf("failed to generate auth header: %w", err)
		}
		httpReq.Header.Set("Authorization", authHeader)
		log.Printf("Added JWT auth header for account %s", req.FromAccountNum)
	} else {
		log.Printf("WARNING: No authenticator available, request will likely fail")
	}

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return nil, NewBankError(resp.StatusCode, errResp.Error, errResp.Message)
		}
		return nil, NewBankError(resp.StatusCode, "transaction_failed", string(respBody))
	}

	// Parse successful response
	var txResponse TransactionResponse
	if err := json.Unmarshal(respBody, &txResponse); err != nil {
		// Bank API might return empty body on 201, so we'll construct a response
		if resp.StatusCode == http.StatusCreated {
			log.Printf("Bank API returned 201 with no body, transaction successful")
			return &TransactionResponse{
				FromAccountNum: req.FromAccountNum,
				FromRoutingNum: req.FromRoutingNum,
				ToAccountNum:   req.ToAccountNum,
				ToRoutingNum:   req.ToRoutingNum,
				Amount:         req.Amount,
				Timestamp:      time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("Bank transaction created successfully: ID=%d", txResponse.TransactionID)
	return &txResponse, nil
}

// CheckBalance checks the balance of an account (for future use)
func (c *Client) CheckBalance(accountNum, routingNum string) (*BalanceResponse, error) {
	url := fmt.Sprintf("%s/balances/%s", c.baseURL, accountNum)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add JWT authentication if available
	if c.authenticator != nil {
		authHeader, err := c.authenticator.GetAuthHeader(accountNum)
		if err != nil {
			return nil, fmt.Errorf("failed to generate auth header: %w", err)
		}
		httpReq.Header.Set("Authorization", authHeader)
	}

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return nil, NewBankError(resp.StatusCode, errResp.Error, errResp.Message)
		}
		return nil, NewBankError(resp.StatusCode, "balance_check_failed", string(respBody))
	}

	// Parse successful response
	var balanceResp BalanceResponse
	if err := json.Unmarshal(respBody, &balanceResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &balanceResp, nil
}

// HealthCheck checks if the Bank API is available
func (c *Client) HealthCheck() error {
	url := fmt.Sprintf("%s/ready", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("bank API unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bank API not ready, status: %d", resp.StatusCode)
	}

	return nil
}
