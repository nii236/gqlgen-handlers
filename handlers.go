package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Common interfaces and types

// Handler is a custom handler type, not made for HTTP but rather for Go only when using gqlgen
type Handler func(ctx context.Context, w Writer, r Reader) error

// Writer provides an agnostic output from the handler methods
type Writer interface {
	io.Writer
}

// Reader will let request structs come in in a generic form
type Reader interface {
	io.Reader
}

// MustRegisterMetrics for prometheus
func MustRegisterMetrics(namespace, subsystem string) (*prometheus.HistogramVec, *prometheus.GaugeVec, *prometheus.CounterVec, *prometheus.CounterVec) {
	metricDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      "Duration histogram of time taken to execute requests",
		Name:      "request_duration_milliseconds",
	}, []string{"method"})
	metricSaturation := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      "Saturation levels",
		Name:      "in_flight_total",
	}, []string{"method"})
	metricCalls := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      "Number of calls",
		Name:      "calls_total",
	}, []string{"method"})
	metricErrors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      "Number of errors",
		Name:      "errors_total",
	}, []string{"method"})

	prometheus.MustRegister(metricDuration)
	prometheus.MustRegister(metricCalls)
	prometheus.MustRegister(metricSaturation)
	prometheus.MustRegister(metricErrors)
	return metricDuration, metricSaturation, metricCalls, metricErrors

}

// UserError holds input and validation errors
type UserError struct {
	ID      string
	Message string
	Err     error
}

func (e *UserError) Error() string {
	return e.Message
}

// SystemError  holds internal system errors
type SystemError struct {
	ID      string
	Message string
	Err     error
}

func (e *SystemError) Error() string {
	return e.Message
}

// Middlewares

// WithAdmin prevents downstream execution if the user in context is not admin
func WithAdmin(next Handler) Handler {
	fn := func(ctx context.Context, w Writer, r Reader) error {
		// Check context for admin

		isAdmin, ok := ctx.Value("is_admin").(bool)
		if !isAdmin {
			return &UserError{Message: "user not admin"}
		}
		if !ok {
			return &SystemError{Message: "is_admin not in context"}
		}
		return next(ctx, w, r)
	}
	return fn
}

// WithMetrics will execute relevant prometheus values
func WithMetrics(method string,
	calls *prometheus.CounterVec,
	errors *prometheus.CounterVec,
	saturation *prometheus.GaugeVec,
	duration *prometheus.HistogramVec,
	next Handler) Handler {
	fn := func(ctx context.Context, w Writer, r Reader) error {
		// Handle metrics
		start := time.Now()
		calls.WithLabelValues(method).Inc()
		saturation.WithLabelValues(method).Inc()
		err := next(ctx, w, r)
		if err != nil {
			errors.WithLabelValues(method).Inc()
		}
		saturation.WithLabelValues(method).Dec()
		duration.WithLabelValues(method).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
		return err
	}
	return fn
}

// WithLogging will add start, end and duration logs
func WithLogging(next Handler) Handler {
	fn := func(ctx context.Context, w Writer, r Reader) error {
		// Handle logging
		now := time.Now()
		err := next(ctx, w, r)
		since := time.Since(now)
		fmt.Println(since.Nanoseconds())
		return err
	}
	return fn
}

// Helper funcs

// MustEncode will write to w the JSON encoded value
// It will panic on error
func MustEncode(w io.Writer, v interface{}) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
}

// MustDecode read from r and hydrate v
// It will panic on error
func MustDecode(r io.Reader, v interface{}) {
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		panic(err)
	}
}

// MustNewReader return a Reader of v
// It will panic on error
func MustNewReader(v interface{}) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}
