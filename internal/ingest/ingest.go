package ingest

import (
	"context"
	"encoding/json"
	"errors"
	"slices"

)

var ErrNotFound = errors.New("no matching file discovered")

type Comparator func([]byte) bool

type Finder interface {
	Discover(ctx context.Context, paths string, cmp Comparator) ([]byte, error)
}


func Membershipbuilder(attribute string, targets []any)Comparator{
	return func(b []byte) bool {
		return memerbshipComparator(b,attribute,targets)
	}
}

func RangeBuilder(attribute string, low, high int)Comparator{
	return func(b []byte) bool {
		return rangeComparator(b,attribute,low,high)
	}
}

var memerbshipComparator = func (rawJson []byte, field string, targets []any) bool {
	var data map[string]any
	if err := json.Unmarshal(rawJson, &data); err != nil {
		return false
	}
	val, exists := data[field]
	if !exists {
		return false
	}

	return slices.Contains(targets, val)
}

var rangeComparator = func (rawJson []byte, field string, lower int, upper int) bool {
	var data map[string]any
	if err := json.Unmarshal(rawJson, &data); err != nil {
		return false
	}
	val, exists := data[field]
	if !exists {
		return false
	}
	// JSON numbers unmarshal as float64. We convert to float for the comparison
	// or cast to int if we're sure it's a whole number.
	num, ok := val.(float64)
	if !ok {
		return false
	}

	intVal := int(num)

	return intVal >= lower && intVal <= upper
}
