package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/alecthomas/jsonschema"
	jsonpatch "github.com/evanphx/json-patch"
	ds "github.com/ipfs/go-datastore"
)

var (
	ErrNotFound = errors.New("instance not found")
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
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

type Txn struct {
	model     *Model
	discarded bool
	commited  bool
	patches   [][]byte
}

func (t *Txn) Discard() {
	t.discarded = true
}

func (t *Txn) Commit() error {
	if t.discarded || t.commited {
		return fmt.Errorf("can't commit discarded or already commited txn")
	}
	// ToDo: Somehow here should merge multiple `.Save()` in the Txn
	// and build whatever necessary to build the Event.
	// Something like folding all t.patches into something for the Event.Body
	return nil
}

func (t *Txn) Add(id string, v interface{}) error {
	key := ds.NewKey(id)
	exists, err := t.model.datastore.Has(key)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("can't add already existing instance id:%s", id)
	}
	return t.persist(key, nil, v)
}

func (t *Txn) Save(id string, v interface{}) error {
	key := ds.NewKey(id)
	actual, err := t.model.datastore.Get(key)
	if err == ds.ErrNotFound {
		return fmt.Errorf("can't save non-existing instance id:%s", id)
	}
	if err != nil {
		return err
	}
	return t.persist(key, actual, v)
}

func (t *Txn) Delete(id string) error {
	key := ds.NewKey(id)
	exists, err := t.model.datastore.Has(key)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("can't remove since doesn't exist: %s", id)
	}
	// ToDo
	panic("Not implemented yet")
}

func (t *Txn) FindByID(id string, v interface{}) error {
	return t.model.FindByID(id, v)
}

func (t *Txn) persist(key ds.Key, actual []byte, new interface{}) error {
	// ToDo: Validate `v` against t.model.schema?
	newBytes, err := json.Marshal(new)
	if err != nil {
		return err
	}
	if actual != nil {
		patch, err := jsonpatch.CreateMergePatch(actual, newBytes)
		if err != nil {
			return err
		}
		t.patches = append(t.patches, patch)
		fmt.Printf("Save() patch (%d bytes): %s\n\n", len(patch), string(patch))
	} else {
		fmt.Printf("Add(): %s\n\n", string(newBytes))
	}
	return t.model.datastore.Put(key, newBytes)
}
