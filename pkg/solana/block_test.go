package solana

import (
	"log"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
)

var blockInstance *Block

func init() {
	blockInstance = InitBlock(rpc.DevNet_RPC)
}

func Test_GetBlock(t *testing.T) {
	err := blockInstance.GetBlock()
	if err != nil {
		log.Fatal(err)
	}
}
