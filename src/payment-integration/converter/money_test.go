package converter

import (
	pb "github.com/gke-hackathon/payment-integration/proto"
	"testing"
)

func TestBoutiqueMoneyToCents(t *testing.T) {
	tests := []struct {
		name     string
		money    *pb.Money
		expected int64
		hasError bool
	}{
		{
			name: "Simple dollar amount",
			money: &pb.Money{
				Units:        10,
				Nanos:        0,
				CurrencyCode: "USD",
			},
			expected: 1000, // $10 = 1000 cents
			hasError: false,
		},
		{
			name: "Dollar with cents",
			money: &pb.Money{
				Units:        15,
				Nanos:        500_000_000, // 0.50
				CurrencyCode: "USD",
			},
			expected: 1550, // $15.50 = 1550 cents
			hasError: false,
		},
		{
			name: "Zero amount",
			money: &pb.Money{
				Units:        0,
				Nanos:        0,
				CurrencyCode: "USD",
			},
			expected: 0,
			hasError: false,
		},
		{
			name: "Only cents",
			money: &pb.Money{
				Units:        0,
				Nanos:        990_000_000, // 0.99
				CurrencyCode: "USD",
			},
			expected: 99, // $0.99 = 99 cents
			hasError: false,
		},
		{
			name: "Negative amount",
			money: &pb.Money{
				Units:        -5,
				Nanos:        -250_000_000, // -0.25
				CurrencyCode: "USD",
			},
			expected: -525, // -$5.25 = -525 cents
			hasError: false,
		},
		{
			name: "Large amount",
			money: &pb.Money{
				Units:        999999,
				Nanos:        990_000_000,
				CurrencyCode: "USD",
			},
			expected: 99999999, // $999,999.99
			hasError: false,
		},
		{
			name:     "Nil money",
			money:    nil,
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BoutiqueMoneyToCents(tt.money)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %d cents, got %d", tt.expected, result)
				}
			}
		})
	}
}

func TestCentsToBoutiqueMoney(t *testing.T) {
	tests := []struct {
		name          string
		cents         int64
		currencyCode  string
		expectedUnits int64
		expectedNanos int32
	}{
		{
			name:          "Simple dollar amount",
			cents:         1000,
			currencyCode:  "USD",
			expectedUnits: 10,
			expectedNanos: 0,
		},
		{
			name:          "Dollar with cents",
			cents:         1550,
			currencyCode:  "USD",
			expectedUnits: 15,
			expectedNanos: 500_000_000,
		},
		{
			name:          "Zero amount",
			cents:         0,
			currencyCode:  "USD",
			expectedUnits: 0,
			expectedNanos: 0,
		},
		{
			name:          "Only cents",
			cents:         99,
			currencyCode:  "USD",
			expectedUnits: 0,
			expectedNanos: 990_000_000,
		},
		{
			name:          "Negative amount",
			cents:         -525,
			currencyCode:  "USD",
			expectedUnits: -5,
			expectedNanos: -250_000_000,
		},
		{
			name:          "One cent",
			cents:         1,
			currencyCode:  "USD",
			expectedUnits: 0,
			expectedNanos: 10_000_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CentsToBoutiqueMoney(tt.cents, tt.currencyCode)

			if result.Units != tt.expectedUnits {
				t.Errorf("Expected units %d, got %d", tt.expectedUnits, result.Units)
			}
			if result.Nanos != tt.expectedNanos {
				t.Errorf("Expected nanos %d, got %d", tt.expectedNanos, result.Nanos)
			}
			if result.CurrencyCode != tt.currencyCode {
				t.Errorf("Expected currency %s, got %s", tt.currencyCode, result.CurrencyCode)
			}
		})
	}
}

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		name     string
		money    *pb.Money
		expected string
	}{
		{
			name: "Simple dollar amount",
			money: &pb.Money{
				Units:        10,
				Nanos:        0,
				CurrencyCode: "USD",
			},
			expected: "USD 10.00",
		},
		{
			name: "Dollar with cents",
			money: &pb.Money{
				Units:        15,
				Nanos:        990_000_000,
				CurrencyCode: "EUR",
			},
			expected: "EUR 15.99",
		},
		{
			name:     "Nil money",
			money:    nil,
			expected: "$0.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatMoney(tt.money)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
