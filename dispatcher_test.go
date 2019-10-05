package eventstore

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

type nullReducer struct{}

func (n *nullReducer) Reduce(event Event) error {
	return nil
}

type errorReducer struct{}

func (n *errorReducer) Reduce(event Event) error {
	return errors.New("error")
}

type slowReducer struct{}

func (n *slowReducer) Reduce(event Event) error {
	time.Sleep(2 * time.Second)
	return nil
}

func setup() *Dispatcher {
	return NewDispatcher(datastore.NewMapDatastore())
}

func TestNewDispatcher(t *testing.T) {
	dispatcher := setup()
	event := &nullEvent{Timestamp: time.Now()}
	dispatcher.Dispatch(event)
}

func TestRegister(t *testing.T) {
	dispatcher := setup()
	token := dispatcher.Register(&nullReducer{})
	if token != 1 {
		t.Error("callback registration failed")
	}
	if len(dispatcher.reducers) < 1 {
		t.Error("expected callbacks map to have non-zero length")
	}
}

func TestDispatchLock(t *testing.T) {
	dispatcher := setup()
	dispatcher.Register(&slowReducer{})
	event := &nullEvent{Timestamp: time.Now()}
	t1 := time.Now()
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		if err := dispatcher.Dispatch(event); err != nil {
			t.Error("unexpected error in dispatch call")
		}
	}()
	if err := dispatcher.Dispatch(event); err != nil {
		t.Error("unexpected error in dispatch call")
	}
	wg.Wait()
	t2 := time.Now()
	if t2.Sub(t1) < (4 * time.Second) {
		t.Error("reached this point too soon")
	}
}

func TestDeregister(t *testing.T) {
	dispatcher := setup()
	dispatcher.Deregister(99) // no-op
	token := dispatcher.Register(&nullReducer{})
	dispatcher.Deregister(token)
	if len(dispatcher.reducers) > 0 {
		t.Error("expected reducers map to have zero length")
	}
}

func TestDispatch(t *testing.T) {
	dispatcher := setup()
	event := &nullEvent{Timestamp: time.Now()}
	if err := dispatcher.Dispatch(event); err != nil {
		t.Error("unexpected error in dispatch call")
	}
	results, err := dispatcher.Query(query.Query{})
	if rest, _ := results.Rest(); len(rest) != 1 {
		t.Errorf("expected 1 result, got %d", len(rest))
	}
	dispatcher.Register(&errorReducer{})
	err = dispatcher.Dispatch(event)
	if errs, ok := err.(*multierror.Error); ok {
		if len(errs.Errors) != 1 {
			t.Error("should be one error")
		}
		if errs.Errors[0].Error() != "warning error" {
			t.Errorf("`%s` should be `warning error`", err)
		}
	} else {
		t.Error("expected error in dispatch call")
	}
	results, err = dispatcher.Query(query.Query{})
	if rest, _ := results.Rest(); len(rest) != 1 {
		t.Errorf("expected 1 result, got %d", len(rest))
	}
}

func TestQuery(t *testing.T) {
	dispatcher := setup()
	var events []Event
	n := 100
	for i := 1; i <= n; i++ {
		events = append(events, &nullEvent{Timestamp: time.Now()})
		time.Sleep(time.Millisecond)
	}
	for _, event := range events {
		if err := dispatcher.Dispatch(event); err != nil {
			t.Error("unexpected error in dispatch call")
		}
	}
	results, err := dispatcher.Query(query.Query{
		Orders: []query.Order{query.OrderByKey{}},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if rest, _ := results.Rest(); len(rest) != n {
		t.Errorf("expected %d result, got %d", n, len(rest))
	}
}
