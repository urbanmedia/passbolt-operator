package util

import (
	"context"
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	idxCount := 0

	type args struct {
		ctx         context.Context
		maxRetries  int
		intervalSec int
		retryFunc   func(ctx context.Context) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "retry 0 times without error",
			args: args{
				ctx:         context.Background(),
				maxRetries:  0,
				intervalSec: 0,
				retryFunc: func(ctx context.Context) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "retry 0 times with error",
			args: args{
				ctx:         context.Background(),
				maxRetries:  0,
				intervalSec: 0,
				retryFunc: func(ctx context.Context) error {
					return errors.New("some error")
				},
			},
			wantErr: true,
		},

		{
			name: "retry 1 times without error",
			args: args{
				ctx:         context.Background(),
				maxRetries:  1,
				intervalSec: 0,
				retryFunc: func(ctx context.Context) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "retry 1 times with error",
			args: args{
				ctx:         context.Background(),
				maxRetries:  1,
				intervalSec: 0,
				retryFunc: func(ctx context.Context) error {
					return errors.New("some error")
				},
			},
			wantErr: true,
		},

		{
			name: "first time error, second time no error",
			args: args{
				ctx:         context.Background(),
				maxRetries:  1,
				intervalSec: 0,
				retryFunc: func(ctx context.Context) error {
					// throw error at first time and return nil at second time
					if idxCount == 0 {
						idxCount++
						return errors.New("some error")
					}
					idxCount = 0
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Retry(tt.args.ctx, tt.args.maxRetries, tt.args.intervalSec, tt.args.retryFunc); (err != nil) != tt.wantErr {
				t.Errorf("Retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
