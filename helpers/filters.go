package helpers

import (
	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/eth"
	"github.com/SoteriaTech/blockchain-functions/store"
)

// FilterBtcTransactionsByAccountAddress filter a list of transactions by a list of btc addresses
func FilterBtcTransactionsByAccountAddress(txs []*btc.Transaction, accs []*store.BtcAccountSchema) map[string]*store.BtcTransactionSchema {
	f := make(map[string]store.BtcAccountSchema, len(accs))
	out := make(map[string]*store.BtcTransactionSchema)
	for _, a := range accs {
		f[a.Address] = *a
	}
	for _, t := range txs {
		if acc, ok := f[t.Address]; ok {
			tx := &store.BtcTransactionSchema{
				To:          t.Address,
				TxHash:      t.Hash,
				Amount:      FromSatoshiToBtc(&t.Value),
				BlockHeight: t.BlockHeight,
				VoutIdx:     t.N,
			}

			out[acc.UID] = tx
		}
	}
	return out
}

// FilterEthTransactionsByAccountAddress filter a list of transactions by a list of btc addresses
func FilterEthTransactionsByAccountAddress(txs []*eth.Transaction, accs []*store.EthAccountSchema) map[string]*store.EthTransactionSchema {
	f := make(map[string]store.EthAccountSchema, len(accs))
	out := make(map[string]*store.EthTransactionSchema)
	for _, a := range accs {
		f[a.Address] = *a
	}
	for _, t := range txs {
		if acc, ok := f[t.Receiver]; ok {
			tx := &store.EthTransactionSchema{
				From:        t.From,
				To:          t.To,
				TxHash:      t.Hash,
				Amount:      t.Value.String(),
				BlockHeight: t.BlockHeight,
				LogIdx:      t.LogIdx,
				Receiver:    t.Receiver,
				Currency:    t.Currency,
			}
			out[acc.UID] = tx
		}
	}
	return out
}

// FilterBtcTransactionsByHash filter transactions by a slice of hashes
func FilterBtcTransactionsByHash(txs []*store.BtcTransactionSchema, hashes []string) (out []*store.BtcTransactionSchema) {
	f := make(map[string]*store.BtcTransactionSchema, len(txs))

	for _, t := range txs {
		f[t.TxHash] = t
	}

	for _, h := range hashes {
		if _, ok := f[h]; ok {
			out = append(out, f[h])
		}
	}
	return
}

// FilterEthTransactionsByHash filter transactions by a slice of hashes
func FilterEthTransactionsByHash(txs []*store.EthTransactionSchema, hashes []string) (out []*store.EthTransactionSchema) {
	f := make(map[string]*store.EthTransactionSchema, len(txs))

	for _, t := range txs {
		f[t.TxHash] = t
	}

	for _, h := range hashes {
		if _, ok := f[h]; ok {
			out = append(out, f[h])
		}
	}
	return
}
