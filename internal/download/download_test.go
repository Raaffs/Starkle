package download

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/Suy56/ProofChain/internal/models"
)
func TestExtractProofValues_WithExtra(t *testing.T) {
	// 1. Setup sample data with both Fixed and Extra fields
	input := DownloadProof{
		Name:      models.LeafFields{Hash: "h1", Key: "Name", Salt: "s1", Value: "Maria"},
		CertificateName:      models.LeafFields{Hash: "hn", Key: "CertificateName", Salt: "sn", Value: "UPDATED CERT"},
		BirthDate: models.LeafFields{Hash: "h2", Key: "BirthDate", Salt: "s2", Value: 1993},
		Address: models.LeafFields{Hash: "h3",Key: "Address", Salt: "s3", Value: "Tokyo, Japan"},
		Extra: map[string]models.LeafFields{
			"MembershipID": {
				Hash:  "m_hash",
				Key:   "MID",
				Salt:  "m_salt",
				Value: "GOLD_99",
			},
		},
	}
	wrapper := certificateWrapper{
        SaltedFields: input,
    }

	inputBytes,err:=json.Marshal(wrapper);if err!=nil{
		t.Fatal("error marshalling json: ",err)
	}
	d,err:=NewDownloader(inputBytes, slog.New(slog.NewJSONHandler(os.Stdout,nil)));if err!=nil{
		t.Fatal("error initializing downloader %w",err)
	}
	if err:=d.Exec();err!=nil{
		t.Fatal("Error downloading: ",err)
	}
}