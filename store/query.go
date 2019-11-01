package store

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	dsquery "github.com/ipfs/go-datastore/query"
)

type operation int

const (
	eq operation = iota
	// ToDo: Further operations here
)

type Query struct {
	ands []*Criterion
	ors  []*Query
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
	resSlice := valRes.Elem()
	resSlice.Set(resSlice.Slice(0, 0)) // ToDo: Document that received slice is niled

	// ToDo: also check `result` is slice of *model type*

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
	for {
		res, ok := dsr.NextSync()
		if !ok {
			break
		}

		instance := reflect.New(t.model.valueType.Elem())
		err = json.Unmarshal(res.Value, instance.Interface())
		if err != nil {
			return fmt.Errorf("error when unmarhsaling query result: %v", err)
		}
		ok, err := q.match(instance)
		if err != nil {
			return fmt.Errorf("error when matching entry with query: %v", err)
		}
		if ok {
			resSlice = reflect.Append(resSlice, instance)
		}
	}
	valRes.Elem().Set(resSlice.Slice(0, resSlice.Len()))
	return nil
}

func traverseFieldPath(value reflect.Value, fieldPath string) (reflect.Value, error) {
	fields := strings.Split(fieldPath, ".")

	current := value // ToDo: Can `current` be deleted?
	for i := range fields {
		if current.Kind() == reflect.Ptr {
			current = current.Elem()
		}
		current = current.FieldByName(fields[i])

		if !current.IsValid() {
			return reflect.Value{}, fmt.Errorf("instance field %s doesn't exist in type %s", fieldPath, value)
		}
	}
	return current, nil
}
