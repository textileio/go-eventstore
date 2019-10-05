package eventstore

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

// @todo: Should we also support a `Transformer` to actually fetch raw event data as part of a pipeline?

// Token is a simple unique ID used to reference a registered callback.
type Token uint

// lastID is the last ID used by the singleton Dispatcher.
var lastID uint

// Dispatcher is used to dispatch events to registered reducers.
//
// This is different from generic pub-sub systems because reducers are not subscribed to particular events.
// Every event is dispatched to every registered reducer. When a given reducer is registered, it returns a `token`,
// which can be used to deregister the reducer later. The dispatcher wraps an underlying event store, which it uses to
// persist events before dispatching to registered reducers. The whole system uses a mutex to prevent additional
// dispatch calls from kicking off reducers until all in-flight reducers complete.
type Dispatcher struct {
	store    datastore.Datastore
	reducers map[Token]Reducer
	lock     sync.Mutex
}

// NewDispatcher returns a new Dispatcher. This should only be called once in an application to ensure a singleton
// displatcher. While it is not enforced here, it is a good idea to a singleton pattern in your own code, for example:
// 	var (
// 		once sync.Once
// 		instance *Dispatcher
// 	)
//
//	func New() *Dispatcher {
// 		once.Do(func() {
// 			singleton = NewDispatcher(store)
// 		})
// 		return singleton
// 	}
func NewDispatcher(store datastore.Datastore) *Dispatcher {
	return &Dispatcher{
		store:    store,
		reducers: make(map[Token]Reducer),
	}
}

// Register takes a reducer to be invoked with each dispatched event and returns a token for de-registration.
func (d *Dispatcher) Register(reducer Reducer) Token {
	d.lock.Lock()
	defer d.lock.Unlock()
	lastID++
	id := Token(lastID)
	d.reducers[id] = reducer
	return id
}

// Deregister removes a reducer based on its token. If the token is invalid (i.e. no associated reducer),
// this is a no-op.
func (d *Dispatcher) Deregister(token Token) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.reducers, token)
}

// Dispatch dispatches a payload to all registered reducers. It returns a multierror object, which may contain
// zero (nil) or more errors. Errors from reducer callbacks can be safely ignored and/or retried (and are prefixed
// with "warning"), whereas errors due to event persistence are "critical" and will be prefixed as such.
func (d *Dispatcher) Dispatch(event Event) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	var result *multierror.Error
	// Key format: <timestamp>/<entity-id>/<type>
	// @todo: This is up for debate, its a 'fake' Event struct right now anyway
	key := datastore.NewKey(string(event.Time())).ChildString(event.EntityID()).ChildString(event.Type())
	// Encode and add an Event to event store
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(event.Body()); err != nil {
		return multierror.Append(multierror.Prefix(err, "critical"))
	}
	if err := d.store.Put(key, b.Bytes()); err != nil {
		return multierror.Append(multierror.Prefix(err, "critical"))
	}
	// Fire off reducers now that event is safely persisted
	wg := sync.WaitGroup{}
	wg.Add(len(d.reducers))
	for _, reducer := range d.reducers {
		// Launch each reducer in a separate goroutine
		go func(r Reducer) {
			if err := r.Reduce(event); err != nil {
				result = multierror.Append(result, multierror.Prefix(err, "warning"))
			}
			wg.Done()
		}(reducer)
	}
	// Wait for all reducers to complete (or error) before releasing our lock
	wg.Wait()
	return result.ErrorOrNil()
}

// Query searches the internal event store and returns a query result. This is a proxy to the underlying event store's
// Query method.
func (d *Dispatcher) Query(query query.Query) (query.Results, error) {
	return d.store.Query(query)
}
