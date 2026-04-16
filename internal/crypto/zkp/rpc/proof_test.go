package rpc

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Suy56/ProofChain/internal/crypto/zkp"
	pb "github.com/Suy56/ProofChain/internal/crypto/zkp/rpc/proto"
)

var allLeavesStrings = []string{
		"f0954e34ed538fdd71816223a5c37d3786c3312f0b39cc3363bef13cbca4c533",
		"18c7606c52b9a88b6d6ecb75bbf15094d01109af31c7480998029de76455d38f",
		"c181f04c87cce736752fc482e70e79734d7026bd81eb49b5351ca84dd22b5e72",
		"28d138788ad5ca9a8935ff5d3e9c62508922cde37e9ed81f8fc4e8bba4d95f8a",
		"b2800bb4502edfae9a7c9bb31309d7a67072559592d70e1a6468aac4f8e1fcf0",
		"79c161309c99aacab26c2dd19f8f4151304a364d02208128d77b1339b7b39f95",
		"c79c9bd4139169d8385b6cd7f19916528d07aad3580e300065b515a3dc972cab",
	}


// Helper to get gRPC client
func getClient(t *testing.T) (pb.ProverServiceClient, func()) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	client := pb.NewProverServiceClient(conn)
	return client, func() { conn.Close() }
}

func TestRangeProof(t *testing.T) {
	client, closeConn := getClient(t)
	defer closeConn()

	validRoot := "d9143e35153df4c3ef748d4b4c1f476bd9c45ce735eef92e515e9dd8f4a605c9"
	validSalt := "afa4b107459a892dc4e26d7227ad3d5b"

	tests := []struct {
		name        string
		actualValue uint32
		salt        string
		root        string
		lower       uint32
		upper       uint32
		wantErr     bool
	}{
		{"ValidProof", 24, validSalt, validRoot, 18, 60, false},
		{"false salt",24,"fkejf",validRoot,18,60,true},
		{"false root",24,validSalt,"fdfdfef",18,60,true},
		{"false actual value",25,validSalt,validRoot,18,60,true},
		{"lower bound violated",24,validSalt,validRoot,25,60,true},
		{"upper bound violated",24,validSalt,validRoot,18,20,true},
	}

	root,siblings:=zkp.GenerateMerklePath(allLeavesStrings,"18c7606c52b9a88b6d6ecb75bbf15094d01109af31c7480998029de76455d38f")
	if root!=validRoot{
		t.Fatalf("merkle path provided invalid root.\nexpected: %s got: %s",validRoot,root)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.ProofRequest{
				ProofData: &pb.ProofRequest_Range{
					Range: &pb.RangeRequest{
						ActualValue: tt.actualValue,
						ActualSalt:  tt.salt,
						Siblings:   siblings,
						LowerBound:  tt.lower,
						UpperBound:  tt.upper,
						PublicRoot:  tt.root,
					},
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			resp, err := client.GenerateProof(ctx, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				t.Logf("Receipt ID: %s | Cycles: %d", resp.ReceiptId, resp.Cycles)
			}
		})
	}
}

func TestSetMembershipProof(t *testing.T) {
	client, closeConn := getClient(t)
	defer closeConn()

	validRoot := "d9143e35153df4c3ef748d4b4c1f476bd9c45ce735eef92e515e9dd8f4a605c9"
	validSalt := "6a8fd0e5db7ca293dfb777099887c100"
	publicList := []string{"Maria", "Mark", "John", "Maquia"}

	tests := []struct {
		name        string
		actualValue string
		salt        string
		root        string
		list        []string
		wantErr     bool
	}{
		{"ValidMember", "Maria", validSalt, validRoot, publicList, false},
	}

	root,siblings:=zkp.GenerateMerklePath(allLeavesStrings,"b2800bb4502edfae9a7c9bb31309d7a67072559592d70e1a6468aac4f8e1fcf0")
	if root!=validRoot{
		t.Fatalf("merkle path provided invalid root.\nexpected: %s got: %s",validRoot,root)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.ProofRequest{
				ProofData: &pb.ProofRequest_Membership{
					Membership: &pb.MembershipRequest{
						ActualValue: tt.actualValue,
						ActualSalt:  tt.salt,
						Siblings:   siblings,
						PublicList:  tt.list,
						PublicRoot:  tt.root,
					},
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			resp, err := client.GenerateProof(ctx, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				t.Logf("Receipt ID: %s | Cycles: %d", resp.ReceiptId, resp.Cycles)
			}
		})
	}
}