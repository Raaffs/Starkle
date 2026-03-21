package main

import (
	"embed"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"github.com/Suy56/ProofChain/internal/crypto/keyUtils"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

type SlogHandler struct {
	slog.Handler
	l io.Writer
}

func (h *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()
	fields := make(map[string]any, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	// Structure the log output
	output := map[string]any{
		"time":    r.Time.Format("2006-01-02 15:04:05"),		"level":   level,
		"message": r.Message,
		"fields":  fields,
	}

	// Marshal with indentation
	b, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintln(h.l, string(b))
	return nil
}

func NewLogger(out io.Writer) *slog.Logger {
	h := &SlogHandler{
		Handler: slog.NewJSONHandler(out, nil),
		l:       out,
	}
	return slog.New(h)
}

func main() {
	// Create an instance of the app structure
	app := &App{
		keys: 	  		&keyUtils.ECKeys{},
		envMap: 		make(map[string]any),
	}
	
	err := wails.Run(&options.App{
		Title:  "ProofChainV3",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
	if err:=godotenv.Load(".env");err!=nil{
		println("Error : ",err.Error())
	}
}
