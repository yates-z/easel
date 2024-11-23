package signal

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

type SlotOption func(*slot)

type Slot interface {
	setID(int64) bool
	ID() int64
	Name() string
	Priority() int
	MaxCalls() int64
	Exceeded() bool
	call(args ...any)
}

var _ Slot = (*slot)(nil)

func WithName(name string) SlotOption {
	return func(s *slot) {
		s.name = name
	}
}

func WithPriority(priority int) SlotOption {
	return func(s *slot) {
		s.priority = priority
	}
}

func WithMaxCalls(maxCalls int64) SlotOption {
	return func(s *slot) {
		s.maxCalls = maxCalls
	}
}

type slot struct {
	id        atomic.Int64
	name      string
	priority  int
	maxCalls  int64
	callCount atomic.Int64
	callback  reflect.Value
}

func SLOT(callback any, opts ...SlotOption) Slot {
	fn := reflect.ValueOf(callback)
	if fn.Kind() != reflect.Func {
		panic("callback must be a function")
	}

	s := &slot{priority: 0, callback: fn}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// setID implements Slot.
func (s *slot) setID(id int64) bool {
	//return s.id.CompareAndSwap(0, id)
	s.id.Store(id)
	return true
}

// ID implements Slot.
func (s *slot) ID() int64 {
	return s.id.Load()
}

// Name implements Slot.
func (s *slot) Name() string {
	return s.name
}

// Priority implements Slot.
func (s *slot) Priority() int {
	return s.priority
}

// MaxCalls implements Slot.
func (s *slot) MaxCalls() int64 {
	return s.maxCalls
}

func (s *slot) Exceeded() bool {
	if s.maxCalls > 0 && s.callCount.Load() >= s.maxCalls {
		return true
	}
	return false
}

func (s *slot) call(args ...any) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error in slot %d: %v\n", s.ID(), r)
		}
	}()

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	s.callback.Call(in)

	s.callCount.Add(1)
}
