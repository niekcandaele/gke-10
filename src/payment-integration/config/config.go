package config

import (
	"log"
	"os"
)

// Config holds the configuration for the payment integration service
type Config struct {
	Port            string
	MerchantAccount string
	RoutingNumber   string
	BankAPIURL      string
	LogLevel        string
	PrivateKeyPath  string
	PublicKeyPath   string
	TokenExpiry     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		Port:            getEnv("PORT", "50051"),
		MerchantAccount: getEnv("MERCHANT_ACCOUNT", "1111111111"),
		RoutingNumber:   getEnv("ROUTING_NUMBER", "123456789"),
		BankAPIURL:      getEnv("BANK_API_URL", "http://ledgerwriter.bank-of-anthos.svc.cluster.local:8080"),
		LogLevel:        getEnv("LOG_LEVEL", "INFO"),
		PrivateKeyPath:  getEnv("PRIV_KEY_PATH", "/tmp/.ssh/privatekey"),
		PublicKeyPath:   getEnv("PUB_KEY_PATH", "/tmp/.ssh/publickey"),
		TokenExpiry:     getEnv("TOKEN_EXPIRY_SECONDS", "3600"),
	}

	log.Printf("Configuration loaded:")
	log.Printf("  Port: %s", cfg.Port)
	log.Printf("  Merchant Account: %s", cfg.MerchantAccount)
	log.Printf("  Routing Number: %s", cfg.RoutingNumber)
	log.Printf("  Bank API URL: %s", cfg.BankAPIURL)
	log.Printf("  Log Level: %s", cfg.LogLevel)
	log.Printf("  JWT Private Key Path: %s", cfg.PrivateKeyPath)
	log.Printf("  JWT Public Key Path: %s", cfg.PublicKeyPath)
	log.Printf("  Token Expiry Seconds: %s", cfg.TokenExpiry)

	return cfg
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
