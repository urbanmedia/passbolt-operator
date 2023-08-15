package util

import (
	"context"
	"time"
)

// Retry will retry the given function until it returns nil or the maxRetries is reached.
// If the maxRetries is reached, the last error will be returned.
// 0 retries means the function will be called once without retrying.
func Retry(ctx context.Context, maxRetries int, intervalSec int, retryFunc func(ctx context.Context) error) error {
	var err error
	for i := 0; i < 1+maxRetries; i++ {
		err = retryFunc(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(intervalSec) * time.Second)
	}
	return err
}
