package solana

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/gagliardetto/solana-go"
)

type SolanaAccout struct {
	account   *solana.Wallet
	solanaCli *rpc.Client
}

func InitSolanaAccount(solanaEnv string) *SolanaAccout {
	return &SolanaAccout{solanaCli: rpc.New(solanaEnv), account: solana.NewWallet()}
}

func (s *SolanaAccout) CreateAccount() error {
	fmt.Println("account private key: ", s.account.PrivateKey)
	fmt.Println("account public key: ", s.account.PublicKey())

	// 投递一个sol币
	out, err := s.solanaCli.RequestAirdrop(context.TODO(), s.account.PublicKey(), solana.LAMPORTS_PER_SOL*1, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("air drop sol to solana chain error: %+v", err)
		return err
	}
	fmt.Println("airdrop transaction signature:", out)
	return nil
}

func (s *SolanaAccout) GetBalance() error {
	fmt.Println("account public key: ", s.account.PublicKey())

	balanceInfo, err := s.solanaCli.GetBalance(context.Background(), s.account.PublicKey(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("query balanc from solana chain error: %+v", err)
		return err
	}
	realSol := realsol(balanceInfo)
	fmt.Println("@", realSol.Text('f', 10))
	return nil

}

func (s *SolanaAccout) GetAccount() error {
	fmt.Println("account public key: ", s.account.PublicKey())
	accountInfo, err := s.solanaCli.GetAccountInfoWithOpts(context.TODO(), s.account.PublicKey(), &rpc.GetAccountInfoOpts{Encoding: solana.EncodingBase58})
	if err != nil {
		log.Fatalf("query accountinfo from solana chain error: %+v", err)
		return err
	}

	spew.Dump(accountInfo)
	return nil
}

// tool

func realsol(balanceInfo *rpc.GetBalanceResult) *big.Float {
	var lamportOnAccount = new(big.Float).SetUint64(uint64(balanceInfo.Value))
	var solBalance = new(big.Float).Quo(lamportOnAccount, new(big.Float)).SetUint64(solana.LAMPORTS_PER_SOL)
	return solBalance
}
