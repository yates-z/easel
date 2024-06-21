package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestLeakyBucket(t *testing.T) {
	wg := sync.WaitGroup{}
	limiter := NewLeakyBucket(2, 2)
	now := time.Now()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e := limiter.Limit(context.Background())
			if e != nil {
				t.Log(i, "---", e, time.Since(now))
				return
			}
			t.Log(i, "---", "success", time.Since(now))
		}()
	}
	wg.Wait()
}
