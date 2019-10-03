package eventstore

import (
	"errors"

	datastore "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"

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

// MemoryModel does stuff...
type MemoryModel struct {
	broadcaster *broadcast.Broadcaster
}

// Reduce does stuff...
func (m *MemoryModel) Reduce(event Event) error {
	return errors.New("not implemented")
}

// Listen does stuff...
func (m *MemoryModel) Listen() *broadcast.Listener {
	return m.broadcaster.Listen()
}

// StoredModel does stuff...
type StoredModel struct {
	*MemoryModel
	store datastore.Datastore
}

// Store returns the internal view store.
func (m StoredModel) Store() datastore.Datastore {
	return m.store
}

// Query searches the internal view store and returns a query result.
// This is a syncronouse version of github.com/ipfs/go-datastore's Query method
func (m StoredModel) Query(query query.Query) ([]query.Entry, error) {
	result, err := m.store.Query(query)
	if err != nil {
		return nil, err
	}
	return result.Rest()
}

// Sanity check
var _ ViewModel = (*MemoryModel)(nil)
var _ ViewModel = (*StoredModel)(nil)
