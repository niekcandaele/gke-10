package converter

import (
	"fmt"
	pb "github.com/gke-hackathon/payment-integration/proto"
)

// BoutiqueMoneyToCents converts Online Boutique Money format to cents (used by Bank of Anthos)
// Boutique Money format: units (dollars) + nanos (billionths of a dollar)
// Bank format: cents (hundredths of a dollar)
func BoutiqueMoneyToCents(money *pb.Money) (int64, error) {
	if money == nil {
		return 0, fmt.Errorf("money cannot be nil")
	}

	// Convert units to cents (1 unit = 100 cents)
	cents := money.Units * 100

	// Convert nanos to cents (1 nano = 10^-9 units, 1 cent = 10^-2 units)
	// So 1 nano = 10^-7 cents, or 10^7 nanos = 1 cent
	nanosInCents := money.Nanos / 10_000_000

	totalCents := cents + int64(nanosInCents)

	// Validate the result is reasonable (prevent overflow/underflow)
	if totalCents < 0 && money.Units >= 0 && money.Nanos >= 0 {
		return 0, fmt.Errorf("money conversion overflow")
	}

	return totalCents, nil
}

// CentsToBoutiqueMoney converts cents (Bank of Anthos) to Online Boutique Money format
func CentsToBoutiqueMoney(cents int64, currencyCode string) *pb.Money {
	// Calculate whole units (dollars)
	units := cents / 100

	// Calculate remaining cents and convert to nanos
	remainingCents := cents % 100
	nanos := remainingCents * 10_000_000

	// Handle negative amounts
	if cents < 0 {
		// For negative amounts, both units and nanos must be negative or zero
		if units == 0 && nanos > 0 {
			nanos = -nanos
		}
	}

	return &pb.Money{
		Units:        units,
		Nanos:        int32(nanos),
		CurrencyCode: currencyCode,
	}
}

// FormatMoney returns a human-readable string representation of Money
func FormatMoney(money *pb.Money) string {
	if money == nil {
		return "$0.00"
	}

	cents, err := BoutiqueMoneyToCents(money)
	if err != nil {
		return "invalid"
	}

	// Format as dollars with 2 decimal places
	dollars := float64(cents) / 100.0
	return fmt.Sprintf("%s %.2f", money.CurrencyCode, dollars)
}
