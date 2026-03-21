package main

import (
	"testing"

	blockchain "github.com/Suy56/ProofChain/chaincore/core"
)

func TestApprove(t *testing.T) {
	pk:="0x4c940bf3f77c3c9251582a3c7b3849a5d08b89ff72f91d0a9e47b74c4338297e"
	contract:="0x1Da26345D7d4e24a148C2B49E9b2C8af2e3e5239"
	c:=	&blockchain.ClientConnection{}
	i:= &blockchain.ContractVerifyOperations{}
	host:="http://localhost:7545"
	if err:=blockchain.Init(
		c,
		i,
		pk[2:],
		contract,
		host,
	);err!=nil{
		t.Fatal(err)
	}
	_, err:=i.Instance.ApproveVerifier(c.TxOpts,"ins4");if err!=nil{
		t.Fatal(err)
	}	
}