package main

import (
	"testing"

	blockchain "github.com/Suy56/ProofChain/chaincore/core"
)

func TestApprove(t *testing.T) {
	pk:="0x4c940bf3f77c3c9251582a3c7b3849a5d08b89ff72f91d0a9e47b74c4338297e"
	contract:="0x56Dd04698844a363CcCae8322C19D42C9A5CE7fF"
	c:=	&blockchain.ClientConnection{}
	i:= &blockchain.ContractVerifyOperations{}
	host:="http://localhost:8545"
	if err:=blockchain.Init(
		c,
		i,
		pk[2:],
		contract,
		host,
	);err!=nil{
		t.Fatal(err)
	}
	_, err:=i.Instance.ApproveVerifier(c.TxOpts,"ins");if err!=nil{
		t.Fatal(err)
	}	
}