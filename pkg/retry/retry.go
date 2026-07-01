package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

type Config struct {
	MaxAttempts   int
	BaseDelay     time.Duration
	MaxDelay      time.Duration
}

func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    2 * time.Second,
	}
}

func Do(ctx context.Context, config Config, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if attempt > 0 {
			delay := time.Duration(float64(config.BaseDelay) * math.Pow(2, float64(attempt-1)))
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		if err := fn(); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return fmt.Errorf("all %d attempts failed: %w", config.MaxAttempts, lastErr)
}

func DoWithResult[T any](ctx context.Context, config Config, fn func() (T, error)) (T, error) {
	var zero T
	var lastErr error

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if attempt > 0 {
			delay := time.Duration(float64(config.BaseDelay) * math.Pow(2, float64(attempt-1)))
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}

			select {
			case <-ctx.Done():
				return zero, ctx.Err()
			case <-time.After(delay):
			}
		}

		result, err := fn()
		if err != nil {
			lastErr = err
			continue
		}

		return result, nil
	}

	return zero, fmt.Errorf("all %d attempts failed: %w", config.MaxAttempts, lastErr)
}
