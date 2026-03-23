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

	"github.com/Suy56/ProofChain/internal/utils"
	"github.com/Suy56/ProofChain/storage/models"
)

type hashedField struct {
	Hash  string `json:"hash"`
	Key   string `json:"key"`
	Salt  string `json:"salt"`
	Value string `json:"value"`
}

// DownloadProof is the structure that holds the necessary data for reconstructing the proofs for each field. 
// It mirrors the CertificateData structure but with hashedField values instead of plain strings.
// Which is why we're using a generic base structure
type DownloadProof = models.CertificateBase[hashedField]

// Downloader manages the lifecycle of exporting proof files
type Downloader struct {
	TargetDir string
	ProofData DownloadProof
	logger    *slog.Logger
}

type DocumentWrapper struct {
	SaltedFields DownloadProof `json:"salted_fields"`
}

// NewDownloader initializes the downloader, determines the path, and unmarshals the data
func NewDownloader(document []byte, logger *slog.Logger) (*Downloader, error) {
	var wrapper DocumentWrapper
	if err := json.Unmarshal(document, &wrapper); err != nil {
		return nil, fmt.Errorf("could not decode certificate proof: %w", err)
	}
	doc := wrapper.SaltedFields
	basePath, err := getDownloadDir()
	if err != nil {
		return nil, err
	}
	// Calculate the specific folder for this certificate
	finalDir := filepath.Join(basePath, doc.CertificateName.Value)

	return &Downloader{
		TargetDir: finalDir,
		ProofData: doc,
		logger:    logger,
	}, nil
}

// Exec starts the download process for all fields in the ProofChain
func (d *Downloader) Exec() error {
	var errs []error

	for k, v := range utils.Walk(d.ProofData) {
		proofK := d.extractProofValues(k, v)
		if err := d.store(k, proofK); err != nil {
			d.logger.Error("Failed to store field proof", "field", k, "directory", d.TargetDir, "error", err)
			errs = append(errs, err)
			continue
		}
		d.logger.Info("Field proof saved", "field", k)
	}
	if len(errs) > 0 {
		return fmt.Errorf("completed with %d failures", len(errs))
	}
	return nil
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

func (d *Downloader) extractProofValues(activeKey string, fullValue any) map[string]hashedField {
	slim := func(f hashedField) hashedField {
		return hashedField{Hash: f.Hash, Key: f.Key}
	}

	v := d.ProofData
	result := map[string]hashedField{
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

	hf, ok := fullValue.(hashedField)
	if ok {
		result[activeKey] = hf
	}
	return result
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
