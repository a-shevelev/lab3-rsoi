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
	mu                sync.Mutex
	state             State
	failureTimes      []time.Time
	threshold         int
	retryAfter        time.Duration
	failureTimer      time.Duration
	halfOpenLimit     int
	openTime          time.Time
	halfOpenSuccesses int
	onStateChange     func(State)
}

func NewBreaker(threshold int, retryAfter, window time.Duration, halfOpenLimit int) *Breaker {
	return &Breaker{
		state:         Closed,
		threshold:     threshold,
		retryAfter:    retryAfter,
		failureTimer:  window,
		halfOpenLimit: halfOpenLimit,
		failureTimes:  make([]time.Time, 0),
	}
}

func (b *Breaker) clearOldFailures() {
	now := time.Now()
	validFailures := make([]time.Time, 0)

	for _, t := range b.failureTimes {
		if now.Sub(t) <= b.failureTimer {
			validFailures = append(validFailures, t)
		}
	}

	b.failureTimes = validFailures
}

func (b *Breaker) getFailureCount() int {
	b.clearOldFailures()
	return len(b.failureTimes)
}

func (b *Breaker) Execute(
	operation func() (any, error),
	fallback func() any,
) (any, error) {

	b.mu.Lock()
	switch b.state {
	case Open:
		if time.Since(b.openTime) < b.retryAfter {
			b.mu.Unlock()
			return fallback(), nil
		}
		b.state = HalfOpen
		b.halfOpenSuccesses = 0
	case HalfOpen:
	case Closed:
	default:
		panic("unhandled default case")
	}
	b.mu.Unlock()

	result, err := operation()
	if err != nil {
		b.recordFailure()
		return fallback(), err
	}

	b.recordSuccess()
	return result, nil
}

func (b *Breaker) recordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	b.failureTimes = append(b.failureTimes, now)
	b.clearOldFailures()

	if b.state == HalfOpen {
		b.state = Open
		b.openTime = now
		b.halfOpenSuccesses = 0
	} else if b.getFailureCount() >= b.threshold {
		b.state = Open
		b.openTime = now
	}

	if b.onStateChange != nil {
		b.onStateChange(b.state)
	}
}

func (b *Breaker) recordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case HalfOpen:
		b.halfOpenSuccesses++
		if b.halfOpenSuccesses >= b.halfOpenLimit {
			b.state = Closed
			b.halfOpenSuccesses = 0
		}
	case Closed:
		b.clearOldFailures()
	}

	if b.onStateChange != nil {
		b.onStateChange(b.state)
	}
}
