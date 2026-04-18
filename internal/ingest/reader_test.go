package ingest

import (
	"context"
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
			name:       "Range Match: ID between 100 and 200",
			comparator: RangeBuilder("id", 100, 200),
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