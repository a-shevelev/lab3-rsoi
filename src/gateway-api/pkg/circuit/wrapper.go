package circuit

import (
	"errors"
	"gateway-api/pkg/ext"
)

func WithCircuitBreaker[T any](
	b *Breaker,
	action func() (T, error),
	fallback func() T,
	isHealthy func() bool,
) (T, error) {
	var zero T
	if !isHealthy() {
		b.recordFailure()
		return fallback(), ext.ServiceUnavailableError
	}

	res, err := b.Execute(
		func() (any, error) {
			return action()
		},
		func() any {
			return fallback()
		},
	)
	if err != nil {
		return fallback(), err
	}

	typedRes, ok := res.(T)
	if !ok {
		return zero, errors.New("type assertion failed")
	}

	return typedRes, nil
}
