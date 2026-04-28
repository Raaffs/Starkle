package main

import (
	"log"
	"testing"

	"github.com/Suy56/ProofChain/chaincore/core"
)
var ETH_CLIENT_URL="http://localhost:8545"

func TestDeploy(t *testing.T){
	app:=&struct{
		conn *blockchain.ClientConnection
		in *blockchain.ContractVerifyOperations
	}{
		conn: &blockchain.ClientConnection{},
		in: &blockchain.ContractVerifyOperations{},
	}

	privateKey:="0x4c940bf3f77c3c9251582a3c7b3849a5d08b89ff72f91d0a9e47b74c4338297e"
	if err:=blockchain.Init(app.conn,app.in,privateKey[2:],"",ETH_CLIENT_URL);err!=nil{
		t.Fatal(err)
	}
	contract,_,err:=blockchain.Deploy(app.conn.TxOpts,app.conn.Client);if err!=nil{
		log.Println("contract:",contract)
		t.Fatal(err)
	}
	log.Println(contract)
}