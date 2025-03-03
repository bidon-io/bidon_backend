package clock_test

import (
	"testing"
	"time"

	"github.com/bidon-io/bidon-backend/pkg/clock"
)

// Ensure that the clock's time matches the standard library.

func TestClock_Now(t *testing.T) {
	r := time.Now().Round(time.Second)
	c := clock.New().Now().Round(time.Second)
	if !r.Equal(c) {
		t.Errorf("not equal: %s != %s", r, c)
	}
}

func TestClock_Since(t *testing.T) {
	c := clock.New()
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	duration := c.Since(start)
	if duration <= 0 {
		t.Errorf("expected positive duration, got %v", duration)
	}
}

func TestMockClock_Now(t *testing.T) {
	m := clock.NewMock()
	now := m.Now()
	if !now.Equal(time.Unix(0, 0)) {
		t.Errorf("expected Unix epoch, got %v", now)
	}
}

func TestMockClock_Add(t *testing.T) {
	m := clock.NewMock()
	m.Add(10 * time.Second)
	expected := time.Unix(10, 0)
	if !m.Now().Equal(expected) {
		t.Errorf("expected %v, got %v", expected, m.Now())
	}
}

func TestMockClock_Set(t *testing.T) {
	m := clock.NewMock()
	newTime := time.Unix(100, 0)
	m.Set(newTime)
	if !m.Now().Equal(newTime) {
		t.Errorf("expected %v, got %v", newTime, m.Now())
	}
}

func TestMockClock_Since(t *testing.T) {
	m := clock.NewMock()
	start := time.Unix(0, 0)
	m.Add(10 * time.Second)
	duration := m.Since(start)
	if duration != 10*time.Second {
		t.Errorf("expected 10s, got %v", duration)
	}
}
