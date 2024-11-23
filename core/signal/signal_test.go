package signal

import (
	"fmt"
	"testing"
	"time"
)

func Test_Concurrency(t *testing.T) {
	s := SLOT(func() {})
	signal := NewSignal()
	for i := 0; i < 10; i++ {
		go func() {
			signal.Connect(s)
		}()
	}
	time.Sleep(2 * time.Second)
	fmt.Println(s.ID())
}

func Test_Concurrency2(t *testing.T) {
	signal := NewSignal()
	for i := range 10 {
		go func() {
			s := SLOT(func() {})
			signal.Connect(s)
		}()
		go func() {
			if i != 5 {
				return
			}
			signal.DisconnectAll()
		}()
	}

	time.Sleep(1 * time.Second)
	// id should start from 0.
	for _, s := range signal.slots.Iterator() {
		fmt.Println(s.ID())
	}
}

func Test_Priority(t *testing.T) {
	s1 := SLOT(func() { fmt.Println("I am slot1") }, WithPriority(1))
	s2 := SLOT(func() { fmt.Println("I am slot2") }, WithPriority(2))
	signal := NewSignal()
	signal.Connect(s1)
	signal.Connect(s2)

	signal.Emit()
	time.Sleep(1 * time.Second)
}

func Test_Signal(t *testing.T) {
	mySignal := NewSignal()

	s1 := SLOT(func() { fmt.Println("I am slot1") }, WithPriority(1))
	s2 := SLOT(func() { fmt.Println("I am slot2") }, WithPriority(2))

	mySignal.Connect(s1)
	mySignal.Connect(s2)
	// 触发信号
	mySignal.Emit()

	time.Sleep(1 * time.Second)
	mySignal.Disconnect(s1)
	mySignal.Emit()
	time.Sleep(1 * time.Second)
}
