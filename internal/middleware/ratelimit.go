package middleware

import (
	"sync"
	"time"

	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	"github.com/gumla-hds/gumla-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	tokens    int
	lastSeen  time.Time
}

type RateLimiter struct {
	mu       sync.RWMutex
	visitors map[string]*visitor
	rate     int
	interval time.Duration
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     requestsPerMinute,
		interval: time.Minute,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &visitor{tokens: rl.rate - 1, lastSeen: now}
		return true
	}

	elapsed := now.Sub(v.lastSeen)
	v.lastSeen = now

	refill := int(elapsed / rl.interval * time.Duration(rl.rate))
	if refill > 0 {
		v.tokens = min(v.tokens+refill, rl.rate)
	}

	if v.tokens <= 0 {
		return false
	}

	v.tokens--
	return true
}

func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip) {
			response.Error(c, apperrors.New(
				apperrors.CodeBadRequest,
				"rate limit exceeded",
				"बहुत अधिक अनुरोध, कृपया कुछ समय बाद पुनः प्रयास करें",
			))
			c.Abort()
			return
		}
		c.Next()
	}
}
