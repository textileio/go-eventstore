package store

import (
	"fmt"
	"reflect"
)

type Criterion struct {
	fieldPath string
	operation operation
	value     interface{}
	query     *Query
}

func (c *Criterion) Eq(value interface{}) *Query {
	c.operation = eq
	c.value = value

	// First Criterion of a query?
	if c.query == nil {
		c.query = &Query{}
	}

	c.query.ands = append(c.query.ands, c)

	return c.query
}

func (c *Criterion) compare(rowValue, criterionValue interface{}) (int, error) {
	if rowValue == nil || criterionValue == nil {
		if rowValue == criterionValue {
			return 0, nil
		}
		return 0, &ErrTypeMismatch{rowValue, criterionValue}
	}

	value := rowValue

	for reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	other := criterionValue
	for reflect.TypeOf(other).Kind() == reflect.Ptr {
		other = reflect.ValueOf(other).Elem().Interface()
	}

	return compare(value, other)
}

func (c *Criterion) match(testValue interface{}) (bool, error) {
	result, err := c.compare(testValue, c.value)
	if err != nil {
		return false, err
	}
	switch c.operation {
	case eq:
		return result == 0, nil
	default:
		panic("invalid operator")
	}
}

type ErrTypeMismatch struct { // ToDo: unexport?
	Value interface{}
	Other interface{}
}

func (e *ErrTypeMismatch) Error() string {
	return fmt.Sprintf("%v (%T) cannot be compared with %v (%T)", e.Value, e.Value, e.Other, e.Other)
}
