package servicex

import (
	"context"
	"errors"
	"time"
)

var ErrWaitTimeout = errors.New("servicex: wait timed out")

const DefaultWaitTimeout = 30 * time.Second

func WaitUntil(ctx context.Context, interval time.Duration, ready func(context.Context) (bool, error)) error {
	if ready == nil {
		return ErrNilReady
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultWaitTimeout)
		defer cancel()
	}
	if interval <= 0 {
		interval = time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.Join(ErrWaitTimeout, ctx.Err())
		default:
		}
		ok, err := ready(ctx)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		select {
		case <-ctx.Done():
			return errors.Join(ErrWaitTimeout, ctx.Err())
		case <-ticker.C:
		}
	}
}
