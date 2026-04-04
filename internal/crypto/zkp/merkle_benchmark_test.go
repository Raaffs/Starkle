package zkp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"github.com/Suy56/ProofChain/internal/models"
)

// Helper to ensure directory exists
func ensureDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
}

func benchmarkCert(totalFields int) models.CertificateBase[any] {
cert := models.CertificateBase[any]{
        Address:         "1",
        Age:             "23",
        BirthDate:       "1990-01-01",
        CertificateName: "Master of Electrical Engineering",
        Name:            "alice",
        PublicAddress:   "0x123...",
        UniqueID:        "ID-999",
        Extra:           make(map[string]any),
    }
	extraNeeded := totalFields - 7
	if extraNeeded > 0 {
		for i := 0; i < extraNeeded; i++ {
			key := fmt.Sprintf("ext_%d", i)
			cert.Extra[key] = fmt.Sprintf("val_%d", i)
		}
	}
	return cert
}

func BenchmarkMerkleScalability(b *testing.B) {
	scenarios := []struct {
		name   string
		target int
	}{
		{"N_5", 5},
		{"N_20", 20},
		{"N_40", 40},
		{"N_70", 70},
		{"N_100", 100},
		{"N_128", 128},
	}

	// Prepare result directory
	resultsDir := "test_results"
	ensureDir(resultsDir)

	for _, s := range scenarios {
		b.Run(s.name, func(b *testing.B) {
			input := benchmarkCert(s.target)
			merkle := NewMerkleProof()
			
			// We only need to save the JSON ONCE per scenario, not every iteration
			// So we do it before the loop
			_, salted, _ := merkle.GenerateRootProof(input)
			
			fileName := filepath.Join(resultsDir, fmt.Sprintf("salted_cert_%d.json", s.target))
			file, _ := json.MarshalIndent(salted, "", "  ")
			_ = os.WriteFile(fileName, file, 0644)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = merkle.GenerateRootProof(input)
			}
		})
	}
}