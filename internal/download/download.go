package download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Suy56/ProofChain/internal/models"
	"github.com/Suy56/ProofChain/internal/utils"
)

// DownloadProof is the structure that holds the necessary data for reconstructing the proofs for each field.
// It mirrors the CertificateData structure but with models.LeafFields values instead of plain strings.
type DownloadProof = models.CertificateBase[models.LeafFields]

type Downloader struct {
	TargetDir string
	ProofData DownloadProof
	logger    *slog.Logger
}

func NewDownloader(certificate []byte, logger *slog.Logger) (*Downloader, error) {
	var cert models.CertificateBase[models.LeafFields]

	if err := json.Unmarshal(certificate, &cert); err != nil {
		return nil, fmt.Errorf("could not decode certificate proof: %w", err)
	}
	basePath, err := getDownloadDir()
	if err != nil {
		return nil, err
	}

	val,ok:=cert.CertificateName.Value.(string);if !ok{
		return nil, fmt.Errorf("invalid certificate name value. Expected: string got %v", cert.CertificateName.Value)
	}

	finalDir := filepath.Join(basePath, val)

	return &Downloader{
		TargetDir: finalDir,
		ProofData: cert,
		logger:    logger,
	}, nil
}

func (d *Downloader) Exec() error {
	var errs []error
	for k, v := range utils.Walk(d.ProofData) {
		proofK := d.extractProofValues(k, v)
		if err := d.store(k, proofK); err != nil {
			d.logger.Error("Failed to store field proof", "field", k, "directory", d.TargetDir, "error", err)
			errs = append(errs, err)
			continue
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("completed with %d failures", len(errs))
	}
	return nil
}

func (d *Downloader) extractProofValues(activeKey string, fullValue any) map[string]models.LeafFields {
	slim := func(f models.LeafFields) models.LeafFields {
		return models.LeafFields{Hash: f.Hash, Key: f.Key, Value: ""}
	}

	v := d.ProofData
	result := map[string]models.LeafFields{
		"Address":         slim(v.Address),
		"Age":             slim(v.Age),
		"BirthDate":       slim(v.BirthDate),
		"CertificateName": slim(v.CertificateName),
		"Name":            slim(v.Name),
		"PublicAddress":   slim(v.PublicAddress),
		"UniqueID":        slim(v.UniqueID),
	}

	for k, val := range v.Extra {
		result[k] = slim(val)
	}
	hf, ok := fullValue.(models.LeafFields)
	if ok {
		result[activeKey] = hf
	}
	return result
}

func (d *Downloader) store(key string, proof any) error {
	if err := os.MkdirAll(d.TargetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(proof, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal proof JSON: %w", err)
	}

	filename := filepath.Join(d.TargetDir, key+".json")
	return os.WriteFile(filename, data, 0644)
}


func getDownloadDir() (string, error) {
	var downloadDir string
	// 1. Try Linux standard
	cmd := exec.Command("xdg-user-dir", "DOWNLOAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		downloadDir = strings.TrimSpace(out.String())
	}

	// 2. Fallback for macOS/Other
	if downloadDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to detect home directory: %w", err)
		}
		downloadDir = filepath.Join(home, "Downloads")
	}

	finalPath := filepath.Join(downloadDir, "ProofChain")
	if err := os.MkdirAll(finalPath, 0755); err != nil {
		return "", fmt.Errorf("failed to ensure ProofChain directory: %w", err)
	}

	return finalPath, nil
}

