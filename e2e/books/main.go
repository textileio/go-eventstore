package main

import (
	ds "github.com/ipfs/go-datastore"
	es "github.com/textileio/go-eventstore"
	"github.com/textileio/go-eventstore/store"
)

type Book struct {
	ID     es.EntityID
	Title  string
	Author string
	Meta   BookStats
}

type BookStats struct {
	TotalReads int
	Rating     float64
}

func main() {
	s := createMemStore()

	model, err := s.RegisterJSONPatcher("Book", &Book{})
	checkErr(err)

	// Bootstrap the model with some books: two from Author1 and one from Author2
	{
		// Create a two books for Author1
		book1 := &Book{ // Notice ID will be autogenerated
			Title:  "Title1",
			Author: "Author1",
			Meta:   BookStats{TotalReads: 100, Rating: 3.2},
		}
		book2 := &Book{
			Title:  "Title2",
			Author: "Author1",
			Meta:   BookStats{TotalReads: 150, Rating: 4.1},
		}
		checkErr(model.Create(book1, book2)) // Note you can create multiple books at the same time (variadic)

		// Create book for Author2
		book3 := &Book{
			Title:  "Title3",
			Author: "Author2",
			Meta:   BookStats{TotalReads: 500, Rating: 4.9},
		}
		checkErr(model.Create(book3))
	}

	// Query all the books
	{
		var books []*Book
		err := model.Find(&books, &store.Query{})
		checkErr(err)
		if len(books) != 3 {
			panic("there should be three books")
		}
	}

	// Query the books from Author2
	{
		var books []*Book
		err := model.Find(&books, store.Where("Author").Eq("Author1"))
		checkErr(err)
		if len(books) != 2 {
			panic("Author1 should have two books")
		}
	}

	// Query book by two conditions
	{
		var books []*Book
		err := model.Find(&books, store.Where("Author").Eq("Author1").And("Title").Eq("Title2"))
		checkErr(err)
		if len(books) != 1 {
			panic("Author1 should have only one book with Title2")
		}
	}

	// Query book by OR condition
	{
		var books []*Book
		err := model.Find(&books, store.Where("Author").Eq("Author1").Or(store.Where("Author").Eq("Author2")))
		checkErr(err)
		if len(books) != 3 {
			panic("Author1 & Author2 have should have 3 books in total")
		}
	}

	// // Sorted query
	// {
	// 	var books []*Book
	// 	err := model.Find(&books, store.Where("Title"))
	// 	checkErr(err)

	// 	// Modify title
	// 	book := books[0]
	// 	book.Title = "ModifiedTitle"
	// 	model.Save(book)
	// 	err = model.Find(&books, store.Where("Title").Eq("Title3"))
	// 	checkErr(err)
	// 	if len(books) != 0 {
	// 		panic("Book with Title3 shouldn't exist")
	// 	}

	// 	// Delete it
	// 	err = model.Find(&books, store.Where("Title").Eq("ModifiedTitle"))
	// 	checkErr(err)
	// 	if len(books) != 1 {
	// 		panic("Book with ModifiedTitle should exist")
	// 	}
	// 	model.Delete(books[0].ID)
	// 	err = model.Find(&books, store.Where("Title").Eq("ModifiedTitle"))
	// 	checkErr(err)
	// 	if len(books) != 0 {
	// 		panic("Book with ModifiedTitle shouldn't exist")
	// 	}
	// }

	// Query, Update, and Save
	{
		var books []*Book
		err := model.Find(&books, store.Where("Title").Eq("Title3"))
		checkErr(err)

		// Modify title
		book := books[0]
		book.Title = "ModifiedTitle"
		model.Save(book)
		err = model.Find(&books, store.Where("Title").Eq("Title3"))
		checkErr(err)
		if len(books) != 0 {
			panic("Book with Title3 shouldn't exist")
		}

		// Delete it
		err = model.Find(&books, store.Where("Title").Eq("ModifiedTitle"))
		checkErr(err)
		if len(books) != 1 {
			panic("Book with ModifiedTitle should exist")
		}
		model.Delete(books[0].ID)
		err = model.Find(&books, store.Where("Title").Eq("ModifiedTitle"))
		checkErr(err)
		if len(books) != 0 {
			panic("Book with ModifiedTitle shouldn't exist")
		}
	}

	// ToDo: Create indexes
	// ToDo: Use indexes for queries
	// ToDo: Sorting?
	// ToDo: Self-referencing conditionals
}

func createMemStore() *store.Store {
	datastore := ds.NewMapDatastore()
	dispatcher := es.NewDispatcher(es.NewTxMapDatastore())
	return store.NewStore(datastore, dispatcher)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
