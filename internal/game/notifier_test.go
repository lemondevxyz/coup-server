package game

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestNewNotifier(t *testing.T) {
	n := NewNotifier()
	is := is.New(t)
	is.True(n.chnls != nil)
}

func TestNotifierSubscribe(t *testing.T) {
	n := NewNotifier()

	id, _ := n.Subscribe()
	is := is.New(t)
	is.True(n.chnls[id] != nil)
}

func TestNotifierUnsubscribe(t *testing.T) {
	n := NewNotifier()

	now := time.Now()
	n.chnls[now] = nil
	n.Unsubscribe(now)

	is := is.New(t)
	_, ok := n.chnls[now]
	is.True(ok == false)
}

func TestNotifierAnnounce(t *testing.T) {
	wait := make(chan struct{})

	notifier := NewNotifier()
	_, ch := notifier.Subscribe()
	go func() {
		val := <-ch
		wait <- val
	}()
	notifier.Announce()

	select {
	case <-wait:
	case <-time.After(time.Millisecond):
		t.Fatalf("Annonuce doesn't work")
	}
}

func TestNotifierGet(t *testing.T) {
	val := interface{}(false)
	notifier := NewNotifier()
	notifier.val = val

	is := is.New(t)

	is.Equal(notifier.Get(), val)
}

func TestNotifierSet(t *testing.T) {
	val := interface{}(false)
	notifier := NewNotifier()
	notifier.Set(val)

	is := is.New(t)
	is.Equal(notifier.val, val)
}
