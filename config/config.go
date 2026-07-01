package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	GinMode       string
	DatabaseURL   string
	JWTSecret     string
	JWTConfig     JWTConfig
	Firebase      FirebaseConfig
	Razorpay      RazorpayConfig
	RateLimit     RateLimitConfig
}

type JWTConfig struct {
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type FirebaseConfig struct {
	CredentialsJSON string
	StorageBucket   string
}

type RazorpayConfig struct {
	KeyID           string
	KeySecret       string
	WebhookSecret   string
}

type RateLimitConfig struct {
	RequestsPerMinute int
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gumla_hds?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production-must-be-64-chars-long-hex-string"),
		JWTConfig: JWTConfig{
			AccessExpiry:  getDurationEnv("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},
		Firebase: FirebaseConfig{
			CredentialsJSON: getEnv("FIREBASE_CREDENTIALS_JSON", ""),
			StorageBucket:   getEnv("FIREBASE_STORAGE_BUCKET", ""),
		},
		Razorpay: RazorpayConfig{
			KeyID:         getEnv("RAZORPAY_KEY_ID", ""),
			KeySecret:     getEnv("RAZORPAY_KEY_SECRET", ""),
			WebhookSecret: getEnv("RAZORPAY_WEBHOOK_SECRET", ""),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getIntEnv("RATE_LIMIT_PER_MINUTE", 100),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}
	return fallback
}
