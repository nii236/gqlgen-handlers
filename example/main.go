package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"handlers"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	ctx := context.Background()
	duration, saturation, calls, errors := handlers.MustRegisterMetrics("example", "onboarding")
	r := &RootResolver{
		userStore:    nil,
		accountStore: nil,

		Duration:   duration,
		Calls:      calls,
		Saturation: saturation,
		Errors:     errors,
	}
	r.RunExampleResolver(ctx)
}

// UserStore is an example of a persistence layer
type UserStore interface{}

// AccountStore is an example of a persistence layer
type AccountStore interface{}

// RootResolver will hold dependencies
type RootResolver struct {
	userStore    UserStore
	accountStore AccountStore

	Duration   *prometheus.HistogramVec
	Calls      *prometheus.CounterVec
	Saturation *prometheus.GaugeVec
	Errors     *prometheus.CounterVec
}

// RunExampleResolver shows how you would setup one resolver to use this pattern
// The resolver would typically be a mutation, as queries have their own implementation
// The GraphQL server acts as a thin proxy that maps and passes requests onto the handler
// - Initialise request structs
// - Pass in the writer and reader
// - Add in the middleware
// - Execute with context
func (r *RootResolver) RunExampleResolver(ctx context.Context) {
	// Prepare args
	req := handlers.MustNewReader(&OnboardStartRequest{})
	resp := &bytes.Buffer{}

	// Prepare func
	fn := handlers.WithLogging(handlers.WithMetrics("start", r.Calls, r.Errors, r.Saturation, r.Duration, OnboardStart(r.userStore, r.accountStore)))

	// Execute func and handle errors
	err := fn(ctx, resp, req)
	var userErr *handlers.UserError
	var sysErr *handlers.SystemError
	if errors.As(err, &userErr) {
		fmt.Println("Input Error:", userErr.Message)
		return
	}
	if errors.As(err, &sysErr) {
		fmt.Println("Internal Error")
		return
	}

}

// Controller methods

// OnboardStartRequest is used to pass arguments into the OnboardStart method
type OnboardStartRequest struct {
}

// OnboardStartResponse is used to respond to the caller of OnboardStart
type OnboardStartResponse struct {
}

// OnboardStart holds the business logic in a signature agnostic way, allowing for middlewares to be applied
// It returns a func, allowing for additional data to be passed via closures
func OnboardStart(userStore UserStore, accountStore AccountStore) handlers.Handler {
	fn := func(ctx context.Context, w handlers.Writer, r handlers.Reader) error {
		req := &OnboardStartRequest{}
		handlers.MustDecode(r, req)
		// Handle validations
		if true {
			return &handlers.UserError{"3c157e31-0a63-4f96-9c9c-19353024ce34", "wrong username", errors.New("example error")}
		}

		// Handle logic
		if true {
			return &handlers.SystemError{"8c24c633-1bba-4990-813c-b21d24d6e7f5", "database connection failure", errors.New("example error")}
		}

		// Handle response
		resp := &OnboardStartResponse{}
		handlers.MustEncode(w, resp)

		// Return user errors and system errors
		return nil
	}
	return fn
}
