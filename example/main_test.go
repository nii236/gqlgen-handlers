package main

import (
	"bytes"
	"context"
	"handlers"
	"testing"
)

func TestOnboardStart(t *testing.T) {
	fn := OnboardStart(nil, nil)

	type args struct {
		ctx context.Context
		w   handlers.Writer
		r   handlers.Reader
	}
	tests := []struct {
		name        string
		args        args
		wantUserErr bool
		wantErr     bool
	}{
		{"happy path", args{context.Background(), &bytes.Buffer{}, handlers.MustNewReader(&OnboardStartRequest{})}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fn(tt.args.ctx, tt.args.w, tt.args.r)
			if !tt.wantErr && err != nil {
				t.Logf("Unexpected err, expected nil, got %v", err)
			}
		})
	}
}
