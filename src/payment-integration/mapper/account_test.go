package mapper

import (
	"testing"
)

func TestCardNumberToAccount(t *testing.T) {
	mapper := NewAccountMapper("MERCHANT-001", "987654321")

	tests := []struct {
		name          string
		cardNumber    string
		expectedAcct  string
		expectedRoute string
	}{
		{
			name:          "Valid 16-digit card",
			cardNumber:    "4532015112830366",
			expectedAcct:  "5112830366",
			expectedRoute: "987654321",
		},
		{
			name:          "Card with spaces",
			cardNumber:    "4532 0151 1283 0366",
			expectedAcct:  "5112830366",
			expectedRoute: "987654321",
		},
		{
			name:          "Card with dashes",
			cardNumber:    "4532-0151-1283-0366",
			expectedAcct:  "5112830366",
			expectedRoute: "987654321",
		},
		{
			name:          "15-digit card (AmEx format)",
			cardNumber:    "378282246310005",
			expectedAcct:  "2246310005",
			expectedRoute: "987654321",
		},
		{
			name:          "19-digit card",
			cardNumber:    "6011111111111117890",
			expectedAcct:  "1111117890",
			expectedRoute: "987654321",
		},
		{
			name:          "Short card number (uses default)",
			cardNumber:    "123456789",
			expectedAcct:  "MERCHANT-001",
			expectedRoute: "987654321",
		},
		{
			name:          "Empty card number",
			cardNumber:    "",
			expectedAcct:  "MERCHANT-001",
			expectedRoute: "987654321",
		},
		{
			name:          "Card ending in zeros",
			cardNumber:    "4000000000000000",
			expectedAcct:  "0000000000",
			expectedRoute: "987654321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acct, route := mapper.CardNumberToAccount(tt.cardNumber)

			if acct != tt.expectedAcct {
				t.Errorf("Expected account %s, got %s", tt.expectedAcct, acct)
			}
			if route != tt.expectedRoute {
				t.Errorf("Expected routing %s, got %s", tt.expectedRoute, route)
			}
		})
	}
}

func TestNewAccountMapper(t *testing.T) {
	tests := []struct {
		name          string
		merchantAcct  string
		routingNum    string
		expectedMerch string
		expectedRoute string
	}{
		{
			name:          "With provided values",
			merchantAcct:  "CUSTOM-123",
			routingNum:    "555666777",
			expectedMerch: "CUSTOM-123",
			expectedRoute: "555666777",
		},
		{
			name:          "With empty values (uses defaults)",
			merchantAcct:  "",
			routingNum:    "",
			expectedMerch: "1111111111",
			expectedRoute: "123456789",
		},
		{
			name:          "Mixed empty and provided",
			merchantAcct:  "ACCT-999",
			routingNum:    "",
			expectedMerch: "ACCT-999",
			expectedRoute: "123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := NewAccountMapper(tt.merchantAcct, tt.routingNum)

			if mapper.DefaultMerchantAccount != tt.expectedMerch {
				t.Errorf("Expected merchant account %s, got %s",
					tt.expectedMerch, mapper.DefaultMerchantAccount)
			}
			if mapper.DefaultRoutingNumber != tt.expectedRoute {
				t.Errorf("Expected routing number %s, got %s",
					tt.expectedRoute, mapper.DefaultRoutingNumber)
			}
		})
	}
}

func TestGetMerchantAccount(t *testing.T) {
	mapper := NewAccountMapper("MERCHANT-XYZ", "111222333")

	acct, route := mapper.GetMerchantAccount()

	if acct != "MERCHANT-XYZ" {
		t.Errorf("Expected merchant account MERCHANT-XYZ, got %s", acct)
	}
	if route != "111222333" {
		t.Errorf("Expected routing number 111222333, got %s", route)
	}
}

func TestValidateCardNumber(t *testing.T) {
	tests := []struct {
		name      string
		cardNum   string
		shouldErr bool
	}{
		{
			name:      "Valid 16-digit",
			cardNum:   "4532015112830366",
			shouldErr: false,
		},
		{
			name:      "Valid 15-digit",
			cardNum:   "378282246310005",
			shouldErr: false,
		},
		{
			name:      "Valid with spaces",
			cardNum:   "4532 0151 1283 0366",
			shouldErr: false,
		},
		{
			name:      "Too short",
			cardNum:   "123456789012",
			shouldErr: true,
		},
		{
			name:      "Too long",
			cardNum:   "12345678901234567890",
			shouldErr: true,
		},
		{
			name:      "Contains letters",
			cardNum:   "4532ABC112830366",
			shouldErr: true,
		},
		{
			name:      "Empty",
			cardNum:   "",
			shouldErr: true,
		},
		{
			name:      "Special characters",
			cardNum:   "4532-0151-1283-036!",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCardNumber(tt.cardNum)

			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for card %s, but got none", tt.cardNum)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for card %s: %v", tt.cardNum, err)
			}
		})
	}
}
