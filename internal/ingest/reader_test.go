package ingest

import (
	"context"
	"encoding/json"
	"slices"
	"testing"
)



func TestDiscover_LocalDir(t *testing.T) {
	// Setup processor
	processor := NewParallel(2)
	ctx := context.Background()
	targetDir := "./testDir" // Ensure this directory exists with test files

	tests := []struct {
		name       string
		comparator Comparator
		expectData bool
	}{
		{
			name:       "Range Match: Age between 100 and 200",
			comparator: RangeBuilder("Age", 18, 2000),
			expectData: true,
		},
		{
			name:       "Membership Match: Status is active",
			comparator: Membershipbuilder("value", []any{"active", "pending"}),
			expectData: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := processor.Discover(ctx, targetDir, tt.comparator)

			if err != nil {
				t.Fatalf("Operation failed: %v", err)
			}

			if tt.expectData && data == nil {
				t.Error("Expected to find data, but got nil")
			}

			if !tt.expectData && data != nil {
				t.Errorf("Expected no data to be found, but got: %s", string(data))
			}

			t.Log(string(data))
		})
	}
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