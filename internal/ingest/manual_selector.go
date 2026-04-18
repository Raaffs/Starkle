package ingest

import (
	"context"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ManualSelectionProcessor struct{}

func NewSelection() *ManualSelectionProcessor {
	return &ManualSelectionProcessor{}
}

func (s *ManualSelectionProcessor)Discover(ctx context.Context, path string, cmp Comparator) ([]byte, error) {
	filePath,err:=runtime.OpenFileDialog(ctx,runtime.OpenDialogOptions{
		Title : "Select Proof",
		Filters: []runtime.FileFilter{{
				DisplayName: "Proof files (.json)",
				Pattern: "*.json",
			},
		},
	})
	if err!=nil {
		return nil, fmt.Errorf("unable to open runtime dialog")
	}
	
	file,err:=os.ReadFile(filePath); if err!=nil {
		return nil, fmt.Errorf("unable to read file %s : %w", filePath, err)
	}
	if !cmp(file){
		return nil, ErrNotFound
	}
	return file,nil
}