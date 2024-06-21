package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity int64
	// rate is the count of tokens will be released per fillDuration.
	rate   int64
	tokens int64
	// fillDuration is the time interval for filling the bucket.
	fillDuration time.Duration
	// fillNum is the number of tokens for filling.
	fillNum    int64
	mu         sync.Mutex
	latestTime time.Time
}

func NewTokenBucket(fillDuration time.Duration, rate, capacity int64) Limiter {
	bucket := &TokenBucket{
		capacity:     capacity,
		rate:         rate,
		fillDuration: fillDuration,
		tokens:       capacity,
		latestTime:   time.Now(),
	}
	return bucket
}

func (b *TokenBucket) Limit(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	interval := int64(now.Sub(b.latestTime) / b.fillDuration)
	if interval > 0 {
		b.latestTime = now
	}

	if b.tokens < b.capacity {
		b.tokens += b.rate * interval
	}
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}

	if b.tokens >= 1 {
		b.tokens--
		return nil
	}

	return errors.New("bucket is empty")
}
