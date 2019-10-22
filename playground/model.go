package playground

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"github.com/alecthomas/jsonschema"
	ds "github.com/ipfs/go-datastore"
)

type Model struct {
	lock      sync.RWMutex
	schema    *jsonschema.Schema
	valueType reflect.Type
	datastore ds.Datastore
}

func (m *Model) Update(f func(txn *Txn) error) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	txn := &Txn{model: m}
	defer txn.Discard()
	if err := f(txn); err != nil {
		return err
	}
	return txn.Commit()
}

func (m *Model) FindByID(id string, v interface{}) error {
	key := ds.NewKey(id)
	bytes, err := m.datastore.Get(key)
	if errors.Is(err, ds.ErrNotFound) {
		v = nil
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}
