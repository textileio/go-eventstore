package store

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"github.com/alecthomas/jsonschema"
	ds "github.com/ipfs/go-datastore"
	es "github.com/textileio/go-eventstore"
)

var (
	baseKey = ds.NewKey("/model/")

	ErrNotFound   = errors.New("instance not found")
	ErrReadonlyTx = errors.New("read only transaction")

	errAlreadyDiscardedCommitedTxn = errors.New("can't commit discarded/commited txn")
	errCantCreateExistingInstance  = errors.New("can't create already existing instance")
	errCantSaveNonExistentInstance = errors.New("can't save unkown instance")
)

type Model struct {
	lock         sync.RWMutex
	schema       *jsonschema.Schema
	valueType    reflect.Type
	datastore    ds.Datastore
	eventcreator es.EventCreator
	dispatcher   *es.Dispatcher
	dsKey        ds.Key
}

func NewModel(name string, defaultInstance interface{}, datastore ds.Datastore, dispatcher *es.Dispatcher, eventcreator es.EventCreator) *Model {
	m := &Model{
		schema:       jsonschema.Reflect(defaultInstance),
		datastore:    datastore,
		valueType:    reflect.TypeOf(defaultInstance),
		dispatcher:   dispatcher,
		eventcreator: eventcreator,
		dsKey:        baseKey.ChildString(name),
	}

	return m
}

func (m *Model) ReadTxn(f func(txn *Txn) error) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	txn := &Txn{model: m, readonly: true}
	defer txn.Discard()
	if err := f(txn); err != nil {
		return err
	}
	return nil
}

func (m *Model) WriteTxn(f func(txn *Txn) error) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	txn := &Txn{model: m}
	defer txn.Discard()
	if err := f(txn); err != nil {
		return err
	}
	return txn.Commit()
}

func (m *Model) FindByID(id es.EntityID, v interface{}) error {
	return m.ReadTxn(func(txn *Txn) error {
		return txn.FindByID(id, v)
	})
}

func (m *Model) Add(v interface{}) error {
	return m.WriteTxn(func(txn *Txn) error {
		return txn.Create(v)
	})
}

func (m *Model) Delete(id es.EntityID) error {
	return m.WriteTxn(func(txn *Txn) error {
		return txn.Delete(id)
	})
}

func (m *Model) Save(v interface{}) error {
	return m.WriteTxn(func(txn *Txn) error {
		return txn.Save(v)
	})
}

func (m *Model) Has(id es.EntityID) (exists bool, err error) {
	m.ReadTxn(func(txn *Txn) error {
		exists, err = txn.Has(id)
		return err
	})
	return
}

type Txn struct {
	model     *Model
	discarded bool
	commited  bool
	readonly  bool

	actions []es.Action
}

type SaveOp struct {
	Before interface{}
	After  interface{}
}

func (t *Txn) Create(new interface{}) error {
	if t.readonly {
		return ErrReadonlyTx
	}
	id := getEntityID(new)
	if id == es.EmptyEntityID {
		id = setNewEntityID(new)
	}
	key := t.model.dsKey.ChildString(id.String())
	exists, err := t.model.datastore.Has(key)
	if err != nil {
		return err
	}
	if exists {
		return errCantCreateExistingInstance
	}

	a := es.Action{
		Type:       es.Create,
		EntityID:   id,
		EntityType: t.model.schema.Ref,
		Previous:   nil,
		Current:    new,
	}
	t.actions = append(t.actions, a)

	return nil
}

func (t *Txn) Save(updated interface{}) error {
	if t.readonly {
		return ErrReadonlyTx
	}

	id := getEntityID(updated)
	key := t.model.dsKey.ChildString(id.String())
	beforeBytes, err := t.model.datastore.Get(key)
	if err == ds.ErrNotFound {
		return errCantSaveNonExistentInstance
	}
	if err != nil {
		return err
	}

	before := reflect.New(t.model.valueType.Elem()).Interface()
	err = json.Unmarshal(beforeBytes, before)
	if err != nil {
		return err
	}
	a := es.Action{
		Type:       es.Save,
		EntityID:   id,
		EntityType: t.model.schema.Ref,
		Previous:   before,
		Current:    updated,
	}
	t.actions = append(t.actions, a)

	return nil
}

func (t *Txn) Delete(id es.EntityID) error {
	if t.readonly {
		return ErrReadonlyTx
	}
	key := t.model.dsKey.ChildString(id.String())
	exists, err := t.model.datastore.Has(key)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	a := es.Action{
		Type:       es.Delete,
		EntityID:   id,
		EntityType: t.model.schema.Ref,
		Previous:   nil,
		Current:    nil,
	}
	t.actions = append(t.actions, a)
	return nil
}

func (t *Txn) Has(id es.EntityID) (bool, error) {
	key := t.model.dsKey.ChildString(id.String())
	exists, err := t.model.datastore.Has(key)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (t *Txn) FindByID(id es.EntityID, v interface{}) error {
	key := t.model.dsKey.ChildString(id.String())
	bytes, err := t.model.datastore.Get(key)
	if errors.Is(err, ds.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func (t *Txn) Commit() error {
	if t.discarded || t.commited {
		return errAlreadyDiscardedCommitedTxn
	}

	events, err := t.model.eventcreator.Create(t.actions)
	if err != nil {
		return err
	}

	for _, e := range events {
		if err := t.model.dispatcher.Dispatch(e); err != nil {
			return err // ToDo/Note: Important to document the implications of a partial dispatch
		}
	}
	return nil
}

func (m *Model) Reduce(event es.Event) error {
	// ToDo: distinguish local and remote events for proper locking
	log.Debugf("reducer %s start", m.schema.Ref)
	if event.Type() != m.schema.Ref {
		log.Debugf("ignoring event from uninteresting type")
		return nil
	}

	return m.eventcreator.Reduce(event, m.datastore, m.dsKey)
}

func (t *Txn) Discard() {
	t.discarded = true
}

func getEntityID(t interface{}) es.EntityID {
	v := reflect.ValueOf(t)
	if v.Type().Kind() != reflect.Ptr {
		v = reflect.New(reflect.TypeOf(v))
	}
	v = v.Elem().FieldByName(idFieldName)
	if !v.IsValid() || v.Type() != reflect.TypeOf(es.EntityID("")) {
		panic("invalid instance: doesn't have EntityID attribute")
	}
	return es.EntityID(v.String())
}

func setNewEntityID(t interface{}) es.EntityID {
	v := reflect.ValueOf(t)
	if v.Type().Kind() != reflect.Ptr {
		v = reflect.New(reflect.TypeOf(v))
	}
	newID := es.NewEntityID()
	v.Elem().FieldByName(idFieldName).Set(reflect.ValueOf(newID))
	return newID
}
