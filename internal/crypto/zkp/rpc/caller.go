package rpc

import (
	"context"
	"time"
	pb "github.com/Suy56/ProofChain/internal/crypto/zkp/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)
type ZKProverClient struct {
	client	pb.ProverServiceClient
	conn	*grpc.ClientConn
}

func NewZKProverClient(endpoint string) (*ZKProverClient, error) {
	kacp:= keepalive.ClientParameters{
		Time: 					10*time.Second,
		Timeout: 				2*time.Second,
		PermitWithoutStream: 	true,
	}

	conn,err:=grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
	)

	if err!=nil{
		return nil, err
	}

	return &ZKProverClient{
		client:pb.NewProverServiceClient(conn), 
		conn:conn,
	}, nil
}

func (c *ZKProverClient) Close() {
	c.conn.Close()
}

func (c *ZKProverClient) RequestMembershipProof(
	ctx context.Context,
	actualValue string,
	actualSalt string,
	allLeaves []string,
	publicList []string,
	publicRoot string,

)(*pb.ProofResponse, error){
	req := &pb.ProofRequest{
		ProofData: &pb.ProofRequest_Membership{
			Membership: &pb.MembershipRequest{
				ActualValue: actualValue,
				ActualSalt:  actualSalt,
				AllLeaves:   allLeaves,
				PublicList:  publicList,
				PublicRoot:  publicRoot,
			},
		},
	}
	return c.client.GenerateProof(ctx, req)
}

func (c *ZKProverClient) RequestRangeProof(
	ctx context.Context,
	actualValue string,
	actualSalt string,
	allLeaves []string,
	lowerBound uint32,
	upperBound uint32,
	publicRoot string,

)(*pb.ProofResponse, error){
	req := &pb.ProofRequest{
		ProofData: &pb.ProofRequest_Range{
			Range: &pb.RangeRequest{
				ActualValue: actualValue,
				ActualSalt:  actualSalt,
				AllLeaves:   allLeaves,
				LowerBound: lowerBound,
				UpperBound: upperBound,
				PublicRoot:  publicRoot,
			},
		},
	}

	return c.client.GenerateProof(ctx, req)
}