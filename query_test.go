package eventstore

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/textileio/go-eventstore/core"
)

type book struct {
	ID     core.EntityID
	Title  string
	Author string
	Meta   bookStats
}

type bookStats struct {
	TotalReads int
	Rating     float64
}

type queryTest struct {
	name    string
	query   *Query
	resIdx  []int // expected idx results from sample data
	ordered bool
}

var (
	sampleData = []book{
		book{Title: "Title1", Author: "Author1", Meta: bookStats{TotalReads: 10, Rating: 3.3}},
		book{Title: "Title2", Author: "Author1", Meta: bookStats{TotalReads: 20, Rating: 3.6}},
		book{Title: "Title3", Author: "Author1", Meta: bookStats{TotalReads: 30, Rating: 3.9}},
		book{Title: "Title4", Author: "Author2", Meta: bookStats{TotalReads: 114, Rating: 4.0}},
		book{Title: "Title5", Author: "Author3", Meta: bookStats{TotalReads: 500, Rating: 4.8}},
	}

	queries = []queryTest{
		queryTest{name: "AllNil", query: nil, resIdx: []int{0, 1, 2, 3, 4}},
		queryTest{name: "AllExplicit", query: &Query{}, resIdx: []int{0, 1, 2, 3, 4}},

		queryTest{name: "FromAuthor1", query: Where("Author").Eq("Author1"), resIdx: []int{0, 1, 2}},
		queryTest{name: "FromAuthor2", query: Where("Author").Eq("Author2"), resIdx: []int{3}},
		queryTest{name: "FromAuthor3", query: Where("Author").Eq("Author3"), resIdx: []int{4}},

		queryTest{name: "FromAuthor1Asc", query: Where("Author").Eq("Author1").OrderBy("Title"), resIdx: []int{0, 1, 2}, ordered: true},
		queryTest{name: "FromAuthor1Desc", query: Where("Author").Eq("Author1").OrderByDesc("Title"), resIdx: []int{2, 1, 0}, ordered: true},
		queryTest{name: "AllDesc", query: (&Query{}).OrderByDesc("Title"), resIdx: []int{4, 3, 2, 1, 0}, ordered: true},

		queryTest{name: "AndAuthor1Title2", query: Where("Author").Eq("Author1").And("Title").Eq("Title2"), resIdx: []int{1}},
		queryTest{name: "AndAuthorNestedTotalReads", query: Where("Author").Eq("Author1").And("Meta.TotalReads").Eq(10), resIdx: []int{0}},

		queryTest{name: "OrAuthor", query: Where("Author").Eq("Author1").Or(Where("Author").Eq("Author3")), resIdx: []int{0, 1, 2, 4}},
	}
)

func TestModelQuery(t *testing.T) {
	m := createModelWithData(t)
	for _, q := range queries {
		q := q
		t.Run(q.name, func(t *testing.T) {
			t.Parallel()
			var res []*book
			if err := m.Find(&res, q.query); err != nil {
				t.Fatal("error when executing query")
			}
			if len(q.resIdx) != len(res) {
				t.Fatalf("query results length doesn't match, expected: %d, got: %d", len(q.resIdx), len(res))
			}

			expectedIdx := make([]int, len(q.resIdx))
			for i := range q.resIdx {
				expectedIdx[i] = q.resIdx[i]
			}
			if !q.ordered {
				sort.Slice(res, func(i, j int) bool {
					return strings.Compare(res[i].ID.String(), res[j].ID.String()) == -1
				})
				sort.Slice(expectedIdx, func(i, j int) bool {
					return strings.Compare(sampleData[expectedIdx[i]].ID.String(), sampleData[expectedIdx[j]].ID.String()) == -1
				})
			}
			for i, idx := range expectedIdx {
				if !reflect.DeepEqual(sampleData[idx], *res[i]) {
					t.Fatalf("wrong query item result, expected: %v, got: %v", sampleData[idx], *res[i])
				}
			}
		})
	}
}

// ToDo: invalid sorting field

func createModelWithData(t *testing.T) *Model {
	store := createTestStore()
	m, err := store.Register("Book", &book{})
	checkErr(t, err)
	for i := range sampleData {
		if err = m.Create(&sampleData[i]); err != nil {
			t.Fatalf("failed to create sample data: %v", err)
		}
	}
	return m
}
