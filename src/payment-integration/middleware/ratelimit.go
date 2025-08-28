package middleware

import (
	"sync"
	"time"
)

// RateLimiter implements a simple rate limiter per account
type RateLimiter struct {
	mu            sync.RWMutex
	limits        map[string]*accountLimit
	maxPerMinute  int
	cleanupTicker *time.Ticker
}

type accountLimit struct {
	count     int
	resetTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		limits:        make(map[string]*accountLimit),
		maxPerMinute:  maxPerMinute,
		cleanupTicker: time.NewTicker(5 * time.Minute),
	}

	// Cleanup old entries periodically
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given account is allowed
func (rl *RateLimiter) Allow(accountNumber string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	limit, exists := rl.limits[accountNumber]

	if !exists {
		// First request from this account
		rl.limits[accountNumber] = &accountLimit{
			count:     1,
			resetTime: now.Add(time.Minute),
		}
		return true
	}

	// Check if we need to reset the counter
	if now.After(limit.resetTime) {
		limit.count = 1
		limit.resetTime = now.Add(time.Minute)
		return true
	}

	// Check if under limit
	if limit.count < rl.maxPerMinute {
		limit.count++
		return true
	}

	// Rate limit exceeded
	return false
}

// GetRemaining returns the number of requests remaining for an account
func (rl *RateLimiter) GetRemaining(accountNumber string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	limit, exists := rl.limits[accountNumber]
	if !exists {
		return rl.maxPerMinute
	}

	if time.Now().After(limit.resetTime) {
		return rl.maxPerMinute
	}

	remaining := rl.maxPerMinute - limit.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// cleanup removes old entries to prevent memory leak
func (rl *RateLimiter) cleanup() {
	for range rl.cleanupTicker.C {
		rl.mu.Lock()
		now := time.Now()

		// Remove entries older than 2 minutes
		for account, limit := range rl.limits {
			if now.After(limit.resetTime.Add(2 * time.Minute)) {
				delete(rl.limits, account)
			}
		}
		rl.mu.Unlock()
	}
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	rl.cleanupTicker.Stop()
}

// GlobalRateLimiter is a singleton instance
var globalRateLimiter *RateLimiter

// InitRateLimiter initializes the global rate limiter
func InitRateLimiter(maxPerMinute int) {
	if maxPerMinute <= 0 {
		maxPerMinute = 10 // Default: 10 transactions per minute per account
	}
	globalRateLimiter = NewRateLimiter(maxPerMinute)
}

// GetRateLimiter returns the global rate limiter instance
func GetRateLimiter() *RateLimiter {
	if globalRateLimiter == nil {
		InitRateLimiter(10) // Initialize with default if not set
	}
	return globalRateLimiter
}
