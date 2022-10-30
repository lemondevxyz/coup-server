package game

import (
	"sync"
	"time"
)

// Notifier is a structure that allows for a Subscriber model.
//
// Essentially, it allows other functions to Subscribe to any changes to
// the underlying value.
//
// This is useful because it allows the structure, Game, to give out
// announcements about which turn is it; or if a Claim has finished
// or not.
//
// An empty notifier will panic.
type Notifier struct {
	mtx   sync.RWMutex
	chnls map[time.Time]chan struct{}
	val   interface{}
}

// Subscribe returns a timestamp and a channel. The timestamp functions
// as the ID of this particular channel. So that in the future, one
// can Unsubscribe with that particular timestamp.
//
// Do note: Unused subscribed channels can freeze a Notifier. Use
//          Notifier.Unsubscribe when you're done using a channel.
func (n *Notifier) Subscribe() (time.Time, chan struct{}) {
	now, ch := time.Now(), make(chan struct{})

	n.mtx.Lock()
	n.chnls[now] = ch
	n.mtx.Unlock()

	return now, ch
}

// Unsubscribe unsubscribers the channel by its timestamp.
func (n *Notifier) Unsubscribe(val time.Time) {
	n.mtx.Lock()
	delete(n.chnls, val)
	n.mtx.Unlock()
}

// Announce sends out an event to every subscribed channel.
func (n *Notifier) Announce() {
	n.mtx.RLock()
	val := struct{}{}
	for _, v := range n.chnls {
		v <- val
	}
	n.mtx.RUnlock()
}

// Get returns the underlying value
func (n *Notifier) Get() interface{} {
	return n.val
}

// Set sets the underlying value
//
// Do note: This does not call Announce; you have to call Announce by
//          yourself to trigger all channels.
func (n *Notifier) Set(val interface{}) {
	n.mtx.Lock()
	n.val = val
	n.mtx.Unlock()
}

// NewNotifier returns a valid Notifier.
func NewNotifier() *Notifier {
	return &Notifier{
		chnls: map[time.Time]chan struct{}{},
	}
}
