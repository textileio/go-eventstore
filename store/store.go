package store

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	es "github.com/textileio/go-eventstore"
	"github.com/textileio/go-eventstore/jsonpatcher"
)

const (
	idFieldName = "ID"
)

var (
	ErrInvalidModel = errors.New("the model is valid")

	log = logging.Logger("store")
)

type Store struct {
	lock       sync.RWMutex
	datastore  ds.Datastore
	dispatcher *es.Dispatcher
	models     map[reflect.Type]*Model
}

// NewStore creates a new Store, which will *own* ds and dispatcher for internal use.
// Saying it differently, ds and dispatcher shouldn't be used externally.
func NewStore(ds ds.Datastore, dispatcher *es.Dispatcher) *Store {
	return &Store{
		datastore:  ds,
		dispatcher: dispatcher,
		models:     make(map[reflect.Type]*Model),
	}
}

func (s *Store) RegisterJSONPatcher(name string, defaultInstance interface{}) (*Model, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.alreadyRegistered(defaultInstance) {
		return nil, fmt.Errorf("already registered model")
	}

	if !isValidModel(defaultInstance) {
		return nil, ErrInvalidModel
	}

	eventcreator := jsonpatcher.New()
	m := NewModel(name, defaultInstance, s.datastore, s.dispatcher, eventcreator, s)
	s.models[m.valueType] = m
	s.dispatcher.Register(m)
	return m, nil
}

func (s *Store) readTxn(m *Model, f func(txn *Txn) error) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	txn := &Txn{model: m, readonly: true}
	defer txn.Discard()
	if err := f(txn); err != nil {
		return err
	}
	return nil
}

func (s *Store) writeTxn(m *Model, f func(txn *Txn) error) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	txn := &Txn{model: m}
	defer txn.Discard()
	if err := f(txn); err != nil {
		return err
	}
	return txn.Commit()
}

// Dispatch applies external events to the store. This function guarantee
// no interference with registered model states, and viceversa.
func (s *Store) Dispatch(e es.Event) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.dispatcher.Dispatch(e)
}

func (s *Store) alreadyRegistered(t interface{}) bool {
	valueType := reflect.TypeOf(t)
	_, ok := s.models[valueType]
	return ok
}

func isValidModel(t interface{}) bool {
	v := reflect.ValueOf(t)
	if v.Type().Kind() != reflect.Ptr {
		v = reflect.New(reflect.TypeOf(v))
	}
	return v.Elem().FieldByName(idFieldName).IsValid()
}
