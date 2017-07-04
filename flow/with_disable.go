package flow

import (
	"sync"

	floc "github.com/workanator/go-floc"
)

type disablableFlow struct {
	sync.Mutex
	parent   floc.Flow
	disabled bool
}

// DisableFunc when invoked disables calls to Complete and Cancel.
type DisableFunc func()

// WithDisable creates a new instance of the flow containing the parent flow
// and a disable function which allows to disable calls to Complete and Cancel.
func WithDisable(parent floc.Flow) (floc.Flow, DisableFunc) {
	flow := &disablableFlow{
		parent: parent,
	}

	disable := func() {
		flow.Mutex.Lock()
		defer flow.Mutex.Unlock()

		flow.disabled = true
	}

	return flow, disable
}

// Done returns a channel that's closed when the flow done.
// Successive calls to Done return the same value.
func (f *disablableFlow) Done() <-chan struct{} {
	return f.parent.Done()
}

// Close finishes the flow and releases all underlying resources.
func (f *disablableFlow) Close() {
	f.parent.Cancel(nil)
}

// Complete finishes the flow with success status and stops
// execution of futher jobs if any.
func (f *disablableFlow) Complete(data interface{}) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	if !f.disabled {
		f.parent.Complete(data)
	}
}

// Cancel cancels the execution of the flow.
func (f *disablableFlow) Cancel(data interface{}) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	if !f.disabled {
		f.parent.Cancel(data)
	}
}

// Tests if the execution of the flow is either completed or canceled.
func (f *disablableFlow) IsFinished() bool {
	return f.parent.IsFinished()
}

// Returns the result code and the result data of the flow.
func (f *disablableFlow) Result() (result floc.Result, data interface{}) {
	return f.parent.Result()
}