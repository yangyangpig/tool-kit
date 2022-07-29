package solana

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	_ "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"
)

// 主要用于对钱包的操作

type Wallet struct {
	solanaCli   *rpc.Client
	wssolanaCli *ws.Client
}

func (w *Wallet) CreateAccount() (solana.PrivateKey, error) {
	account := solana.NewWallet()
	fmt.Println("account private key:", account.PrivateKey)
	fmt.Println("account public key:", account.PublicKey())
	return account.PrivateKey, nil
}

func InitWallet(solanaEnv string) *Wallet {

	wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
	if err != nil {
		panic(err)
	}

	return &Wallet{
		solanaCli:   rpc.New(solanaEnv),
		wssolanaCli: wsClient,
	}
}

// 投放
func (w *Wallet) AirDropAndBalance(account solana.PrivateKey, num uint64) error {

	// 投放5sol币到新账号
	out, err := w.solanaCli.RequestAirdrop(context.TODO(), account.PublicKey(), solana.LAMPORTS_PER_SOL*num, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("airdrop happen error %+v", err)
		return err
	}

	fmt.Println("airdrop transaction signature:", out)

	return nil
}

// 获取账号信息
func (w *Wallet) GetBalance(pubKey solana.PublicKey) error {
	out, err := w.solanaCli.GetBalance(context.TODO(), pubKey, rpc.CommitmentProcessed)
	if err != nil {
		return err
	}
	spew.Dump(out)
	spew.Dump(out.Value)
	solBalance := getRealSol(out)
	fmt.Println("◎", solBalance.Text('f', 10))
	return nil
}

func (w *Wallet) GetAccountInfo(pubKey solana.PublicKey) error {
	// 获取账号信息
	fmt.Println("account public key:", pubKey)
	accountInfo, err := w.solanaCli.GetAccountInfoWithOpts(context.Background(), pubKey,
		&rpc.GetAccountInfoOpts{Encoding: solana.EncodingBase58})
	if err != nil {
		log.Fatalf("airdrop after accountInfo happen error %+v", err)
		return err
	}
	spew.Dump(accountInfo)
	return nil
}

// 转账
func (w *Wallet) transfer(fromAccount solana.PrivateKey, toAccountPubKey solana.PublicKey) error {
	amount := uint64(10000000)
	fmt.Println("from account public key:", fromAccount.PublicKey())
	fmt.Println("to account public key:", toAccountPubKey)
	recent, err := w.solanaCli.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("transfer get recent block hash happen error %+v", err)
		return err
	}
	// 创建一个交易，交易数量为 3333 lamports, 1 sol = 1000000000 lamports
	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(amount, fromAccount.PublicKey(), toAccountPubKey).Build(),
	}, recent.Value.Blockhash, solana.TransactionPayer(fromAccount.PublicKey()))
	if err != nil {
		log.Fatalf("transfer get recent block hash happen error %+v", err)
		return err
	}
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if fromAccount.PublicKey().Equals(key) {
			return &fromAccount
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Sign happen error %+v", err)
		return err
	}
	spew.Dump(tx)

	// 打印交易信息
	tx.EncodeToTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))

	// 发送交易并且等待确认
	// sig, err := w.solanaCli.SendTransactionWithOpts(
	// 	context.TODO(),
	// 	tx,
	// 	true,
	// 	rpc.CommitmentFinalized,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	sig, err := confirm.SendAndConfirmTransaction(context.TODO(), w.solanaCli, w.wssolanaCli, tx)
	if err != nil {
		log.Fatalf("send and confirm transaction happen error %+v", err)
		return err
	}
	spew.Dump(sig)
	return nil

}



func getRealSol(balanceInfo *rpc.GetBalanceResult) *big.Float {
	var lamportsOnAccount = new(big.Float).SetUint64(uint64(balanceInfo.Value))
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
	return solBalance
}
