package store

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FireStoreStore struct for firestore DB
type FireStoreStore struct {
	Client *firestore.Client
	ctx    context.Context
}

// Firestore instance of Firestore store
var Firestore *FireStoreStore

//InitFirestoreStore initialize a new firestore client
func InitFirestoreStore(projectID string, keyPath string) {
	ctx := context.Background()
	client := newFireStoreClient(ctx, projectID, keyPath)

	Firestore = &FireStoreStore{
		Client: client,
		ctx:    ctx,
	}
}

func newFireStoreClient(ctx context.Context, projectID string, keyPath string) *firestore.Client {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create firestore client %v", err)
	}
	return client
}

// FindBtcAccount find btc account from a user UID
func (f *FireStoreStore) FindBtcAccount(uid string) (*BtcAccountSchema, error) {
	var btcAccount *BtcAccountSchema

	doc, err := f.Client.Collection("btc_accounts").Doc(uid).Get(f.ctx)
	if err != nil {
		return nil, err
	}
	if err := doc.DataTo(&btcAccount); err != nil {
		return nil, err
	}
	btcAccount.UID = uid

	return btcAccount, nil
}

// FindBtcBalance find the btc balance of a user UID
func (f *FireStoreStore) FindBtcBalance(uid string) (float64, error) {
	data := make(map[string]float64)

	doc, err := f.Client.Collection("balances").Doc(uid).Get(f.ctx)
	if err != nil {
		return 0, err
	}

	if err := doc.DataTo(&data); err != nil {
		return 0, err
	}

	return data["BTC"], nil
}

// FindEthBalance find the eth balance of a user UID
func (f *FireStoreStore) FindEthBalance(uid string, curr string) (*big.Float, error) {
	data := &EthAccountSchema{}

	doc, err := f.Client.Collection("balances").Doc(uid).Get(f.ctx)
	if err != nil {
		return nil, err
	}

	if err := doc.DataTo(&data); err != nil {
		log.Fatalf("error here %v", err)
		return nil, err
	}
	// TODO: that's bad, need to refactor
	val := data.ETH
	if curr == "USDC" {
		val = data.USDC
	}

	if curr == "USDT" {
		val = data.USDT
	}

	return big.NewFloat(val), nil
}

// UpdateBtcBalance update the btc balance of a user's account
func (f *FireStoreStore) UpdateBtcBalance(uid string, newBalance *big.Float) (float64, error) {
	doc := make(map[string]interface{})
	flBalance, _ := newBalance.Float64()
	doc["BTC"] = flBalance

	_, err := f.Client.Collection("balances").Doc(uid).Set(f.ctx, doc, firestore.MergeAll)
	if err != nil {
		return flBalance, err
	}
	return flBalance, nil
}

// UpdateEthBalance update the eth balance of a user's account
func (f *FireStoreStore) UpdateEthBalance(uid string, bal *big.Float, curr string) (float64, error) {
	doc := make(map[string]interface{})
	flBalance, _ := bal.Float64()
	doc[curr] = flBalance
	_, err := f.Client.Collection("balances").Doc(uid).Set(f.ctx, doc, firestore.MergeAll)
	if err != nil {
		return flBalance, err
	}
	return flBalance, nil
}

// GetAllBtcAccountAddresses get all the current bitcoin accounts and addresses from Soteria
func (f *FireStoreStore) GetAllBtcAccountAddresses() ([]*BtcAccountSchema, error) {
	var accs []*BtcAccountSchema
	iter := f.Client.Collection("btc_accounts").Documents(f.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var acc *BtcAccountSchema
		doc.DataTo(&acc)
		acc.UID = doc.Ref.ID

		accs = append(accs, acc)
	}
	return accs, nil
}

// GetAllEthAccountAddresses get all the current ethereum accounts and addresses from Soteria
func (f *FireStoreStore) GetAllEthAccountAddresses() ([]*EthAccountSchema, error) {
	var accs []*EthAccountSchema
	iter := f.Client.Collection("eth_accounts").Documents(f.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		var acc *EthAccountSchema
		doc.DataTo(&acc)
		acc.UID = doc.Ref.ID

		accs = append(accs, acc)
	}
	return accs, nil
}

// GetConvertRequests get all convert requests of a given account
func (f *FireStoreStore) GetConvertRequests(uid string) (map[string]interface{}, error) {
	docs := make(map[string]interface{})
	iter := f.Client.Collection("convert_history").Doc(uid).Collection("history").Documents(f.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		docs[uid] = doc.Data()
	}

	return docs, nil
}

// FindBtcTransaction find a btc transaction by hash
func (f *FireStoreStore) FindBtcTransaction(idx string) (t *BtcTransactionSchema, err error) {
	doc, errStore := f.Client.Collection("btc_transactions").Doc(idx).Get(f.ctx)

	if errStore != nil && status.Code(errStore) != codes.NotFound {
		err = errStore
		return
	}

	if doc.Exists() {
		errDoc := doc.DataTo(&t)
		if errDoc != nil {
			err = errDoc
		}
	}
	return
}

// FindEthTransaction find an ethereum transaction by hash
func (f *FireStoreStore) FindEthTransaction(idx string) (t *EthTransactionSchema, err error) {
	doc, errStore := f.Client.Collection("eth_transactions").Doc(idx).Get(f.ctx)

	if errStore != nil && status.Code(errStore) != codes.NotFound {
		err = errStore
		return
	}

	if doc.Exists() {
		errDoc := doc.DataTo(&t)
		if errDoc != nil {
			err = errDoc
		}
	}
	return
}

// CreateBtcTransaction create a btc transaction
func (f *FireStoreStore) CreateBtcTransaction(t *BtcTransactionSchema) (err error) {
	_, err = f.Client.Collection("btc_transactions").Doc(t.TxHash+strconv.Itoa(t.VoutIdx)).Create(f.ctx, &t)
	return
}

// CreateEthTransaction create an eth transactions
func (f *FireStoreStore) CreateEthTransaction(t *EthTransactionSchema) (err error) {
	_, err = f.Client.Collection("eth_transactions").Doc(t.TxHash+t.LogIdx).Create(f.ctx, &t)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// GetChainState get the latest block data of the given chain from the store
func (f *FireStoreStore) GetChainState(chain string) (map[string]interface{}, error) {
	hs := make(map[string]interface{})
	doc, err := f.Client.Collection("chain_state").Doc(chain).Get(f.ctx)
	if err != nil {
		return nil, err
	}
	errData := doc.DataTo(&hs)
	if errData != nil {
		return nil, errData
	}
	return hs, err
}

// UpdateChainState update the latest block data of the given chain
func (f *FireStoreStore) UpdateChainState(chain string, data map[string]interface{}) (err error) {

	_, err = f.Client.Collection("chain_state").Doc(chain).Set(f.ctx, data)
	return
}

// FindBtcTransactionsFromBlockHeight find transactions that have been recorded from a specific block height
func (f *FireStoreStore) FindBtcTransactionsFromBlockHeight(h int) (txs []*BtcTransactionSchema, err error) {
	iter := f.Client.Collection("btc_transactions").Where("block_height", "==", h).Where("confirmed", "==", false).Documents(f.ctx)
	for {
		doc, errIter := iter.Next()
		if errIter == iterator.Done {
			break
		}
		if errIter != nil {
			err = errIter
			return
		}
		var tx *BtcTransactionSchema
		if err = doc.DataTo(&tx); err != nil {
			return
		}
		txs = append(txs, tx)
	}

	return
}

// FindEthTransactionsFromBlockHeight find transactions that have been recorded from a specific block height
func (f *FireStoreStore) FindEthTransactionsFromBlockHeight(h int) ([]*EthTransactionSchema, error) {
	var txs []*EthTransactionSchema
	iter := f.Client.Collection("eth_transactions").Where("block_height", "==", h).Where("confirmed", "==", false).Documents(f.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			return nil, err
		}
		var tx *EthTransactionSchema
		if err = doc.DataTo(&tx); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

// UpdateBtcTransactionsConfirmation update confirmation for each given transaction
func (f *FireStoreStore) UpdateBtcTransactionsConfirmation(txs []*BtcTransactionSchema) (err error) {
	for _, t := range txs {
		uid := t.TxHash
		if t.VoutIdx >= 0 {
			uid = uid + strconv.Itoa(t.VoutIdx)
		}
		_, errSet := f.Client.Collection("btc_transactions").Doc(uid).Set(f.ctx, BtcTransactionSchema{Confirmed: true}, firestore.Merge([]string{"confirmed"}))
		if err != nil {
			err = errSet
			continue
		}
	}

	return
}

// UpdateEthTransactionsConfirmation update confirmation for each given transaction
func (f *FireStoreStore) UpdateEthTransactionsConfirmation(txs []*EthTransactionSchema) (err error) {
	for _, t := range txs {
		uid := t.TxHash
		if t.LogIdx != "" {
			uid = uid + t.LogIdx
		}
		_, errSet := f.Client.Collection("eth_transactions").Doc(uid).Set(f.ctx, BtcTransactionSchema{Confirmed: true}, firestore.Merge([]string{"confirmed"}))
		if err != nil {
			err = errSet
			continue
		}
	}

	return
}

// FindBtcAccountByAddress find a firestore bitcoin account from an address
func (f *FireStoreStore) FindBtcAccountByAddress(addr string) (a *BtcAccountSchema, err error) {
	doc, errQ := f.Client.Collection("btc_accounts").Where("address", "==", addr).Documents(f.ctx).Next()
	if errQ != nil {
		err = errQ
		return
	}
	if doc.Exists() {
		doc.DataTo(&a)
		a.UID = doc.Ref.ID
	}
	return
}

// FindEthAccountByAddress find a firestore bitcoin account from an address
func (f *FireStoreStore) FindEthAccountByAddress(addr string) (a *EthAccountSchema, err error) {
	doc, errQ := f.Client.Collection("eth_accounts").Where("address", "==", addr).Documents(f.ctx).Next()
	if errQ != nil {
		err = errQ
		return
	}
	if doc.Exists() {
		doc.DataTo(&a)
		a.UID = doc.Ref.ID
	}
	return
}

// FindAndUpdateLatestHistories find latest histories
func (f *FireStoreStore) FindAndUpdateLatestHistories(timestamp string, target int) (res []map[string]interface{}, err error) {
	iter := f.Client.Collection("interest_payment_histories").Documents(f.ctx)
	for {
		doc, errIter := iter.Next()
		if errIter == iterator.Done {
			break
		}
		if errIter != nil {
			err = errIter
			return
		}

		hIter := doc.Ref.Collection("histories").Where("timestamp", "==", target).Documents(f.ctx)
		for {
			hDoc, errD := hIter.Next()
			if errD == iterator.Done {
				break
			}
			if errD != nil {
				err = errD
				return
			}
			hist := hDoc.Data()
			ts, errT := time.Parse("2006-01-02 15:04:05", timestamp)
			if errT != nil {
				err = errT
				return
			}
			hist["timestamp"] = ts.Unix()
			hist["datetime"] = ts.Format("2006-01-02 15:04:05")
			_, errS := hDoc.Ref.Set(f.ctx, hist)
			if errS != nil {
				err = errS
				return
			}
			res = append(res, hist)
		}
	}
	return
}
