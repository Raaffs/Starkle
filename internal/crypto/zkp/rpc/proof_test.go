package rpc

import (
	"context"
	"time"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Suy56/ProofChain/internal/crypto/zkp/rpc/proto"
)

var allLeavesStrings = []string{
		"12e1970ebf67cc91c0a76aabe0107bd1e0427b89c5ec6e21d888dbd408652cb8",
		"55cece5b7c3eb9b4972758d38dd8763efa67d991c5311cd96ef05889a441c74e",
		"d05387dd84b0432a9fb3f24bf26e8e25da656d9e2241f2a089db5fb732a23fcc",
		"2686c2fca941f44f9dc41d733dd4ff9a960acbf64115c981ada87045f2f98e64",
		"8750ac38375e61bc2ad011d3888517642e9c53fb6ddac078880372e744349c27",
		"5d4be9af4ebcf4eb740e59a775f362873647dc41f236c8319bf5f4bcbc1cab69",
		"9abb011986c668d6ef31a58fca1ac09380ef0cd6cb8eb7f25c481165fa76d182",
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

	validRoot := "747dec436d75a6912e913324672a78dd7a14c40c6ef3acc2825c6adebb0116bc"
	validSalt := "8ee3d96cd121b2fc3235ce0e23ad800d"

	tests := []struct {
		name        string
		actualValue uint32
		salt        string
		root        string
		lower       uint32
		upper       uint32
		wantErr     bool
	}{
		{"ValidProof", 28, validSalt, validRoot, 18, 60, false},
		{"ValueTooLow", 10, validSalt, validRoot, 18, 60, true},
		{"ValueTooHigh", 70, validSalt, validRoot, 18, 60, true},
		{"WrongSalt", 28, "wrong_salt_12345", validRoot, 18, 60, true},
		{"WrongRootHash", 28, validSalt, "deadbeef12345678", 18, 60, true},
		{"BoundaryLower", 18, validSalt, validRoot, 18, 60, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.ProofRequest{
				ProofData: &pb.ProofRequest_Range{
					Range: &pb.RangeRequest{
						ActualValue: tt.actualValue,
						ActualSalt:  tt.salt,
						AllLeaves:   allLeavesStrings,
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
				t.Logf("✅ Receipt ID: %s | Cycles: %d", resp.ReceiptId, resp.Cycles)
			}
		})
	}
}

func TestSetMembershipProof(t *testing.T) {
	client, closeConn := getClient(t)
	defer closeConn()

	validRoot := "747dec436d75a6912e913324672a78dd7a14c40c6ef3acc2825c6adebb0116bc"
	validSalt := "e47a395cd43ec2ab68f0f902336053bc"
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
		{"ValueNotInList", "Steve", validSalt, validRoot, publicList, true},
		{"IncorrectSalt", "Maria", "bad_salt_999", validRoot, publicList, true},
		{"TamperedRoot", "Maria", validSalt, "fake_root_hash", publicList, true},
		{"EmptyPublicList", "Maria", validSalt, validRoot, []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.ProofRequest{
				ProofData: &pb.ProofRequest_Membership{
					Membership: &pb.MembershipRequest{
						ActualValue: tt.actualValue,
						ActualSalt:  tt.salt,
						AllLeaves:   allLeavesStrings,
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
				t.Logf("✅ Receipt ID: %s | Cycles: %d", resp.ReceiptId, resp.Cycles)
			}
		})
	}
}