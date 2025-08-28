package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestRSAKeys(t *testing.T) (privateKeyPEM, publicKeyPEM []byte) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Encode private key to PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key to PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	publicKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM
}

func setupTestKeys(t *testing.T) (privateKeyPath, publicKeyPath string) {
	// Generate test key pair
	privateKey, publicKey := generateTestRSAKeys(t)

	// Create temp files for keys
	privFile, err := os.CreateTemp("", "test-private-*.key")
	if err != nil {
		t.Fatal(err)
	}

	pubFile, err := os.CreateTemp("", "test-public-*.key")
	if err != nil {
		t.Fatal(err)
	}

	// Write keys to files
	if _, err := privFile.Write(privateKey); err != nil {
		t.Fatal(err)
	}
	if _, err := pubFile.Write(publicKey); err != nil {
		t.Fatal(err)
	}

	privFile.Close()
	pubFile.Close()

	return privFile.Name(), pubFile.Name()
}

func TestNewServiceAuthenticator(t *testing.T) {
	privateKeyPath, publicKeyPath := setupTestKeys(t)
	defer os.Remove(privateKeyPath)
	defer os.Remove(publicKeyPath)

	authenticator, err := NewServiceAuthenticator(privateKeyPath, publicKeyPath, 3600)
	if err != nil {
		t.Fatalf("Failed to create authenticator: %v", err)
	}

	if authenticator.privateKey == nil {
		t.Error("Private key not loaded")
	}
	if authenticator.publicKey == nil {
		t.Error("Public key not loaded")
	}
	if authenticator.expiryTime != 3600*time.Second {
		t.Errorf("Expiry time mismatch: got %v, want %v", authenticator.expiryTime, 3600*time.Second)
	}
}

func TestGenerateServiceToken(t *testing.T) {
	privateKeyPath, publicKeyPath := setupTestKeys(t)
	defer os.Remove(privateKeyPath)
	defer os.Remove(publicKeyPath)

	authenticator, err := NewServiceAuthenticator(privateKeyPath, publicKeyPath, 3600)
	if err != nil {
		t.Fatalf("Failed to create authenticator: %v", err)
	}

	token, err := authenticator.GenerateServiceToken("1234567890")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token is empty")
	}

	// Verify token structure (should be 3 parts separated by dots)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Invalid JWT structure: expected 3 parts, got %d", len(parts))
	}
}

func TestValidateToken(t *testing.T) {
	privateKeyPath, publicKeyPath := setupTestKeys(t)
	defer os.Remove(privateKeyPath)
	defer os.Remove(publicKeyPath)

	authenticator, err := NewServiceAuthenticator(privateKeyPath, publicKeyPath, 3600)
	if err != nil {
		t.Fatalf("Failed to create authenticator: %v", err)
	}

	// Generate a token
	token, err := authenticator.GenerateServiceToken("1234567890")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := authenticator.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check claims
	if claims.User != "payment-service" {
		t.Errorf("Invalid user claim: got %s, want payment-service", claims.User)
	}
	if claims.Acct != "1234567890" {
		t.Errorf("Invalid account claim: got %s, want 1234567890", claims.Acct)
	}
	if claims.Name != "Payment Integration Service" {
		t.Errorf("Invalid name claim: got %s, want Payment Integration Service", claims.Name)
	}
}

func TestGetAuthHeader(t *testing.T) {
	privateKeyPath, publicKeyPath := setupTestKeys(t)
	defer os.Remove(privateKeyPath)
	defer os.Remove(publicKeyPath)

	authenticator, err := NewServiceAuthenticator(privateKeyPath, publicKeyPath, 3600)
	if err != nil {
		t.Fatalf("Failed to create authenticator: %v", err)
	}

	authHeader, err := authenticator.GetAuthHeader("1234567890")
	if err != nil {
		t.Fatalf("Failed to get auth header: %v", err)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		t.Errorf("Auth header should start with 'Bearer ', got: %s", authHeader[:10])
	}
}

func TestTokenExpiry(t *testing.T) {
	privateKeyPath, publicKeyPath := setupTestKeys(t)
	defer os.Remove(privateKeyPath)
	defer os.Remove(publicKeyPath)

	// Create authenticator with 1 second expiry
	authenticator, err := NewServiceAuthenticator(privateKeyPath, publicKeyPath, 1)
	if err != nil {
		t.Fatalf("Failed to create authenticator: %v", err)
	}

	token, err := authenticator.GenerateServiceToken("1234567890")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Token should be valid immediately
	_, err = authenticator.ValidateToken(token)
	if err != nil {
		t.Errorf("Token should be valid immediately after generation: %v", err)
	}

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Token should now be expired
	_, err = authenticator.ValidateToken(token)
	if err == nil {
		t.Error("Token should be expired after 2 seconds")
	}
	if !strings.Contains(err.Error(), "token is expired") && !strings.Contains(err.Error(), "token has invalid claims") {
		t.Errorf("Expected token expired error, got: %v", err)
	}
}

func TestInvalidToken(t *testing.T) {
	privateKeyPath, publicKeyPath := setupTestKeys(t)
	defer os.Remove(privateKeyPath)
	defer os.Remove(publicKeyPath)

	authenticator, err := NewServiceAuthenticator(privateKeyPath, publicKeyPath, 3600)
	if err != nil {
		t.Fatalf("Failed to create authenticator: %v", err)
	}

	// Test with invalid token
	_, err = authenticator.ValidateToken("invalid.token.string")
	if err == nil {
		t.Error("Should fail with invalid token")
	}

	// Test with empty token
	_, err = authenticator.ValidateToken("")
	if err == nil {
		t.Error("Should fail with empty token")
	}
}

func TestServiceClaimsFormat(t *testing.T) {
	privateKeyPath, publicKeyPath := setupTestKeys(t)
	defer os.Remove(privateKeyPath)
	defer os.Remove(publicKeyPath)

	authenticator, err := NewServiceAuthenticator(privateKeyPath, publicKeyPath, 3600)
	if err != nil {
		t.Fatalf("Failed to create authenticator: %v", err)
	}

	tokenString, err := authenticator.GenerateServiceToken("TEST_ACCOUNT")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Parse token to check claims structure
	token, err := jwt.ParseWithClaims(tokenString, &ServiceClaims{}, func(token *jwt.Token) (interface{}, error) {
		return authenticator.publicKey, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*ServiceClaims)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	// Check all required fields are present
	if claims.User == "" {
		t.Error("User claim is empty")
	}
	if claims.Acct == "" {
		t.Error("Account claim is empty")
	}
	if claims.Name == "" {
		t.Error("Name claim is empty")
	}
	if claims.IssuedAt == nil {
		t.Error("IssuedAt claim is missing")
	}
	if claims.ExpiresAt == nil {
		t.Error("ExpiresAt claim is missing")
	}
}

func TestNonExistentKeyFiles(t *testing.T) {
	// Test with non-existent files
	_, err := NewServiceAuthenticator("/non/existent/private.key", "/non/existent/public.key", 3600)
	if err == nil {
		t.Error("Should fail with non-existent key files")
	}
}

func TestInvalidKeyFormat(t *testing.T) {
	// Create temp files with invalid content
	privFile, _ := os.CreateTemp("", "invalid-private-*.key")
	pubFile, _ := os.CreateTemp("", "invalid-public-*.key")
	defer os.Remove(privFile.Name())
	defer os.Remove(pubFile.Name())

	privFile.Write([]byte("invalid key content"))
	pubFile.Write([]byte("invalid key content"))
	privFile.Close()
	pubFile.Close()

	_, err := NewServiceAuthenticator(privFile.Name(), pubFile.Name(), 3600)
	if err == nil {
		t.Error("Should fail with invalid key format")
	}
}
