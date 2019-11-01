package jsonpatcher

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	es "github.com/textileio/go-eventstore"
)

type operationType int

const (
	create operationType = iota
	save
	delete
)

var (
	log                           = logging.Logger("jsonpatcher")
	errSavingNonExistentInstance  = errors.New("can't save nonexistent instance")
	errCantCreateExistingInstance = errors.New("cant't create already existent instance")
	errUnknownOperation           = errors.New("unknown operation type")
)

type operation struct {
	Type      operationType
	EntityID  es.EntityID
	JSONPatch []byte
}

type patcher struct {
}

var _ es.EventCodec = (*patcher)(nil)

func New() es.EventCodec {
	return &patcher{}
}

func (m *patcher) Create(actions []es.Action) ([]es.Event, error) {
	events := make([]es.Event, len(actions))
	for i := range actions {
		var eventPayload []byte
		var err error
		switch actions[i].Type {
		case es.Create:
			eventPayload, err = createEvent(actions[i].EntityID, actions[i].Current)
		case es.Save:
			eventPayload, err = saveEvent(actions[i].EntityID, actions[i].Previous, actions[i].Current)
		case es.Delete:
			eventPayload, err = deleteEvent(actions[i].EntityID)
		default:
			panic("unkown action type")
		}
		if err != nil {
			return nil, err
		}
		events[i] = patchEvent{
			Timestamp: time.Now(),
			ID:        actions[i].EntityID,
			TypeName:  actions[i].EntityType,
			Patch:     eventPayload,
		}
	}
	return events, nil
}

func (p *patcher) Reduce(e es.Event, datastore ds.Datastore, baseKey ds.Key) error {
	var op operation
	if err := json.Unmarshal(e.Body(), &op); err != nil {
		return err
	}

	key := baseKey.ChildString(e.EntityID().String())
	switch op.Type {
	case create:
		exist, err := datastore.Has(key)
		if err != nil {
			return err
		}
		if exist {
			return errCantCreateExistingInstance
		}
		if err := datastore.Put(key, op.JSONPatch); err != nil {
			return fmt.Errorf("error when reducing create event: %v", err)
		}
		log.Debug("\tcreate operation applied")
	case save:
		value, err := datastore.Get(key)
		if errors.Is(err, ds.ErrNotFound) {
			return errSavingNonExistentInstance
		}
		if err != nil {
			return err
		}
		patchedValue, err := jsonpatch.MergePatch(value, op.JSONPatch)
		if err != nil {
			return fmt.Errorf("error when reducing save event: %v", err)
		}
		if err = datastore.Put(key, patchedValue); err != nil {
			return err
		}
		log.Debug("\tsave operation applied")
	case delete:
		if err := datastore.Delete(key); err != nil {
			return err
		}
		log.Debug("\tdelete operation applied")
	default:
		return errUnknownOperation
	}

	return nil
}

func createEvent(id es.EntityID, v interface{}) ([]byte, error) {
	opBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	op := operation{
		Type:      create,
		EntityID:  id,
		JSONPatch: opBytes,
	}
	eventPayload, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	return eventPayload, nil
}

func saveEvent(id es.EntityID, prev interface{}, curr interface{}) ([]byte, error) {
	prevBytes, err := json.Marshal(prev)
	if err != nil {
		return nil, err
	}
	currBytes, err := json.Marshal(curr)
	if err != nil {
		return nil, err
	}
	jsonPatch, err := jsonpatch.CreateMergePatch(prevBytes, currBytes)
	if err != nil {
		return nil, err
	}
	op := operation{
		Type:      save,
		EntityID:  id,
		JSONPatch: jsonPatch,
	}
	eventPayload, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	return eventPayload, nil
}

func deleteEvent(id es.EntityID) ([]byte, error) {
	op := operation{
		Type:      delete,
		EntityID:  id,
		JSONPatch: nil,
	}
	eventPayload, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	return eventPayload, nil
}

type patchEvent struct {
	Timestamp time.Time
	ID        es.EntityID
	TypeName  string
	Patch     []byte
}

func (je patchEvent) Body() []byte {
	return je.Patch
}

func (je patchEvent) Time() []byte {
	t := je.Timestamp.UnixNano()
	buf := new(bytes.Buffer)
	// Use big endian to preserve lexicographic sorting
	binary.Write(buf, binary.BigEndian, t)
	return buf.Bytes()
}

func (je patchEvent) EntityID() es.EntityID {
	return je.ID
}

func (je patchEvent) Type() string {
	return je.TypeName
}

var _ es.Event = (*patchEvent)(nil)
