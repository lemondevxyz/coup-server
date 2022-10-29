package game

import (
	"sync"
	"time"
)

type Notifier struct {
	mtx   sync.RWMutex
	chnls map[time.Time]chan struct{}
	val   interface{}
}

func (n *Notifier) Subscribe() (time.Time, chan struct{}) {
	now, ch := time.Now(), make(chan struct{})

	n.mtx.Lock()
	n.chnls[now] = ch
	n.mtx.Unlock()

	return now, ch
}

func (n *Notifier) Unsubscribe(val time.Time) {
	n.mtx.Lock()
	delete(n.chnls, val)
	n.mtx.Unlock()
}

func (n *Notifier) Announce() {
	n.mtx.RLock()
	val := struct{}{}
	for _, v := range n.chnls {
		v <- val
	}
	n.mtx.RUnlock()
}

func (n *Notifier) Get() interface{} {
	return n.val
}

func (n *Notifier) Set(val interface{}) {
	n.mtx.Lock()
	n.val = val
	n.mtx.Unlock()
}

func NewNotifier() *Notifier {
	return &Notifier{
		chnls: map[time.Time]chan struct{}{},
	}
}
