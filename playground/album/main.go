package main

import (
	"fmt"

	"github.com/google/uuid"
	ds "github.com/ipfs/go-datastore"
	"github.com/textileio/go-eventstore/playground"
)

type Comment struct {
	Body string
}

var dogType = &Dog{}

type Dog struct {
	ID       string
	Name     string
	Age      int
	Comments []Comment
}

func main() {
	// Now used MapDatstore which is in-mem, but can be switched
	store := playground.NewStore(ds.NewMapDatastore())

	// Register model
	dogModel, err := store.Register("Dog", dogType)
	checkErr(err)

	// Create new model instance
	// The `model.Update(..)` function gives guarantees that nobody is messing around during internal func.
	// Saying it differently, gives isolation guarantees(serialization) during update logic. Also on commit, an 'all or nothing'
	// update is done (e.g: two instances of the model could be updated atomically event-wise, not sure will be used though so can be simpler)
	id := uuid.New().String()
	err = dogModel.Update(func(txn *playground.Txn) error {
		dog := Dog{
			ID:   id,
			Name: "Bob",
			Age:  2,
			Comments: []Comment{
				{
					Body: "Nice dog",
				},
			},
		}
		fmt.Printf("Original Dog: %#v\n\n", dog)
		return txn.Add(dog.ID, dog)
	})
	checkErr(err)

	// under the hood `txn` is attached to `dogModel`, so `FindByID` and `Save` already
	// assumes we're talking about `Dog` model.
	// Notes:
	// - `Save()` could validate schema to enforce validation
	err = dogModel.Update(func(txn *playground.Txn) error {
		dog := &Dog{}

		err := txn.FindByID(id, dog)
		checkErr(err)

		dog.Name = "Norton"
		dog.Age = 10
		return txn.Save(dog.ID, dog)
	})
	checkErr(err)

	dog := Dog{}
	checkErr(dogModel.FindByID(id, &dog))
	fmt.Printf("Final FindByID(): %#v", dog)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
