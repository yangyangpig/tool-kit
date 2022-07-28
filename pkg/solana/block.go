package solana

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Block struct {
	solanaCli *rpc.Client
}

func InitBlock(solanaEnv string) *Block {
	return &Block{solanaCli: rpc.New(solanaEnv)}
}

func (s *Block) GetBlock() error {
	// 返回最新的block的hash数据，This method is only available in solana-core v1.9 or newer. Please use getRecentBlockhash for solana-core v1.8
	example, err := s.solanaCli.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("query recent blocks hash happen error %+v", err)
		return err
	}
	spew.Dump(example)
	// 获取block信息
	out, err := s.solanaCli.GetBlock(context.TODO(), uint64(example.Context.Slot))
	if err != nil {
		log.Fatalf("query blocks info happen error %+v", err)
		return err
	}

	spew.Dump(len(out.Transactions))

	// 根据筛选条件获取block信息
	includeRewards := false
	outWithOpts, err := s.solanaCli.GetBlockWithOpts(context.TODO(), uint64(example.Context.Slot), &rpc.GetBlockOpts{
		Encoding:           solana.EncodingBase64,
		Commitment:         rpc.CommitmentFinalized,
		TransactionDetails: rpc.TransactionDetailsSignatures,
		Rewards:            &includeRewards,
	})
	if err != nil {
		log.Fatalf("query blocks info with opts happen error %+v", err)
		return err
	}
	spew.Dump(outWithOpts)
	return nil

}
