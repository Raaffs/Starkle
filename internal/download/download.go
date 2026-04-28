package download

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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

func New(certificate DownloadProof, logger *slog.Logger, user string) (*Downloader, error) {
	basePath, err := utils.GetDirPath("Downloads")
	if err != nil {
		return nil, err
	}

	certName,ok:=certificate.CertificateName.Value.(string);if !ok{
		return nil, fmt.Errorf("invalid certificate name value. Expected: string got %v", certificate.CertificateName.Value)
	}

	finalDir := filepath.Join(basePath, user , certName)

	return &Downloader{
		TargetDir: finalDir,
		ProofData: certificate,
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