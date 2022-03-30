package functions

import (
	"log"

	"github.com/SoteriaTech/blockchain-functions/env"
	"github.com/SoteriaTech/blockchain-functions/eth"
	"github.com/SoteriaTech/blockchain-functions/helpers"
	"github.com/SoteriaTech/blockchain-functions/store"
	"github.com/SoteriaTech/blockchain-functions/utils"
)

// ScanEthBlock scan an ethereum block for transactions
func ScanEthBlock(h uint64, config *env.ChainConfig) (*eth.Header, error) {

	// get transactions from 3 blocks earlier from store
	conf := int(h) - config.Confirmations
	prevTxs, _ := store.Firestore.FindEthTransactionsFromBlockHeight(conf)
	if len(prevTxs) > 0 {
		var tbc []string
		for _, t := range prevTxs {
			tbc = append(tbc, t.TxHash)
		}
		hashes, _ := eth.ConfirmTransactions(tbc, conf)
		if len(hashes) > 0 {
			cTxs := helpers.FilterEthTransactionsByHash(prevTxs, hashes)
			if err := helpers.ConfirmEthTransactions(cTxs, config); err != nil {
				utils.ErrorReport.LogAndPrintError(err)
			}
		}
	}

	b, errB := eth.ScanBlock(h)
	if errB != nil {
		log.Fatal("error scan block")
		return nil, errB
	}

	accs, errAccs := store.Firestore.GetAllEthAccountAddresses()
	if errAccs != nil {
		return nil, errAccs
	}

	walletTxs := helpers.FilterEthTransactionsByAccountAddress(b.Txs, accs)
	for uid, t := range walletTxs {
		t.Confirmed = false
		exists, errTx := helpers.FindOrCreateEthTransaction(t)
		if errTx != nil {
			log.Println("error find or create tx &v", t)
			return nil, errTx
		}
		if exists != nil {
			continue
		}
		a := &store.EthAccountSchema{UID: uid, Address: t.To}
		log.Println("Transaction found for account: &v", helpers.SetCurrencyAmount(a, t))
	}

	return &b.Meta, nil
}
