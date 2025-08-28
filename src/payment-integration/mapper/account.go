package mapper

import (
	"fmt"
	"strings"
)

// AccountMapper handles mapping between credit card numbers and bank account numbers
type AccountMapper struct {
	DefaultMerchantAccount string
	DefaultRoutingNumber   string
}

// NewAccountMapper creates a new account mapper instance
func NewAccountMapper(merchantAccount, routingNumber string) *AccountMapper {
	// Use defaults if not provided
	if merchantAccount == "" {
		merchantAccount = "1111111111" // Default merchant account
	}
	if routingNumber == "" {
		routingNumber = "123456789" // Default routing number
	}

	return &AccountMapper{
		DefaultMerchantAccount: merchantAccount,
		DefaultRoutingNumber:   routingNumber,
	}
}

// CardNumberToAccount maps a credit card number to a bank account number
// Initial implementation: use last 10 digits of card number as account number
func (m *AccountMapper) CardNumberToAccount(cardNumber string) (accountNum string, routingNum string) {
	// Remove spaces and dashes from card number
	cleaned := strings.ReplaceAll(cardNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Validate card number length (should be 13-19 digits for valid cards)
	if len(cleaned) < 10 {
		// If card number is too short, use default merchant account
		return m.DefaultMerchantAccount, m.DefaultRoutingNumber
	}

	// Extract last 10 digits as account number
	if len(cleaned) >= 10 {
		accountNum = cleaned[len(cleaned)-10:]
	} else {
		// Pad with zeros if needed (shouldn't happen with valid cards)
		accountNum = fmt.Sprintf("%010s", cleaned)
	}

	// For now, use the default routing number for all accounts
	// In a real system, this would be looked up from a database
	routingNum = m.DefaultRoutingNumber

	return accountNum, routingNum
}

// GetMerchantAccount returns the merchant account details for receiving payments
func (m *AccountMapper) GetMerchantAccount() (accountNum string, routingNum string) {
	return m.DefaultMerchantAccount, m.DefaultRoutingNumber
}

// ValidateCardNumber performs basic validation on a credit card number
// This is a simplified check - real systems would use Luhn algorithm
func ValidateCardNumber(cardNumber string) error {
	// Remove spaces and dashes
	cleaned := strings.ReplaceAll(cardNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Check length (13-19 digits for valid cards)
	if len(cleaned) < 13 || len(cleaned) > 19 {
		return fmt.Errorf("invalid card number length: %d", len(cleaned))
	}

	// Check that all characters are digits
	for _, ch := range cleaned {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("card number contains non-digit characters")
		}
	}

	return nil
}
