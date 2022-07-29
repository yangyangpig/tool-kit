package solana

import (
	"context"
	"os"
	"reflect"

	"log"

	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"
)

type Instruction struct {
	solanaCli   *rpc.Client
	wssolanaCli *ws.Client
}

func InitInstruction(solanaEnv string) *Instruction {
	wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
	if err != nil {
		panic(err)
	}
	return &Instruction{solanaCli: rpc.New(solanaEnv), wssolanaCli: wsClient}
}

// 异步发送
func (s *Instruction) TXSync(name string, commitment rpc.CommitmentType, instr []solana.Instruction,
	signerFunc func(key solana.PublicKey) *solana.PrivateKey, payer solana.PublicKey) error {
	recent, err := s.solanaCli.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("GetRecentBlockhash happen error %+v", err)
		return err
	}
	tx, err := solana.NewTransaction(
		instr,
		recent.Value.Blockhash,
		solana.TransactionPayer(payer),
	)
	if err != nil {
		log.Fatalf("NewTransaction happen error %+v", err)
		return err
	}
	log.Println("tx")
	//spew.Dump(tx)
	if _, err = tx.EncodeTree(text.NewTreeEncoder(os.Stdout, name)); err != nil {
		log.Fatalf("EncodeTree happen error %+v", err)
		return err
	}
	if _, err = tx.Sign(signerFunc); err != nil {
		return err
	}
	sig, err := s.solanaCli.SendTransactionWithOpts(
		context.Background(),
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: commitment,
		},
	)
	if err != nil {
		log.Fatalf("SendTransactionWithOpts happen error %+v", err)
		return err
	}
	log.Println("send come here-------------------------")
	spew.Dump(tx)
	sub, err := s.wssolanaCli.SignatureSubscribe(
		sig,
		commitment,
	)
	if err != nil {
		log.Fatalf("SignatureSubscribe happen error %+v", err)
		return err
	}
	defer sub.Unsubscribe()
	res, err := sub.Recv()
	if err != nil {
		return err
	}
	log.Println("res")
	spew.Dump(res)
	return nil
}
func (s *Instruction) CreateAccInstr(acc *solana.Wallet, accSize uint64,
	ownerPubKey solana.PublicKey, payer solana.PublicKey) (solana.Instruction, error) {
	rentMin, err := s.solanaCli.GetMinimumBalanceForRentExemption(
		context.TODO(),
		accSize,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return nil, err
	}
	return system.NewCreateAccountInstruction(
		rentMin,
		accSize,
		ownerPubKey,
		payer,
		acc.PublicKey(),
	).Build(), nil
}

func exampleFromGetTransaction() {
	endpoint := rpc.TestNet_RPC
	client := rpc.New(endpoint)

	txSig := solana.MustSignatureFromBase58("3pByJJ2ff7EQANKd2bgetmnYQxknk3QUib1xLMnrg6aCvg5hS78peaGMoceC9AFckomqrsgo38DpzrG2LPW9zj3g")
	{
		out, err := client.GetTransaction(
			context.TODO(),
			txSig,
			&rpc.GetTransactionOpts{
				Encoding: solana.EncodingBase64,
			},
		)
		if err != nil {
			panic(err)
		}

		tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
		if err != nil {
			panic(err)
		}

		decodeSystemTransfer(tx)
	}
}

func decodeSystemTransfer(tx *solana.Transaction) {
	spew.Dump(tx)

	// we know that the first instruction of the transaction is a `system` program instruction:
	i0 := tx.Message.Instructions[0]

	// parse a system program instruction:
	inst, err := system.DecodeInstruction(i0.ResolveInstructionAccounts(&tx.Message), i0.Data)
	if err != nil {
		panic(err)
	}
	// inst.Impl contains the specific instruction type (in this case, `inst.Impl` is a `*system.Transfer`)
	spew.Dump(inst)
	if _, ok := inst.Impl.(*system.Transfer); !ok {
		panic("the instruction is not a *system.Transfer")
	}

	// OR
	{
		// There is a more general instruction decoder: `solana.DecodeInstruction`.
		// But before you can use `solana.DecodeInstruction`,
		// you must register a decoder for each program ID beforehand
		// by using `solana.RegisterInstructionDecoder` (all solana-go program clients do it automatically with the default program IDs).
		decodedInstruction, err := solana.DecodeInstruction(
			system.ProgramID,
			i0.ResolveInstructionAccounts(&tx.Message),
			i0.Data,
		)
		if err != nil {
			panic(err)
		}
		spew.Dump(decodedInstruction)

		// decodedInstruction == inst
		if !reflect.DeepEqual(inst, decodedInstruction) {
			panic("they are NOT equal (this would never happen)")
		}

		// To register other (not yet registered decoders), you can add them with
		// `solana.RegisterInstructionDecoder` function.
	}

	{
		// pretty-print whole transaction:
		_, err := tx.EncodeTree(text.NewTreeEncoder(os.Stdout, text.Bold("TEST TRANSACTION")))
		if err != nil {
			panic(err)
		}
	}
}
