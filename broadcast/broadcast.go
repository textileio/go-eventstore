// Package broadcast implements multi-listener broadcast channels.
// See https://godoc.org/github.com/tjgq/broadcast for original implementation.
//
// To create an un-buffered broadcast channel, just declare a Broadcaster:
//
//     var b broadcast.Broadcaster
//
// To create a buffered broadcast channel with capacity n, call New:
//
//     b := broadcast.New(n)
//
// To add a listener to a channel, call Listen and read from Channel():
//
//     l := b.Listen()
//     for v := range l.Channel() {
//         // ...
//     }
//
//
// To send to the channel, call Send:
//
//     b.Send("Hello world!")
//     v <- l.Channel() // returns interface{}("Hello world!")
//
// To remove a listener, call Discard.
//
//     l.Discard()
//
// To close the broadcast channel, call Discard. Any existing or future listeners
// will read from a closed channel:
//
//     b.Discard()
//     v, ok <- l.Channel() // returns ok == false
package broadcast

import "sync"

// Broadcaster implements a Publisher. The zero value is a usable un-buffered channel.
type Broadcaster struct {
	m         sync.Mutex
	listeners map[int]chan<- interface{} // lazy init
	nextID    int
	capacity  int
	closed    bool
}

// NewBroadcaster returns a new Broadcaster with the given capacity (0 means un-buffered).
func NewBroadcaster(n int) *Broadcaster {
	return &Broadcaster{capacity: n}
}

// Send broadcasts a message to the channel.
// Sending on a closed channel causes a runtime panic.
func (b *Broadcaster) Send(v interface{}) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.closed {
		panic("broadcast: send after close")
	}
	for _, l := range b.listeners {
		l <- v
	}
}

// Discard closes the channel, disabling the sending of further messages.
func (b *Broadcaster) Discard() {
	b.m.Lock()
	defer b.m.Unlock()
	b.closed = true
	for _, l := range b.listeners {
		close(l)
	}
}

// Listen returns a Listener for the broadcast channel.
func (b *Broadcaster) Listen() *Listener {
	b.m.Lock()
	defer b.m.Unlock()
	if b.listeners == nil {
		b.listeners = make(map[int]chan<- interface{})
	}
	for b.listeners[b.nextID] != nil {
		b.nextID++
	}
	ch := make(chan interface{}, b.capacity)
	if b.closed {
		close(ch)
	}
	b.listeners[b.nextID] = ch
	return &Listener{ch, b, b.nextID}
}

// Listener implements a Subscriber to broadcast channel.
type Listener struct {
	ch <-chan interface{}
	b  *Broadcaster
	id int
}

// Discard closes the Listener, disabling the reception of further messages.
func (l *Listener) Discard() {
	l.b.m.Lock()
	defer l.b.m.Unlock()
	delete(l.b.listeners, l.id)
}

// Channel returns the channel that receives broadcast messages
func (l *Listener) Channel() <-chan interface{} {
	return l.ch
}
