package store

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/alecthomas/jsonschema"
	ds "github.com/ipfs/go-datastore"
	kt "github.com/ipfs/go-datastore/keytransform"
	logging "github.com/ipfs/go-log"
	"github.com/textileio/go-eventstore"
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
	datastore  ds.Datastore
	dispatcher *eventstore.Dispatcher
	models     map[reflect.Type]*Model
}

func NewStore(ds ds.Datastore, dispatcher *eventstore.Dispatcher) *Store {
	return &Store{
		datastore:  ds,
		dispatcher: dispatcher,
		models:     make(map[reflect.Type]*Model),
	}
}

func (s *Store) Register(name string, t interface{}) (*Model, error) {
	if s.alreadyRegistered(t) {
		return nil, fmt.Errorf("already registered model")
	}
	if !isValidModel(t) {
		return nil, ErrInvalidModel
	}

	m := s.createModel(name, t)
	s.models[m.valueType] = m

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

func (s *Store) createModel(name string, t interface{}) *Model {
	baseModelKey := baseKey.ChildString(name)
	pair := &kt.Pair{
		Convert: func(k ds.Key) ds.Key {
			return baseModelKey.Child(k)
		},
		Invert: func(k ds.Key) ds.Key {
			l := k.List()
			if !k.IsDescendantOf(baseModelKey) {
				panic("huh!!") // ToDo: Reconsider the keytransformation thing. may backfire in queries, see later.
			}
			return ds.KeyWithNamespaces(l[2:])
		},
	}
	m := &Model{
		schema:     jsonschema.Reflect(t),
		datastore:  kt.Wrap(s.datastore, pair), // Make models don't worry about namespaces
		valueType:  reflect.TypeOf(t),
		dispatcher: s.dispatcher,
	}
	m.regToken = s.dispatcher.Register(m)

	return m
}
