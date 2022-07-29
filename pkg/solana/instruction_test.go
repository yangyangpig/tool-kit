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
	fromAccountPub := solana.MustPublicKeyFromBase58("5JThDmMMbc7QKRLuNkcqEoVAcDv6LyGqifveSRDe9FQW")
	fromAccountPrivate := solana.MustPrivateKeyFromBase58("4uRpmrXMToJs55d15bhzzxDnNKeW2buYW6Xdn1AThs1fdmKnrV5mkmkp4K6zVY9W1oYpLUGKPwXJi5rSAoUe3Ygi")
	playerprivate := solana.MustPrivateKeyFromBase58("isDH5DXkvT2SUZE1QLeia7xWUiin2hbfuRwarBgT8eNEXzfJqYxGYRgjuFPmVXusXCFxSE89cPLUoQrnoAn759b")
	playerPubKey := solana.MustPublicKeyFromBase58("GDiHmJFergf2pFPSkpAQNNMqwPQbhbG1cU4ULoeKx3ff")
	stateAccount := solana.NewWallet()
	//toAccountInstr, _ := solanaInstruct.CreateAccInstr(stateAccount, AccessControllerStateAccountSize, toAccountPubKey, fromAccountPub)
	err := solanaInstruct.TXSync(
		"Hello OnChain",
		rpc.CommitmentFinalized,
		[]solana.Instruction{
			update_data.NewInitializeInstruction(300000, stateAccount.PublicKey(), fromAccountPub,
				solana.SystemProgramID).Build(),
		},
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(fromAccountPub) {
				return &fromAccountPrivate
			}
			if key.Equals(stateAccount.PublicKey()) {
				return &stateAccount.PrivateKey
			}
			if key.Equals(playerPubKey) {
				return &playerprivate
			}
			return nil
		},
		playerPubKey,
	)
	if err != nil {
		log.Fatal(err)
	}
}
