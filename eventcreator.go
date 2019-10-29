package eventstore

import (
	ds "github.com/ipfs/go-datastore"
)

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
	// Before is the instance before operation was done
	Previous interface{}
	// After is the instance after
	Current interface{}
}

type EventCreator interface {
	Reduce(e Event, datastore ds.Datastore, baseKey ds.Key) error
	// Create corresponding events to be dispatched
	Create(ops []Action) ([]Event, error)
}
