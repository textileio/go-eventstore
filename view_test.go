package eventstore

import (
	"testing"
	"time"

	datastore "github.com/ipfs/go-datastore"
	"github.com/textileio/go-eventstore/broadcast"
)

// MapModel implements StoredModel using a map.
func NewMapModel() *MapModel {
	return &MapModel{
		StoredModel: &StoredModel{
			MemoryModel: &MemoryModel{
				broadcaster: &broadcast.Broadcaster{},
			},
			store: datastore.NewMapDatastore(),
		},
	}
}

type MapModel struct {
	*StoredModel
}

// Reduce puts the `Body` of the incoming event at the `EntityID` key in the map store.
func (m *MapModel) Reduce(event Event) error {
	err := m.Store().Put(datastore.NewKey(event.EntityID()), event.Body())
	m.broadcaster.Send(event)
	return err
}

// Sanity check
var _ ViewModel = (*MapModel)(nil)

// BodyModel implements MemoryModel with a single Body property.
type BodyModel struct {
	MemoryModel
	Body []byte
}

// Reduce replaces `Body` with the event's `Body`
func (m *BodyModel) Reduce(event Event) error {
	m.Body = event.Body()
	return nil
}

// Sanity check
var _ ViewModel = (*BodyModel)(nil)

func TestNewMemoryModel(t *testing.T) {
	viewmodel := &BodyModel{}
	event := &nullEvent{Timestamp: time.Now()}
	if err := viewmodel.Reduce(event); err != nil {
		t.Errorf("error calling reduce: %s", err.Error())
	}
}

func TestReduceNotImplemented(t *testing.T) {
	viewmodel := &MemoryModel{}
	event := &nullEvent{Timestamp: time.Now()}
	if err := viewmodel.Reduce(event); err == nil {
		t.Error("should return error")
	}
}

func TestNewStoredModel(t *testing.T) {
	viewmodel := NewMapModel()
	event := &nullEvent{Timestamp: time.Now()}
	if err := viewmodel.Reduce(event); err != nil {
		t.Errorf("error calling reduce: %s", err.Error())
	}
	if viewmodel.Store() == nil {
		t.Error("nil store encountered")
	}
}

// @todo: More tests!
