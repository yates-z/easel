package retry

import (
	"google.golang.org/grpc"
	"time"
)

type Option func(*options)

// WithMax sets the maximum number of retries on this call, or this interceptor.
func WithMax(maxRetries uint) Option {
	return func(o *options) {
		o.max = maxRetries
	}
}

// WithInterval sets the interval per call.
func WithInterval(interval time.Duration) Option {
	return func(o *options) {
		o.interval = interval
	}
}

// WithPerRetryTimeout sets the RPC timeout per call (including initial call) on this call, or this interceptor.
//
// The context.Deadline of the call takes precedence and sets the maximum time the whole invocation
// will take, but WithPerRetryTimeout can be used to limit the RPC time per each call.
//
// For example, with context.Deadline = now + 10s, and WithPerRetryTimeout(3 * time.Seconds), each
// of the retry calls (including the initial one) will have a deadline of now + 3s.
//
// A value of 0 disables the timeout overrides completely and returns to each retry call using the
// parent `context.Deadline`.
//
// Note that when this is enabled, any DeadlineExceeded errors that are propagated up will be retried.
func WithPerRetryTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.perCallTimeout = timeout
	}
}

// WithHeader denotes if set the parent header to new call.
func WithHeader(includeHeader bool) Option {
	return func(o *options) {
		o.includeHeader = includeHeader
	}
}

type options struct {
	max            uint
	interval       time.Duration
	skipPanic      bool
	perCallTimeout time.Duration
	includeHeader  bool
}

// CallOption is a grpc.CallOption that is local to grpc_retry.
type CallOption struct {
	grpc.EmptyCallOption // make sure we implement private after() and before() fields so we don't panic.
	applyFunc            func(opt *options)
}
