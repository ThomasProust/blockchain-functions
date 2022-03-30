package api

import (
	"context"
	"math/big"
	"net/http"
	"time"

	ethinfura "github.com/INFURA/go-ethlibs/eth"
	"github.com/INFURA/go-ethlibs/node"
	"github.com/SoteriaTech/blockchain-functions/eth"
)

// InfuraClient structure of the InfuraClient
type InfuraClient struct {
	client node.Client
	ctx    context.Context
	*http.Client
}

// Infura instance of the infura client
var Infura *InfuraClient

// InitInfuraClient initialize an instance of InfuraClient
func InitInfuraClient(endpoint string) {
	ctx := context.Background()
	client, _ := node.NewClient(ctx, endpoint)
	Infura = &InfuraClient{
		client: client,
		ctx:    ctx,
		Client: &http.Client{},
	}
}

// GetBlockHeader get the height of the latest block
func (i *InfuraClient) GetBlockHeader() (uint64, error) {

	bh, err := i.client.BlockNumber(i.ctx)
	if err != nil {
		return 0, err
	}
	return bh, nil
}

// GetTransactionByHash get a transaction from its hash
func (i *InfuraClient) GetTransactionByHash(hash string, b int) (*eth.Transaction, error) {
	t, err := i.client.TransactionByHash(i.ctx, hash)
	if err != nil {
		return nil, err
	}
	tx := &eth.Transaction{
		Hash:        t.Hash.String(),
		From:        t.From.String(),
		To:          t.To.String(),
		Value:       t.Value.Big(),
		BlockHeight: int(t.BlockNumber.UInt64()),
		Receiver:    t.To.String(),
		Currency:    "ETH",
	}

	return tx, nil
}

// GetTransactionsFromBlock get the transactions from the block body
func (i *InfuraClient) GetTransactionsFromBlock(h *big.Int) (*eth.BlockData, error) {

	block, err := i.client.BlockByNumber(i.ctx, h.Uint64(), true)
	if err != nil {
		return nil, err
	}

	var txs []*eth.Transaction
	for _, t := range block.Transactions {
		tx := &eth.Transaction{
			Hash:        t.Hash.String(),
			From:        t.From.String(),
			To:          to(t),
			Value:       t.Value.Big(),
			BlockHeight: int(t.BlockNumber.UInt64()),
			Receiver:    to(t),
		}
		txs = append(txs, tx)
	}
	bd := &eth.BlockData{
		Txs: txs,
		Meta: eth.Header{
			Hash:        block.Hash.String(),
			Height:      int(block.Number.Int64()),
			LastUpdated: time.Now(),
			Time:        int(block.Timestamp.UInt64()),
			Nonce:       block.Nonce.String(),
		},
	}

	return bd, nil
}

// GetReceipt get the receipt of a transaction
func (i *InfuraClient) GetReceipt(hash string) (*ethinfura.TransactionReceipt, error) {
	r, err := i.client.TransactionReceipt(i.ctx, hash)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func to(t ethinfura.TxOrHash) string {
	if t.To == nil {
		return ""
	}
	return t.To.String()
}
