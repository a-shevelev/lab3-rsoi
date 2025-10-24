package circuit

import (
	"sync"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type Breaker struct {
	mu            sync.Mutex
	state         State
	failureCount  int
	threshold     int
	retryAfter    time.Duration
	lastFailTime  time.Time
	onStateChange func(State)
}

func NewBreaker(threshold int, retryAfter time.Duration) *Breaker {
	return &Breaker{
		state:      Closed,
		threshold:  threshold,
		retryAfter: retryAfter,
	}
}

func (b *Breaker) Execute(
	operation func() (any, error),
	fallback func() any,
) (any, error) {

	b.mu.Lock()
	switch b.state {
	case Open:
		if time.Since(b.lastFailTime) < b.retryAfter {
			b.mu.Unlock()
			return fallback(), nil
		}

		b.state = HalfOpen
	case Closed, HalfOpen:
	default:
		panic("unhandled default case")
	}
	b.mu.Unlock()

	result, err := operation()
	if err != nil {
		b.recordFailure()
		return fallback(), err
	}
	b.reset()
	return result, nil
}

func (b *Breaker) recordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failureCount++
	b.lastFailTime = time.Now()
	if b.failureCount >= b.threshold {
		b.state = Open
		if b.onStateChange != nil {
			b.onStateChange(b.state)
		}
	}
}

func (b *Breaker) reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failureCount = 0
	b.state = Closed
	if b.onStateChange != nil {
		b.onStateChange(b.state)
	}

}
