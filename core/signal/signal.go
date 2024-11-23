package signal

import (
	"github.com/yates-z/easel/core/container/queue"
)

type Signal struct {
	slots *queue.PriorityQueue[Slot]
}

// NewSignal create a new Signal
func NewSignal() *Signal {
	s := &Signal{
		slots: queue.NewPriorityQueue[Slot](queue.WithEqualFunc(func(a, b Slot) bool {
			if a.ID() > 0 && b.ID() > 0 {
				return a.ID() == b.ID()
			}
			return a.Name() == b.Name()
		})),
	}
	return s
}

func (s *Signal) Connect(slot Slot) {
	if slot.ID() != 0 {
		panic("Slot is not reusable!")
	}
	slot.setID(s.slots.Enqueue(slot, slot.Priority()))
}

func (s *Signal) Emit(args ...interface{}) {
	for _, _slot := range s.slots.Iterator() {
		if _slot.Exceeded() {
			s.Disconnect(_slot)
			continue
		}
		go _slot.call(args...)
	}
}

func (s *Signal) Disconnect(slot Slot) bool {
	return s.slots.Remove(slot)
}

func (s *Signal) DisconnectAll() {
	s.slots.Clear()
}
