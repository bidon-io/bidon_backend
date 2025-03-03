package clock

import (
	"sync"
	"time"
)

// Clock represents an interface to the functions in the standard library time
// package. Two implementations are available in the clock package. The first
// is a real-time clock which simply wraps the time package's functions. The
// second is a mock clock which will only change when programmatically adjusted.
// Contains only the methods that are used in the project. Feel free to add more methods if needed.
type Clock interface {
	Now() time.Time
	Since(t time.Time) time.Duration
}

// New returns an instance of a real-time clock.
func New() Clock {
	return &clock{}
}

// clock implements a real-time clock by simply wrapping the time package functions.
type clock struct{}

func (c *clock) Now() time.Time { return time.Now() }

func (c *clock) Since(t time.Time) time.Duration { return time.Since(t) }

// Mock represents a mock clock that only moves forward programmatically.
// It can be preferable to a real-time clock when testing time-based functionality.
type Mock struct {
	mu  sync.RWMutex
	now time.Time // current time
}

// NewMock returns an instance of a mock clock.
// The current time of the mock clock on initialization is the Unix epoch.
func NewMock() *Mock {
	return &Mock{now: time.Unix(0, 0)}
}

// Add moves the current time of the mock clock forward by the specified duration.
func (m *Mock) Add(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.now = m.now.Add(d)
}

// Set sets the current time of the mock clock to a specific one.
func (m *Mock) Set(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.now = t
}

// Now returns the current wall time on the mock clock.
func (m *Mock) Now() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.now
}

// Since returns time since `t` using the mock clock's wall time.
func (m *Mock) Since(t time.Time) time.Duration {
	return m.Now().Sub(t)
}

var (
	// type checking
	_ Clock = &Mock{}
)
