package solana

import (
	"log"
	"testing"
	"toolkit/pkg/solana/idl/update_data"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var solanaInstruct *Instruction

func init() {
	solanaInstruct = InitInstruction(rpc.DevNet_RPC)
}

const (
	Discriminator = 8
	// TokenMintAccountSize default size of data required for a new mint account
	TokenMintAccountSize             = uint64(82)
	TokenAccountSize                 = uint64(165)
	AccessControllerStateAccountSize = uint64(Discriminator + solana.PublicKeyLength + solana.PublicKeyLength + 8 + 32*64)
	StoreAccountSize                 = uint64(Discriminator + solana.PublicKeyLength*3)
	OCRTransmissionsAccountSize      = uint64(Discriminator + 192 + 8192*48)
	OCRProposalAccountSize           = Discriminator + 1 + 32 + 1 + 1 + (1 + 4) + 32 + ProposedOraclesSize + OCROffChainConfigSize
	ProposedOracleSize               = uint64(solana.PublicKeyLength + 20 + 4 + solana.PublicKeyLength)
	ProposedOraclesSize              = ProposedOracleSize*19 + 8
	OCROracle                        = uint64(solana.PublicKeyLength + 20 + solana.PublicKeyLength + solana.PublicKeyLength + 4 + 8)
	OCROraclesSize                   = OCROracle*19 + 8
	OCROffChainConfigSize            = uint64(8 + 4096 + 8)
	OCRConfigSize                    = 32 + 32 + 32 + 32 + 32 + 32 + 16 + 16 + (1 + 1 + 2 + 4 + 4 + 32) + (4 + 32 + 8) + (4 + 4)
	OCRAccountSize                   = Discriminator + 1 + 1 + 2 + 4 + solana.PublicKeyLength + OCRConfigSize + OCROffChainConfigSize + OCROraclesSize
)

func Test_HandleOnChain(t *testing.T) {

	fromAccount, _ := solana.WalletFromPrivateKeyBase58("4kgmdJWXK3gLdVKJd9oyD2exuuj3k7dvSPBcekhBxAhUCkLHMd7tf65B4JjBjt7agbZUrfuZnJjUm7NfNvRBXtHL")
	payerAccount, _ := solana.WalletFromPrivateKeyBase58("6rk91t9suXoRvhbRwnviUeiEeo6UcUbWKkaPsERg29Z9h1Tucwkr1pYBXgiHYiUazq7U32DcexthL7MGi1M5U2e")
	// programId
	//programIdPubKey := solana.MustPublicKeyFromBase58("AUJTadTPH9fCeV2egpcWfmJF2gSLTukXZd414RPRmuiX")
	// seedInstruction, _ := CreateSeedInstruction(payerAccount, programIdPubKey, "helloworld")
	// seedAccount, _ := CreateSeedAccount(payerAccount, programIdPubKey, "helloworld")
	// fmt.Println("seedAccount", seedAccount)
	// _ = solanaInstruct.AirDropAndBalance(seedAccount, 1)

	// 首先要创建一个子账号，子账号需要关联合约账号
	seedAccount := solana.NewWallet()
	//myAccountInstr, _ := solanaInstruct.CreateAccInstr(seedAccount, TokenAccountSize, programIdPubKey, payerAccount.PublicKey())
	err := solanaInstruct.TXSync(
		"Hello OnChain",
		rpc.CommitmentFinalized,
		[]solana.Instruction{
			//myAccountInstr,
			update_data.NewInitializeInstruction(300, seedAccount.PublicKey(), fromAccount.PublicKey(),
				solana.SystemProgramID).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(fromAccount.PublicKey()) {
				return &fromAccount.PrivateKey
			}
			if key.Equals(payerAccount.PublicKey()) {
				return &payerAccount.PrivateKey
			}
			if key.Equals(seedAccount.PublicKey()) {
				return &seedAccount.PrivateKey
			}
			return nil
		},
		payerAccount.PublicKey(),
	)
	if err != nil {
		log.Fatal(err)
	}
}
