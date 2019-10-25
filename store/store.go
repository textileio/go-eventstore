package store

import (
	"errors"
	"fmt"
	"reflect"

	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	es "github.com/textileio/go-eventstore"
)

const (
	idFieldName = "ID"
)

var (
	ErrInvalidModel = errors.New("the model is valid")

	baseKey = ds.NewKey("/model/")
	log     = logging.Logger("store")
)

type Store struct {
	datastore    ds.Datastore
	dispatcher   *es.Dispatcher
	eventcreator es.EventCreator
	models       map[reflect.Type]*Model
}

func NewStore(ds ds.Datastore, dispatcher *es.Dispatcher, eventcreator es.EventCreator) *Store {
	return &Store{
		datastore:    ds,
		dispatcher:   dispatcher,
		models:       make(map[reflect.Type]*Model),
		eventcreator: eventcreator,
	}
}

func (s *Store) Register(name string, t interface{}) (*Model, error) {
	if s.alreadyRegistered(t) {
		return nil, fmt.Errorf("already registered model")
	}

	if !isValidModel(t) {
		return nil, ErrInvalidModel
	}

	m := NewModel(name, t, s.datastore, s.dispatcher, s.eventcreator)
	s.models[m.valueType] = m
	s.dispatcher.Register(m) // ToDo: find good place for reg token, prob will register eventcreator

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
