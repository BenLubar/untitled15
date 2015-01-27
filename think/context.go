package think

import (
	"container/heap"
	"time"
)

// Thinker is the main interface for package think.
type Thinker interface {
	// Think is called by the worker goroutines of Context. Only the given
	// time should be used to update state. Calling time.Now and related
	// functions will break mocking.
	Think(*Context, time.Time)
}

type thinker struct {
	Thinker
	next time.Time
}

// Registration is a helper type returned by the Close method of Context.
type Registration struct {
	Thinker
	time.Duration
}

// Context manages thinkers.
type Context struct {
	register chan thinker
	done     chan chan<- []Registration
}

// NewContext returns a new context with the given number of worker goroutines
// and the given timer.
func NewContext(workers int, timer Timer) *Context {
	var c Context
	c.register = make(chan thinker)
	c.done = make(chan chan<- []Registration)
	go c.master(workers, timer)
	return &c
}

// Register queues a Thinker to be run at the given time. Register panics if
// Close has been called and returned, but it is safe to call Register from
// a Thinker at any time. It is safe to call Register from any goroutine if
// Close has not been called.
func (c *Context) Register(t Thinker, next time.Time) {
	c.register <- thinker{t, next}
}

// Close waits for any currently running Thinkers to finish, then returns
// a sorted slice of thinkers that were registered but did not yet Think.
// After Close has returned, no method on Context may be called. Once Close
// has been called, it is no longer safe to call methods on Context from
// anywhere other than a Thinker.
func (c *Context) Close() []Registration {
	ch := make(chan []Registration)
	c.done <- ch
	return <-ch
}

func (c *Context) master(workers int, timer Timer) {
	var h thinkHeap
	now := timer.Now()

	input, done := make(chan thinker), make(chan struct{}, workers)
	for i := 0; i < workers; i++ {
		go c.slave(input, done)
	}

	var wait <-chan time.Time
	var output chan<- []Registration
	for {
		var in chan<- thinker
		var next thinker
		if len(h) != 0 {
			if !now.Before(h[0].next) {
				in, next = input, h[0]
			} else if wait == nil {
				wait = timer.Wait(h[0].next)
			}
		}

		select {
		case t := <-c.register:
			heap.Push(&h, t)

		case in <- next:
			heap.Pop(&h)

		case output = <-c.done:
			close(c.done)
			close(input)
			input = nil

		case <-done:
			if workers--; workers == 0 {
				close(c.register)
				var r []Registration
				for len(h) != 0 {
					t := heap.Pop(&h).(thinker)
					r = append(r, Registration{t.Thinker, t.next.Sub(timer.Now())})
				}
				output <- r
				return
			}

		case now = <-wait:
			wait = nil
		}
	}
}

func (c *Context) slave(input <-chan thinker, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	for t := range input {
		t.Think(c, t.next)
	}
}
