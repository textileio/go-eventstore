package eventstore

import (
	"bytes"
	"encoding/binary"
	"time"
)

// Event is a generic structure for adding events to the Event Store
//@todo: Decide on what this should actually look like!
type Event interface {
	Body() []byte
	Time() []byte
	EntityID() string
	Type() string
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

func (n *nullEvent) EntityID() string {
	return "null"
}

func (n *nullEvent) Type() string {
	return "null"
}

// Sanity check
var _ Event = (*nullEvent)(nil)

type jsonPatchEvent struct {
	Timestamp time.Time
	entityID  string
	TypeName  string
	Patch     []byte
}

func NewJsonPatchEvent(t time.Time, entityID, typeName string, patch []byte) Event {
	return jsonPatchEvent{
		Timestamp: t,
		entityID:  entityID,
		TypeName:  typeName,
		Patch:     patch,
	}
}

func (je jsonPatchEvent) Body() []byte {
	return je.Patch
}

func (je jsonPatchEvent) Time() []byte {
	t := je.Timestamp.UnixNano()
	buf := new(bytes.Buffer)
	// Use big endian to preserve lexicographic sorting
	binary.Write(buf, binary.BigEndian, t)
	return buf.Bytes()
}

func (je jsonPatchEvent) EntityID() string {
	return je.entityID
}

func (je jsonPatchEvent) Type() string {
	return je.TypeName
}
