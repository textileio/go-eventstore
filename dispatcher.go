package eventstore

import (
	"bytes"
	"encoding/gob"
	"sync"

	"context"
	"fmt"

	datastore "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"golang.org/x/sync/errgroup"
)

type Reducer interface {
	Reduce(event Event) error
}

// @todo: Should we also support a `Transformer` to actually fetch raw event data as part of a pipeline?

// Token is a simple unique ID used to reference a registered callback.
type Token string

const prefix = "ID"

var lastID = 0

// Dispatcher is used to dispatch events to registered reducers.
//
// This is different from generic pub-sub systems because reducers are not subscribed to particular events.
// Every event is dispatched to every registered reducer. When a given reducer is registered, it returns a `token`,
// which can be used to deregister the reducer later.
type Dispatcher struct {
	store    datastore.TxnDatastore
	reducers map[Token]Reducer
	lock     sync.Mutex
}

// NewDispatcher creates a new EventDispatcher
func NewDispatcher(store datastore.TxnDatastore) *Dispatcher {
	return &Dispatcher{
		store:    store,
		reducers: make(map[Token]Reducer),
	}
}

// Store returns the internal event store.
func (d *Dispatcher) Store() datastore.TxnDatastore {
	return d.store
}

// Register takes a reducer to be invoked with each dispatched event and returns a token for de-registration.
func (d *Dispatcher) Register(reducer Reducer) Token {
	d.lock.Lock()
	defer d.lock.Unlock()
	lastID++
	id := Token(fmt.Sprintf("%s-%d", prefix, lastID))
	d.reducers[id] = reducer
	return id
}

// Deregister removes a reducer based on its token.
func (d *Dispatcher) Deregister(token Token) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if _, ok := d.reducers[token]; !ok {
		return fmt.Errorf("`%s` does not map to a registered callback", token)
	}
	delete(d.reducers, token)
	return nil
}

// Dispatch dispatches a payload to all registered reducers.
func (d *Dispatcher) Dispatch(event Event) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	// Key format: <timestamp>/<entity-id>/<type>
	// @todo: This is up for debate, its a 'fake' Event struct right now anyway
	key := datastore.NewKey(string(event.Time())).ChildString(event.EntityID().String()).ChildString(event.Type())
	// Encode and add an Event to event store
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(event); err != nil {
		return err
	}
	if err := d.Store().Put(key, b.Bytes()); err != nil {
		return err
	}
	// Safe to fire off reducers now that event is persisted
	g, _ := errgroup.WithContext(context.Background())
	for _, reducer := range d.reducers {
		// Launch each reducer in a separate goroutine
		g.Go(func() error {
			return reducer.Reduce(event)
		})
	}
	// Wait for all reducers to complete or error out
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// Query searches the internal event store and returns a query result.
// This is a syncronouse version of github.com/ipfs/go-datastore's Query method
func (d *Dispatcher) Query(query query.Query) ([]query.Entry, error) {
	result, err := d.store.Query(query)
	if err != nil {
		return nil, err
	}
	return result.Rest()
}
