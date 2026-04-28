package ingest

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("no matching file discovered")

type Comparator func([]byte) bool

type Finder interface {
	Discover(ctx context.Context, paths string, cmp Comparator) ([]byte, error)
}


