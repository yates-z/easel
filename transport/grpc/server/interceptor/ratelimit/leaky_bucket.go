package ratelimit

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// Note: This file is inspired by:
// https://github.com/uber-go/ratelimit

var leakDuration = time.Second

type LeakyBucket struct {
	// prePadding is created to avoid false sharing.
	prePadding   [64]byte
	lastTimeNano int64
	// postPadding is created to avoid false sharing.cache line size = 64 - 8.
	postPadding [56]byte
	// maxWaitingNano is the max time for waiting.
	maxWaitingNano int64
	perRequest     time.Duration
}

func NewLeakyBucket(rate, capacity int64) Limiter {
	if rate <= 0 || capacity <= 0 {
		panic(errors.New("invalid rate or capacity"))
	}
	perRequest := leakDuration / time.Duration(rate)
	bucket := &LeakyBucket{
		perRequest:     perRequest,
		maxWaitingNano: capacity * int64(perRequest),
	}
	atomic.StoreInt64(&bucket.lastTimeNano, 0)
	return bucket
}

func (b *LeakyBucket) Limit(ctx context.Context) error {
	var newLastTimeNano int64
	var now int64
	for {
		now = time.Now().UnixNano()
		lastTimeNano := atomic.LoadInt64(&b.lastTimeNano)
		if now-lastTimeNano > int64(b.perRequest) {
			newLastTimeNano = now
		} else {
			newLastTimeNano = lastTimeNano + int64(b.perRequest)
		}
		//fmt.Println(b.maxWaitingNano, newLastTimeNano-now)
		if newLastTimeNano-now >= b.maxWaitingNano {
			return errors.New("")
		}
		// Swapping may be failed, try again.
		if atomic.CompareAndSwapInt64(&b.lastTimeNano, lastTimeNano, newLastTimeNano) {
			break
		}
	}
	sleepDuration := time.Duration(newLastTimeNano - now)

	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
	return nil
}
