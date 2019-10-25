package store

import (
	"errors"
	"fmt"
	"reflect"

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
	datastore  ds.Datastore
	dispatcher *es.Dispatcher
	models     map[reflect.Type]*Model
}

func NewStore(ds ds.Datastore, dispatcher *es.Dispatcher) *Store {
	return &Store{
		datastore:  ds,
		dispatcher: dispatcher,
		models:     make(map[reflect.Type]*Model),
	}
}

func (s *Store) RegisterJSONPatcher(name string, defaultInstance interface{}) (*Model, error) {
	if s.alreadyRegistered(defaultInstance) {
		return nil, fmt.Errorf("already registered model")
	}

	if !isValidModel(defaultInstance) {
		return nil, ErrInvalidModel
	}

	eventcreator := jsonpatcher.New()
	m := NewModel(name, defaultInstance, s.datastore, s.dispatcher, eventcreator)
	s.models[m.valueType] = m
	s.dispatcher.Register(m) // ToDo: find good place for reg token, with proper unregistering

	// Debug
	// dbgJSON, _ := json.MarshalIndent(m.schema, "", "  ")
	// log.Debugf("registered model %q: %s", name, string(dbgJSON))

	return m, nil
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
