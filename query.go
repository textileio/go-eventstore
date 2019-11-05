package eventstore

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	dsquery "github.com/ipfs/go-datastore/query"
)

var (
	panicInvalidSortingField = "sorting field doesn't correspond to instance type"
	panicCantCompareOnSort   = "can't compare while sorting"
)

type Query struct {
	ands []*Criterion
	ors  []*Query
	sort struct {
		field string
		desc  bool
	}
}

func Where(field string) *Criterion {
	return &Criterion{
		fieldPath: field,
	}
}

func (q *Query) And(field string) *Criterion {
	return &Criterion{
		fieldPath: field,
		query:     q,
	}
}

func (q *Query) Or(orQuery *Query) *Query {
	q.ors = append(q.ors, orQuery)
	return q
}

func (q *Query) OrderBy(field string) *Query {
	q.sort.field = field
	q.sort.desc = false
	return q
}

func (q *Query) OrderByDesc(field string) *Query {
	q.sort.field = field
	q.sort.desc = true
	return q
}

func (q *Query) match(v reflect.Value) (bool, error) {
	if q == nil {
		panic("query can't be nil")
	}

	andOk := true
	for _, criterion := range q.ands {
		fieldForMatch, err := traverseFieldPath(v, criterion.fieldPath)
		if err != nil {
			return false, err
		}
		ok, err := criterion.match(fieldForMatch)
		if err != nil {
			return false, err
		}
		andOk = andOk && ok
		if !andOk {
			break
		}
	}
	if andOk {
		return true, nil
	}

	for _, q := range q.ors {
		ok, err := q.match(v)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (t *Txn) Find(res interface{}, q *Query) error {
	// ToDo: context cancellation? (to call dsr.Close())
	valRes := reflect.ValueOf(res)
	if valRes.Kind() != reflect.Ptr || valRes.Elem().Kind() != reflect.Slice {
		panic("result should be a slice")
	}
	if q == nil {
		q = &Query{}
	}
	dsq := dsquery.Query{
		Prefix: t.model.dsKey.String(),
	}
	dsr, err := t.model.datastore.Query(dsq)
	if err != nil {
		return fmt.Errorf("error when internal query: %v", err)
	}

	resSlice := valRes.Elem()
	resSlice.Set(resSlice.Slice(0, 0)) // ToDo: Document that received slice is niled
	// ToDo: also check `res` is slice of *model type*
	var unsorted []reflect.Value
	for {
		res, ok := dsr.NextSync()
		if !ok {
			break
		}
		instance := reflect.New(t.model.valueType.Elem())
		if err = json.Unmarshal(res.Value, instance.Interface()); err != nil {
			return fmt.Errorf("error when unmarshaling query result: %v", err)
		}
		ok, err = q.match(instance)
		if err != nil {
			return fmt.Errorf("error when matching entry with query: %v", err)
		}
		if ok {
			unsorted = append(unsorted, instance)
		}
	}
	if q.sort.field != "" {
		sort.Slice(unsorted, func(i, j int) bool {
			fieldI, err := traverseFieldPath(unsorted[i], q.sort.field)
			if err != nil {
				panic(panicInvalidSortingField)
			}
			fieldJ, err := traverseFieldPath(unsorted[j], q.sort.field)
			if err != nil {
				panic(panicInvalidSortingField)
			}
			res, err := compare(fieldI, fieldJ)
			if err != nil {
				panic(panicCantCompareOnSort)
			}
			if q.sort.desc {
				res *= -1
			}
			return res < 0
		})
	}

	for i := range unsorted {
		resSlice = reflect.Append(resSlice, unsorted[i])
	}
	valRes.Elem().Set(resSlice)
	return nil
}
