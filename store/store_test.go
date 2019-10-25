package store

import (
	"os"
	"reflect"
	"testing"

	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	es "github.com/textileio/go-eventstore"
)

const (
	errInvalidInstanceState = "invalid instance state"
)

type Person struct {
	ID   es.EntityID
	Name string
	Age  int
}

type Dog struct {
	ID       es.EntityID
	Name     string
	Comments []Comment
}
type Comment struct {
	Body string
}

func TestMain(m *testing.M) {
	logging.SetLogLevel("*", "debug")
	os.Exit(m.Run())
}

func TestSchemaRegistration(t *testing.T) {
	t.Parallel()
	t.Run("Single", func(t *testing.T) {
		t.Parallel()
		store := createTestStore()
		_, err := store.RegisterJSONPatcher("Dog", &Dog{})
		checkErr(t, err)
	})
	t.Run("Multiple", func(t *testing.T) {
		t.Parallel()
		store := createTestStore()
		_, err := store.RegisterJSONPatcher("Dog", &Dog{})
		checkErr(t, err)
		_, err = store.RegisterJSONPatcher("Person", &Person{})
		checkErr(t, err)
	})
	t.Run("Fail/WithoutEntityID", func(t *testing.T) {
		t.Parallel()
		type FailingModel struct {
			IDontHaveAnIDField int
		}
		store := createTestStore()
		if _, err := store.RegisterJSONPatcher("FailingModel", &FailingModel{}); err != ErrInvalidModel {
			t.Fatal("the model should be invalid")
		}
	})
}

func TestAddInstance(t *testing.T) {
	t.Parallel()
	t.Run("Single", func(t *testing.T) {
		t.Parallel()
		store := createTestStore()
		model, err := store.RegisterJSONPatcher("Person", &Person{})
		checkErr(t, err)

		t.Run("WithImplicitTx", func(t *testing.T) {
			newPerson := &Person{Name: "Foo", Age: 42}
			err = model.Add(newPerson)
			checkErr(t, err)
			assertPersonInModel(t, model, newPerson)
		})
		t.Run("WithTx", func(t *testing.T) {
			newPerson := &Person{Name: "Foo", Age: 42}
			err = model.WriteTxn(func(txn *Txn) error {
				return txn.Create(newPerson)
			})
			checkErr(t, err)
			assertPersonInModel(t, model, newPerson)
		})
	})
	t.Run("Multiple", func(t *testing.T) {
		t.Parallel()
		store := createTestStore()
		model, err := store.RegisterJSONPatcher("Person", &Person{})
		checkErr(t, err)

		newPerson1 := &Person{Name: "Foo1", Age: 42}
		newPerson2 := &Person{Name: "Foo2", Age: 43}
		err = model.WriteTxn(func(txn *Txn) error {
			err := txn.Create(newPerson1)
			if err != nil {
				return err
			}
			return txn.Create(newPerson2)
		})
		checkErr(t, err)
		assertPersonInModel(t, model, newPerson1)
		assertPersonInModel(t, model, newPerson2)
	})

	// ToDo: Add test for `.Create` on an instance with an
	// assigned ID (shouldn't be overwritten)
}

func TestGetInstance(t *testing.T) {
	t.Parallel()

	store := createTestStore()
	model, err := store.RegisterJSONPatcher("Person", &Person{})
	checkErr(t, err)

	newPerson := &Person{Name: "Foo", Age: 42}
	err = model.WriteTxn(func(txn *Txn) error {
		return txn.Create(newPerson)
	})
	checkErr(t, err)

	t.Run("WithImplicitTx", func(t *testing.T) {
		person := &Person{}
		err = model.FindByID(newPerson.ID, person)
		checkErr(t, err)
		if !reflect.DeepEqual(newPerson, person) {
			t.Fatalf(errInvalidInstanceState)
		}
	})
	t.Run("WithReadTx", func(t *testing.T) {
		person := &Person{}
		err = model.ReadTxn(func(txn *Txn) error {
			txn.FindByID(newPerson.ID, person)
			checkErr(t, err)
			if !reflect.DeepEqual(newPerson, person) {
				t.Fatalf(errInvalidInstanceState)
			}
			return nil
		})
	})
	t.Run("WithUpdateTx", func(t *testing.T) {
		person := &Person{}
		err = model.WriteTxn(func(txn *Txn) error {
			txn.FindByID(newPerson.ID, person)
			checkErr(t, err)
			if !reflect.DeepEqual(newPerson, person) {
				t.Fatalf(errInvalidInstanceState)
			}
			return nil
		})
	})
}

func TestUpdateInstance(t *testing.T) {
	t.Parallel()

	store := createTestStore()
	model, err := store.RegisterJSONPatcher("Person", &Person{})
	checkErr(t, err)

	newPerson := &Person{Name: "Alice", Age: 42}
	err = model.WriteTxn(func(txn *Txn) error {
		return txn.Create(newPerson)
	})
	checkErr(t, err)

	err = model.WriteTxn(func(txn *Txn) error {
		p := &Person{}
		err := txn.FindByID(newPerson.ID, p)
		checkErr(t, err)

		p.Name = "Bob"
		return txn.Save(p)
	})
	checkErr(t, err)

	// Under the hood here the instance update went through
	// the dispatcher, then the reducer, which will ultimately
	// apply the change to the current instance state that
	// should make the code below behave as expected

	person := &Person{}
	err = model.FindByID(newPerson.ID, person)
	checkErr(t, err)
	if person.ID != newPerson.ID || person.Age != 42 || person.Name != "Bob" {
		t.Fatalf(errInvalidInstanceState)
	}
}

func TestDeleteInstance(t *testing.T) {
	t.Parallel()

	store := createTestStore()
	model, err := store.RegisterJSONPatcher("Person", &Person{})
	checkErr(t, err)

	newPerson := &Person{Name: "Alice", Age: 42}
	err = model.WriteTxn(func(txn *Txn) error {
		return txn.Create(newPerson)
	})
	checkErr(t, err)

	err = model.Delete(newPerson.ID)
	checkErr(t, err)

	if err = model.FindByID(newPerson.ID, &Person{}); err != ErrNotFound {
		t.Fatalf("FindByID: instance shouldn't exist")
	}
	if exist, err := model.Has(newPerson.ID); exist || err != nil {
		t.Fatalf("Has: instance shouldn't exist")
	}

	// Try to delete again
	if err = model.Delete(newPerson.ID); err != ErrNotFound {
		t.Fatalf("cant't delete non-existent instance")
	}
}

// ToDo
func TestInvalidActions(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		// Compared to schema
	})
	t.Run("Update", func(t *testing.T) {
		// Compared to schema
	})
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertPersonInModel(t *testing.T, model *Model, person *Person) {
	t.Helper()
	p := &Person{}
	err := model.FindByID(person.ID, p)
	checkErr(t, err)
	if !reflect.DeepEqual(person, p) {
		t.Fatalf(errInvalidInstanceState)
	}
}

func createTestStore() *Store {
	datastore := ds.NewMapDatastore()
	dispatcher := es.NewDispatcher(es.NewTxMapDatastore())
	return NewStore(datastore, dispatcher)
}
