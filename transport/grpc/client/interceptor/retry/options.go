package retry

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"slices"
	"time"
)

var DefaultRetriableCodes = []codes.Code{codes.ResourceExhausted, codes.Unavailable, codes.Unknown}

type Option func(*options)

// BackoffFunc denotes a family of functions that control the backoff duration between call retries.
//
// They are called with an identifier of the attempt, and should return a time the system client should
// hold off for. If the time returned is longer than the `context.Context.Deadline` of the request
// the deadline of the request takes precedence and the wait will be interrupted before proceeding
// with the next iteration. The context can be used to extract request scoped metadata and context values.
type BackoffFunc func(ctx context.Context, attempt uint) time.Duration

// OnRetryCallback is the type of function called when a retry occurs.
type OnRetryCallback func(ctx context.Context, attempt uint, err error)

// RetriableFunc denotes a family of functions that control which error should be retried.
type RetriableFunc func(err error) bool

// WithMax sets the maximum number of retries on this call, or this interceptor.
func WithMax(maxRetries uint) Option {
	return func(o *options) {
		o.max = maxRetries
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

// WithBackoff sets the `BackoffFunc` used to control time between retries.
func WithBackoff(bf BackoffFunc) Option {
	return func(o *options) {
		o.backoffFunc = bf
	}
}

// WithOnRetryCallback sets the callback to use when a retry occurs.
//
// By default, when no callback function provided, we will just print a log to trace
func WithOnRetryCallback(fn OnRetryCallback) Option {
	return func(o *options) {
		o.onRetryCallback = fn
	}
}

// WithRetriable sets which error should be retried.
func WithRetriable(retriableFunc RetriableFunc) Option {
	return func(o *options) {
		o.retriableFunc = retriableFunc
	}
}

type options struct {
	max             uint
	perCallTimeout  time.Duration
	includeHeader   bool
	backoffFunc     BackoffFunc
	onRetryCallback OnRetryCallback
	retriableFunc   RetriableFunc
}

// CallOption is a grpc.CallOption that is local to grpc_retry.
type CallOption struct {
	grpc.EmptyCallOption // make sure we implement private after() and before() fields so we don't panic.
	applyFunc            func(opt *options)
}

// newRetriableFuncForCodes returns retriable function for specific Codes.
func newRetriableFuncForCodes(codes []codes.Code) func(err error) bool {
	return func(err error) bool {
		errCode := status.Code(err)
		if isContextError(err) {
			// context errors are not retriable based on user settings.
			return false
		}
		if slices.Contains(codes, errCode) {
			return true
		}
		return false
	}
}
