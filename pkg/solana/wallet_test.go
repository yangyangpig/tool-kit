package solana

import (
	"fmt"
	"log"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var myWallet *Wallet

func init() {
	myWallet = InitWallet(rpc.DevNet_RPC)
}
func Test_CreateWallet(t *testing.T) {

	account, _ := myWallet.CreateAccount()
	err := myWallet.AirDropAndBalance(account, 1)
	if err != nil {
		log.Fatal(err)
	}

}

func Test_GetBalance(t *testing.T) {
	myWallet.GetBalance(solana.MustPublicKeyFromBase58("GDiHmJFergf2pFPSkpAQNNMqwPQbhbG1cU4ULoeKx3ff"))
}

func Test_AccountInfo(t *testing.T) {
	myWallet.GetAccountInfo(solana.MustPublicKeyFromBase58("AUJTadTPH9fCeV2egpcWfmJF2gSLTukXZd414RPRmuiX"))
}

/**
account private key: 4uRpmrXMToJs55d15bhzzxDnNKeW2buYW6Xdn1AThs1fdmKnrV5mkmkp4K6zVY9W1oYpLUGKPwXJi5rSAoUe3Ygi
account public key: 5JThDmMMbc7QKRLuNkcqEoVAcDv6LyGqifveSRDe9FQW
account private key: isDH5DXkvT2SUZE1QLeia7xWUiin2hbfuRwarBgT8eNEXzfJqYxGYRgjuFPmVXusXCFxSE89cPLUoQrnoAn759b
account public key: GDiHmJFergf2pFPSkpAQNNMqwPQbhbG1cU4ULoeKx3ff
*/

func Test_Transfer(t *testing.T) {
	// 进行交易或者获取账号信息，不能够创建账号，而是使用已经创建的账号pubkey，进行访问
	// fromAccount, _ := myWallet.CreateAccount()
	// toAccount, _ := myWallet.CreateAccount()
	fromAccount := solana.MustPrivateKeyFromBase58("4uRpmrXMToJs55d15bhzzxDnNKeW2buYW6Xdn1AThs1fdmKnrV5mkmkp4K6zVY9W1oYpLUGKPwXJi5rSAoUe3Ygi")
	// toAccount := solana.MustPrivateKeyFromBase58("isDH5DXkvT2SUZE1QLeia7xWUiin2hbfuRwarBgT8eNEXzfJqYxGYRgjuFPmVXusXCFxSE89cPLUoQrnoAn759b")
	toAccount := solana.MustPublicKeyFromBase58("AUJTadTPH9fCeV2egpcWfmJF2gSLTukXZd414RPRmuiX")
	fmt.Println("account public key:", fromAccount.PublicKey())
	fmt.Println("to account public key:", fromAccount.PublicKey())

	myWallet.transfer(fromAccount, toAccount)
}
