package ingest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

)

type FileProcessor struct {
	MaxWorkers int
}

func NewParallel(maxWorkers int) *FileProcessor {
	if maxWorkers <= 0 {
		maxWorkers = 5
	}
	return &FileProcessor{MaxWorkers: maxWorkers}
}

func (f *FileProcessor)Discover(ctx context.Context, path string, cmp Comparator) ([]byte, error) {
	entries,err:=os.ReadDir(path); if err!=nil {
		return nil, fmt.Errorf("unable to files at path %s : %w", path, err)
	}

	var files []string

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name())==".json" {
			files = append(files, filepath.Join(path, entry.Name()))
		}	
	}

	searchCtx,cancel:=context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	semaphore:= make(chan struct{}, f.MaxWorkers)
	resultChan := make(chan []byte, 1)

	for _, p:=range files{
		wg.Add(1)
		go func(p string)  {
			defer wg.Done()
			select{
				case <-searchCtx.Done():
					return
				case semaphore<- struct{}{}:
					defer func() { <-semaphore }()
			}
			content,err:=os.ReadFile(p); if err!=nil {
				return
			}
			if cmp(content){
				select{
					case resultChan<- content:
						cancel()
					default:
				}
			}
		}(p)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

		select {
	case data, ok := <-resultChan:
		if ok {
			return data, nil
		}
		return nil, ErrNotFound
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}