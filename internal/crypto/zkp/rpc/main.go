package rpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "zkRPC/proto" 
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
		"1fc1ff4c3b4719c360ef5ddf39f197eb1dab304066fde2b2473fca7e26cc26de",
		"8743145b0aee03a6afa201b1ad1c1de7a68de62b506c3c150109d3b451e71e4a",
		"a8abc69d0802ca42a4412e9f90bb945a8256e954d8ca05764b947dc54ef5a5af",
		"9b532ff95737c1d24351afc8f902c2d38b371ef9ee5d54ef61a0242eeb50747a",
		"7e24da039f02a4bb16f3ef80f85827b501eb5ee88e087e69742c9a5660c3a2b5",
		"2e609d70aaf224556761b9f0938bf0c98d824ff577818f5486f252aec0e9c7dd",
		"96f5f60da8d3ff30d0bb8831e444725621e1cf426cfe41210cf4b6c887c541f1",
	}

	actualValue := "Maria"
	saltStr := "b9970a8b9940e7110ad29dae9029ae1c"
	rootStr := "94ac7ca53356bf8cf6901d73b747315f10f83426e545ef40a4f3df59ebc11c6c"

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