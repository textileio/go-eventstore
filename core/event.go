package core

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/google/uuid"
	ds "github.com/ipfs/go-datastore"
)

const (
	EmptyEntityID = EntityID("")
)

type EntityID string

func NewEntityID() EntityID {
	return EntityID(uuid.New().String())
}

func (e EntityID) String() string {
	return string(e)
}

type Event interface {
	Body() []byte
	Time() []byte
	EntityID() EntityID
	Type() string
}

func NewNullEvent(t time.Time) Event {
	return &nullEvent{Timestamp: t}
}

type nullEvent struct {
	Timestamp time.Time
}

func (n *nullEvent) Body() []byte {
	return nil
}

func (n *nullEvent) Time() []byte {
	t := n.Timestamp.UnixNano()
	buf := new(bytes.Buffer)
	// Use big endian to preserve lexicographic sorting
	binary.Write(buf, binary.BigEndian, t)
	return buf.Bytes()
}

func (n *nullEvent) EntityID() EntityID {
	return "null"
}

func (n *nullEvent) Type() string {
	return "null"
}

// Sanity check
var _ Event = (*nullEvent)(nil)

type ActionType int

const (
	Create ActionType = iota
	Save
	Delete
)

// Action is a operation done in the model
type Action struct {
	// Type of the action
	Type ActionType
	// EntityID of the instance in action
	EntityID EntityID
	// EntityType of the instance in action
	EntityType string
	// Previous is the instance before the action
	Previous interface{}
	// Current is the instance after the action was done
	Current interface{}
}

type EventCodec interface {
	// Reduce applies generated events into state
	Reduce(e Event, datastore ds.Datastore, baseKey ds.Key) error
	// Create corresponding events to be dispatched
	Create(ops []Action) ([]Event, error)
}
