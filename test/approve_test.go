package main

import (
	"testing"

	blockchain "github.com/Suy56/ProofChain/chaincore/core"
)

func TestApprove(t *testing.T) {
	pk:="0x4c940bf3f77c3c9251582a3c7b3849a5d08b89ff72f91d0a9e47b74c4338297e"
	contract:="0xD4E5D3582E8c82b9a7e51AaC380a0A717df1217c"
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
	_, err:=i.Instance.ApproveVerifier(c.TxOpts,"inst");if err!=nil{
		t.Fatal(err)
	}	
}