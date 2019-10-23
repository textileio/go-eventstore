package store

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/alecthomas/jsonschema"
	ds "github.com/ipfs/go-datastore"
	kt "github.com/ipfs/go-datastore/keytransform"
	logging "github.com/ipfs/go-log"
	"github.com/textileio/go-eventstore"
)

var (
	log = logging.Logger("store")
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
	valueType := reflect.TypeOf(t)
	if _, ok := s.models[valueType]; ok {
		return nil, fmt.Errorf("already registered model")
	}

	baseKey := ds.NewKey("/model/").ChildString(name)
	pair := &kt.Pair{
		Convert: func(k ds.Key) ds.Key {
			return baseKey.Child(k)
		},
		Invert: func(k ds.Key) ds.Key {
			l := k.List()
			if !k.IsDescendantOf(baseKey) {
				panic("huh!!") // ToDo
			}
			return ds.KeyWithNamespaces(l[2:])
		},
	}

	m := &Model{
		schema:     jsonschema.Reflect(t),
		datastore:  kt.Wrap(s.datastore, pair), // Make models don't worry about namespaces
		valueType:  valueType,
		dispatcher: s.dispatcher,
	}
	s.models[valueType] = m
	regToken := s.dispatcher.Register(m)
	m.dispatcherToken = regToken

	actualJSON, _ := json.MarshalIndent(m.schema, "", "  ")
	log.Debugf("registered model %q: %s", name, string(actualJSON))

	return m, nil
}
