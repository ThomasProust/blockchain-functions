package helpers

import (
	"log"
	"math/big"
	"strconv"

	"github.com/SoteriaTech/blockchain-functions/env"
	"github.com/SoteriaTech/blockchain-functions/store"
)

// FindOrCreateBtcTransaction find a btc transaction and returns it, or create it if not exist and returns nothing
func FindOrCreateBtcTransaction(t *store.BtcTransactionSchema) (tx *store.BtcTransactionSchema, err error) {
	tx, err = store.Firestore.FindBtcTransaction(t.TxHash + strconv.Itoa(t.VoutIdx))
	if err != nil {
		log.Fatal(err)
		return
	}

	err = store.Firestore.CreateBtcTransaction(t)
	return
}

// FindOrCreateEthTransaction find an ethereum transaction and returns it, or create it if not exist and returns nothing
func FindOrCreateEthTransaction(t *store.EthTransactionSchema) (tx *store.EthTransactionSchema, err error) {
	tx, err = store.Firestore.FindEthTransaction(t.TxHash + t.LogIdx)
	if err != nil {
		log.Fatal(err)
		return
	}
	if tx != nil {
		return
	}
	err = store.Firestore.CreateEthTransaction(t)
	return
}

// UpdateAccountBtcBalance update the btc balance of a user UID by a given amount
func UpdateAccountBtcBalance(uid string, amount *big.Float) (*big.Float, error) {
	bal, errBal := store.Firestore.FindBtcBalance(uid)
	if errBal != nil {
		return nil, errBal
	}
	newBalance := new(big.Float).Add(amount, big.NewFloat(bal))

	updatedBalance, errUpdate := store.Firestore.UpdateBtcBalance(uid, newBalance)
	if errUpdate != nil {
		return nil, errUpdate
	}

	return big.NewFloat(updatedBalance), nil
}

// ConfirmBtcTransactions confirm transactions and update corresponding balances
func ConfirmBtcTransactions(txs []*store.BtcTransactionSchema) (err error) {
	if err = store.Firestore.UpdateBtcTransactionsConfirmation(txs); err != nil {
		log.Fatal(err)
		return
	}

	for _, t := range txs {
		a, err := store.Firestore.FindBtcAccountByAddress(t.To)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if _, err = UpdateAccountBtcBalance(a.UID, big.NewFloat(t.Amount)); err != nil {
			log.Fatal(err)
			continue
		}
	}

	return
}

// ConfirmEthTransactions confirm transactions and update corresponding balances
func ConfirmEthTransactions(txs []*store.EthTransactionSchema, config *env.ChainConfig) (err error) {
	if err = store.Firestore.UpdateEthTransactionsConfirmation(txs); err != nil {
		log.Fatal(err)
		return
	}

	for _, t := range txs {
		// if the transaction if from the gas station then we do not update the balance
		if t.From == config.GasStation {
			continue
		}
		a, err := store.Firestore.FindEthAccountByAddress(t.Receiver)
		if err != nil {
			log.Print(err)
			continue
		}
		weiAmount, _ := new(big.Int).SetString(t.Amount, 10)
		curr := env.FindCurrency(t.To, config.Currencies)
		if _, err = updateAccountEthBalance(a.UID, updateEthDecimal(weiAmount, curr.Decimals), t.Currency); err != nil {
			log.Print(err)
			continue
		}
	}
	return
}

// UpdateAccountEthBalance update the eth balance of a user UID by a given amount
func updateAccountEthBalance(uid string, toAdd *big.Float, curr string) (float64, error) {
	base, errBal := store.Firestore.FindEthBalance(uid, curr)
	if errBal != nil {
		return 0, errBal
	}

	if base == nil {
		base = big.NewFloat(0)
	}
	total := new(big.Float).Add(base, toAdd)
	updatedBalance, errUpdate := store.Firestore.UpdateEthBalance(uid, total, curr)
	if errUpdate != nil {
		return 0, errUpdate
	}

	return updatedBalance, nil
}
