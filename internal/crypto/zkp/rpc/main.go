package rpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Suy56/ProofChain/internal/crypto/zkp/rpc/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewProverServiceClient(conn)

	// 1. Keep hashes as hex strings
	allLeavesStrings := []string{
		"12e1970ebf67cc91c0a76aabe0107bd1e0427b89c5ec6e21d888dbd408652cb8",
		"55cece5b7c3eb9b4972758d38dd8763efa67d991c5311cd96ef05889a441c74e",
		"d05387dd84b0432a9fb3f24bf26e8e25da656d9e2241f2a089db5fb732a23fcc",
		"2686c2fca941f44f9dc41d733dd4ff9a960acbf64115c981ada87045f2f98e64",
		"8750ac38375e61bc2ad011d3888517642e9c53fb6ddac078880372e744349c27",
		"5d4be9af4ebcf4eb740e59a775f362873647dc41f236c8319bf5f4bcbc1cab69",
		"9abb011986c668d6ef31a58fca1ac09380ef0cd6cb8eb7f25c481165fa76d182",
	}

	actualValue := "Maria"
	saltStr := "e47a395cd43ec2ab68f0f902336053bc"
	rootStr := "747dec436d75a6912e913324672a78dd7a14c40c6ef3acc2825c6adebb0116bc"

	// 2. Construct Request using pure strings
	req := &pb.ProofRequest{
		ProofData: &pb.ProofRequest_Membership{
			Membership: &pb.MembershipRequest{
				ActualValue: actualValue,
				ActualSalt:  saltStr,
				AllLeaves:   allLeavesStrings,
				PublicList:  []string{"Mark", "Maria", "John", "Sam"},
				PublicRoot:  rootStr,
			},
		},
	}

	fmt.Printf("🚀 Requesting Membership Proof (String Mode) for: %s\n", actualValue)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	start := time.Now()
	resp, err := client.GenerateProof(ctx, req)
	if err != nil {
		log.Fatalf("❌ Prover Error: %v", err)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("✅ Success! Proof Generated.\n")
	fmt.Printf("Receipt ID:  %s\n", resp.ReceiptId)
	fmt.Printf("Cycles Used: %d\n", resp.Cycles)
	fmt.Printf("Duration:    %v\n", time.Since(start))
	fmt.Println("--------------------------------------------------")
}