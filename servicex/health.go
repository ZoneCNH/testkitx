package servicex

import (
	"context"
	"errors"
)

var ErrNilReady = errors.New("servicex: ready function is required")

type HealthChecker interface{ Healthy(context.Context) error }

func CheckHealth(ctx context.Context, checker HealthChecker) error {
	if checker == nil {
		return errors.New("servicex: health checker is required")
	}
	return checker.Healthy(ctx)
}
