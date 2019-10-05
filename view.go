package eventstore

import (
	"github.com/textileio/go-eventstore/broadcast"
)

// Reducer does some stuff...
type Reducer interface {
	Reduce(event Event) error
}

// ViewModel does some stuff...
type ViewModel interface {
	Reducer
	Listen() *broadcast.Listener
}
