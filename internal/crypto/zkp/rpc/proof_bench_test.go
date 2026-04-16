package rpc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/Suy56/ProofChain/internal/crypto/zkp"
	pb "github.com/Suy56/ProofChain/internal/crypto/zkp/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SaltedField struct {
	Hash  string `json:"hash"`
	Value string `json:"value"`
	Salt  string `json:"salt"`
}

type CertificateData struct {
	SaltedFields map[string]SaltedField `json:"salted_fields"`
}

var benchmarkPublicList = []string{"val_0", "alice", "ID-999", "Master of Electrical Engineering"}

func HashData(data ...[]byte) string {
	h := sha256.New()
	for _, d := range data {
		h.Write(d)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateMerklePath builds the tree and returns the root + specific siblings for the target
func BenchmarkProverMembership(b *testing.B) {
	conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	client := pb.NewProverServiceClient(conn)

	scenarios := []int{20, 40, 128}

	for _, count := range scenarios {
		b.Run(fmt.Sprintf("Leaves-%d", count), func(b *testing.B) {
			path := filepath.Join("test_sample", fmt.Sprintf("salted_cert_%d.json", count))
			fileData, err := os.ReadFile(path)
			if err != nil { b.Fatalf("File missing: %v", err) }

			var cert CertificateData
			json.Unmarshal(fileData, &cert)

			var allLeaves []string
			for _, f := range cert.SaltedFields {
				allLeaves = append(allLeaves, f.Hash)
			}
			
			target := cert.SaltedFields["Name"]
			root, siblings := zkp.GenerateMerklePath(allLeaves, target.Hash);
			log.Print(root,siblings)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				req := &pb.ProofRequest{
					ProofData: &pb.ProofRequest_Membership{
						Membership: &pb.MembershipRequest{
							ActualValue: target.Value,
							ActualSalt:  target.Salt,
							Siblings:   siblings, // ONLY THE PATH
							PublicList:  benchmarkPublicList,
							PublicRoot:  root,
						},
					},
				}
				_, err := client.GenerateProof(context.Background(), req)
				if err != nil { b.Fatalf("Fail at %d: %v", count, err) }
			}
		})
	}
}

