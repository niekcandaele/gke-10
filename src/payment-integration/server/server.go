package server

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gke-hackathon/payment-integration/auth"
	"github.com/gke-hackathon/payment-integration/bank"
	"github.com/gke-hackathon/payment-integration/converter"
	"github.com/gke-hackathon/payment-integration/logging"
	"github.com/gke-hackathon/payment-integration/mapper"
	"github.com/gke-hackathon/payment-integration/metrics"
	"github.com/gke-hackathon/payment-integration/middleware"
	pb "github.com/gke-hackathon/payment-integration/proto"
	"github.com/gke-hackathon/payment-integration/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PaymentServer implements the PaymentService gRPC server
type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	accountMapper      *mapper.AccountMapper
	authenticator      *auth.ServiceAuthenticator
	bankClient         *bank.Client
	transactionCounter int64
	logger             *logging.Logger
}

// NewPaymentServer creates a new instance of PaymentServer
func NewPaymentServer() *PaymentServer {
	// Initialize logger
	logger := logging.NewLogger("payment-integration")

	// Initialize rate limiter
	rateLimit := 10 // Default: 10 transactions per minute per account
	if rateLimitStr := os.Getenv("RATE_LIMIT_PER_MINUTE"); rateLimitStr != "" {
		if limit, err := strconv.Atoi(rateLimitStr); err == nil && limit > 0 {
			rateLimit = limit
		}
	}
	middleware.InitRateLimiter(rateLimit)
	logger.Info("Rate limiter initialized", map[string]interface{}{"limit_per_minute": rateLimit})

	// Get merchant account from environment or use defaults
	merchantAccount := os.Getenv("MERCHANT_ACCOUNT")
	if merchantAccount == "" {
		merchantAccount = "1111111111"
	}

	routingNumber := os.Getenv("ROUTING_NUMBER")
	if routingNumber == "" {
		routingNumber = "123456789"
	}

	// Initialize service authenticator
	privateKeyPath := os.Getenv("PRIV_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "/tmp/.ssh/privatekey"
	}

	publicKeyPath := os.Getenv("PUB_KEY_PATH")
	if publicKeyPath == "" {
		publicKeyPath = "/tmp/.ssh/publickey"
	}

	tokenExpiryStr := os.Getenv("TOKEN_EXPIRY_SECONDS")
	if tokenExpiryStr == "" {
		tokenExpiryStr = "3600"
	}
	tokenExpiry, _ := strconv.Atoi(tokenExpiryStr)

	var authenticator *auth.ServiceAuthenticator
	var authErr error

	// Try to initialize authenticator (non-fatal if it fails during local dev)
	authenticator, authErr = auth.NewServiceAuthenticator(privateKeyPath, publicKeyPath, tokenExpiry)
	if authErr != nil {
		logger.Warn("Failed to initialize service authenticator", map[string]interface{}{
			"error": authErr.Error(),
			"note":  "Service will run without JWT authentication capability",
		})
	} else {
		logger.Info("Service authenticator initialized successfully", nil)
		// Test token generation
		if testToken, err := authenticator.GenerateServiceToken("TEST_ACCOUNT"); err != nil {
			logger.Warn("Failed to generate test token", map[string]interface{}{"error": err.Error()})
		} else {
			logger.Info("Successfully generated test service token", map[string]interface{}{"token_length": len(testToken)})
		}
	}

	// Initialize Bank client
	bankAPIURL := os.Getenv("BANK_API_URL")
	if bankAPIURL == "" {
		bankAPIURL = "http://ledgerwriter.bank-of-anthos.svc.cluster.local:8080"
	}

	var bankClient *bank.Client
	if authenticator != nil {
		bankClient = bank.NewClient(bankAPIURL, authenticator)
		logger.Info("Bank client initialized", map[string]interface{}{"bank_api_url": bankAPIURL})

		// Test bank connectivity
		if err := bankClient.HealthCheck(); err != nil {
			logger.Warn("Bank API health check failed", map[string]interface{}{
				"error": err.Error(),
				"note":  "Transactions may fail until Bank API is available",
			})
		} else {
			logger.Info("Bank API health check successful", nil)
		}
	} else {
		logger.Warn("Bank client not initialized due to missing authenticator", nil)
	}

	return &PaymentServer{
		accountMapper:      mapper.NewAccountMapper(merchantAccount, routingNumber),
		authenticator:      authenticator,
		bankClient:         bankClient,
		transactionCounter: 0,
		logger:             logger,
	}
}

// RegisterPaymentServiceServer registers the payment service with the gRPC server
func RegisterPaymentServiceServer(s *grpc.Server, srv *PaymentServer) {
	pb.RegisterPaymentServiceServer(s, srv)
}

// Charge processes a payment request
func (s *PaymentServer) Charge(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	start := time.Now()

	// Validate request
	if req.Amount == nil {
		return nil, status.Error(codes.InvalidArgument, "amount is required")
	}

	if req.CreditCard == nil {
		return nil, status.Error(codes.InvalidArgument, "credit card info is required")
	}

	// Validate card number
	if err := mapper.ValidateCardNumber(req.CreditCard.CreditCardNumber); err != nil {
		s.logger.Warn("Invalid card number", map[string]interface{}{"error": err.Error()})
		return nil, status.Errorf(codes.InvalidArgument, "invalid card number: %v", err)
	}

	// Convert money format to cents for Bank of Anthos
	cents, err := converter.BoutiqueMoneyToCents(req.Amount)
	if err != nil {
		s.logger.Error("Error converting money", err, nil)
		return nil, status.Errorf(codes.InvalidArgument, "invalid amount: %v", err)
	}

	// Map credit card to bank account
	fromAccount, fromRouting := s.accountMapper.CardNumberToAccount(req.CreditCard.CreditCardNumber)
	toAccount, toRouting := s.accountMapper.GetMerchantAccount()

	// Check rate limit for this account
	rateLimiter := middleware.GetRateLimiter()
	if !rateLimiter.Allow(fromAccount) {
		s.logger.Warn("Rate limit exceeded", map[string]interface{}{
			"account":   fromAccount,
			"remaining": rateLimiter.GetRemaining(fromAccount),
		})
		metrics.GetInstance().RecordRejection()
		return nil, status.Error(codes.ResourceExhausted, "too many payment requests, please try again later")
	}

	// Generate unique transaction UUID for Bank API
	transactionUUID := utils.GenerateUUID()

	// Log the payment request
	s.logger.LogPaymentRequest(ctx, transactionUUID, cents,
		req.Amount.CurrencyCode, getLastFourDigits(req.CreditCard.CreditCardNumber))

	s.logger.Debug("Payment details", map[string]interface{}{
		"transaction_id": transactionUUID,
		"amount_cents":   cents,
		"from_account":   fromAccount,
		"from_routing":   fromRouting,
		"to_account":     toAccount,
		"to_routing":     toRouting,
	})

	// Call the Bank of Anthos API to process the real transaction
	if s.bankClient != nil {
		bankReq := &bank.TransactionRequest{
			FromAccountNum: fromAccount,
			FromRoutingNum: fromRouting,
			ToAccountNum:   toAccount,
			ToRoutingNum:   toRouting,
			Amount:         cents,
			UUID:           transactionUUID,
		}

		bankStart := time.Now()
		_, err := s.bankClient.CreateTransaction(bankReq)
		bankDuration := time.Since(bankStart)

		s.logger.LogBankAPICall(transactionUUID, fromAccount, cents, bankDuration, err)

		if err != nil {
			s.logger.LogPaymentResponse(ctx, transactionUUID, false, time.Since(start), err)

			// Record error metrics
			metrics.GetInstance().RecordRequest(false, time.Since(start), 0, "")
			metrics.GetInstance().RecordError("bank_api_error")

			// Handle specific bank errors
			if bankErr, ok := err.(*bank.BankError); ok {
				return nil, bankErr.ToGRPCError()
			}
			return nil, bank.HandleBankError(err)
		}

		s.logger.LogTransaction(transactionUUID, fromAccount, toAccount, cents,
			req.Amount.CurrencyCode, "Bank transaction successful")
		s.logger.LogPaymentResponse(ctx, transactionUUID, true, time.Since(start), nil)

		// Record metrics
		metrics.GetInstance().RecordRequest(true, time.Since(start), cents, getLastFourDigits(req.CreditCard.CreditCardNumber))

		// Use the transaction UUID as the response ID
		response := &pb.ChargeResponse{
			TransactionId: transactionUUID,
		}
		return response, nil
	} else {
		// Fallback to simulation if Bank client is not available
		s.logger.Warn("Bank client not available, simulating payment", nil)
		s.transactionCounter++
		transactionID := fmt.Sprintf("SIM-%d-%d", time.Now().Unix(), s.transactionCounter)

		s.logger.LogPaymentResponse(ctx, transactionID, true, time.Since(start), nil)

		response := &pb.ChargeResponse{
			TransactionId: transactionID,
		}
		return response, nil
	}
}

// getLastFourDigits returns the last 4 digits of a card number for logging
func getLastFourDigits(cardNumber string) string {
	// Remove spaces and dashes
	cleaned := ""
	for _, ch := range cardNumber {
		if ch >= '0' && ch <= '9' {
			cleaned += string(ch)
		}
	}

	if len(cleaned) >= 4 {
		return cleaned[len(cleaned)-4:]
	}
	return "****"
}
