package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/alecthomas/jsonschema"
	jsonpatch "github.com/evanphx/json-patch"
	ds "github.com/ipfs/go-datastore"
	"github.com/textileio/go-eventstore"
)

var (
	ErrNotFound = errors.New("instance not found")
)

type operationType string

const (
	upsert operationType = "upsert"
	delete operationType = "delete"
)

type Model struct {
	lock            sync.RWMutex
	schema          *jsonschema.Schema
	valueType       reflect.Type
	datastore       ds.Datastore
	dispatcher      *eventstore.Dispatcher
	dispatcherToken eventstore.Token
}

func (m *Model) Update(f func(txn *Txn) error) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	txn := &Txn{model: m, ops: make(map[ds.Key]operation)}
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

func (m *Model) Reduce(event eventstore.Event) error {
	log.Debug("Reduce() start")

	log.Debug("Reduce() end")
	return nil
}

type Txn struct {
	model     *Model
	discarded bool
	commited  bool
	ops       map[ds.Key]operation
}

type operation struct {
	opType    operationType
	entityID  string
	jsonPatch []byte
}

func (t *Txn) Discard() {
	t.discarded = true
}

func (t *Txn) Commit() error {
	if t.discarded || t.commited {
		return fmt.Errorf("can't commit discarded/commited txn")
	}
	now := time.Now()

	//  ToDo/Important: As first approximation, each key change is a separate event
	for _, op := range t.ops {
		event := eventstore.NewJsonPatchEvent(now, op.entityID, t.model.schema.Ref, op.jsonPatch)
		if err := t.model.dispatcher.Dispatch(event); err != nil {
			return err // Ugh! partial failure, think about what this means for application state
		}
	}
	return nil
}

func (t *Txn) Add(id string, new interface{}) error {
	key := ds.NewKey(id)
	exists, err := t.model.datastore.Has(key)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("can't add already existing instance id:%s", id)
	}
	newBytes, err := json.Marshal(new)
	if err != nil {
		return err
	}
	t.ops[key] = operation{opType: upsert, jsonPatch: newBytes}
	return nil
}

func (t *Txn) Save(id string, updated interface{}) error {
	key := ds.NewKey(id)
	actual, err := t.model.datastore.Get(key)
	if err == ds.ErrNotFound {
		return fmt.Errorf("can't save non-existing instance id:%s", id)
	}
	if err != nil {
		return err
	}
	newBytes, err := json.Marshal(updated)
	if err != nil {
		return err
	}
	jsonPatch, err := jsonpatch.CreateMergePatch(actual, newBytes)
	if err != nil {
		return err
	}
	t.ops[key] = operation{opType: upsert, jsonPatch: jsonPatch}
	return nil
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
	t.ops[key] = operation{opType: delete}
	return nil
}

func (t *Txn) FindByID(id string, v interface{}) error {
	return t.model.FindByID(id, v)
}
