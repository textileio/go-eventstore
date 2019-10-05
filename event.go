package eventstore

// Event is a generic structure for adding events to the Event Store
//@todo: Decide on what this should actually look like!
type Event interface {
	Body() []byte
	Time() []byte
	EntityID() string
	Type() string
}
