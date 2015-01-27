package think

import "time"

// Timer is a wrapper around the time package's functions.
type Timer interface {
	// Now returns the current time.
	Now() time.Time

	// Wait returns a channel that recieves the current time when it is
	// not before the given time.
	Wait(time.Time) <-chan time.Time
}

// RealTimer is a Timer that calls the functions of the time package directly.
var RealTimer realTimer

type realTimer struct{}

// Now implements Timer.
func (realTimer) Now() time.Time {
	return time.Now()
}

// Wait implements Timer.
func (realTimer) Wait(until time.Time) <-chan time.Time {
	d := until.Sub(time.Now())
	if d < 0 {
		d = 0
	}
	return time.After(d)
}

// FakeTimer is a Timer that advances instantly.
type FakeTimer struct {
	now time.Time
}

// Now implements Timer.
func (t *FakeTimer) Now() time.Time {
	if t.now.IsZero() {
		t.now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	return t.now
}

// Wait implements Timer.
func (t *FakeTimer) Wait(until time.Time) <-chan time.Time {
	ch := make(chan time.Time, 1)
	if t.Now().Before(until) {
		t.now = until
	}
	ch <- t.now
	return ch
}
