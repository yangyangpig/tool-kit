package solana

import (
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
)

var initSolanaAccount *SolanaAccout

func init() {
	initSolanaAccount = InitSolanaAccount(rpc.DevNet_RPC)
}

func Test_CreateSolanaAccount(t *testing.T) {
	_ = initSolanaAccount.CreateAccount()

	// 获取剩余资产
	_ = initSolanaAccount.GetBalance()
	// 获取帐号信息
	_ = initSolanaAccount.GetAccount()
	
}