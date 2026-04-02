package ingest

import (
	"bytes"
	"context"
	"fmt"
	"testing"
)

func TestDiscover_LocalDir(t *testing.T) {
	// Points to your actual directory
	targetDir := "./testDir"
	
	// We'll use a worker pool of 3
	processor := NewParallel(2)
	ctx := context.Background()

	// Define the search: Look for the 'beta' object
	comparator := func(b []byte) bool {
		return bytes.Contains(b, []byte(`"name": "zeta"`))
	}

	// Run the discovery
	data, err := processor.Discover(ctx, targetDir, comparator)

	// Assertions
	if err != nil {
		t.Fatalf("Expected to discover a file, but got error: %v", err)
	}

	if data == nil {
		t.Fatal("Expected data to be returned, but got nil")
	}

	// Verify we got the right one
	if !bytes.Contains(data, []byte(`"id": 100`)) {
		t.Errorf("Discovered the wrong file! Content: %s", string(data))
	}

	fmt.Println("Successfully discovered and read. Content:", string(data))
}

