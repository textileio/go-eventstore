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
