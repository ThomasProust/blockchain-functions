package store

import "time"

//BtcAccountSchema firestore schema of a firebase bitcoin account
type BtcAccountSchema struct {
	UID     string  `firestore:"uid"`
	Address string  `firestore:"address"`
	BTC     float64 `firestore:"BTC"`
}

//EthAccountSchema firestore schema of a firebase ETH account
type EthAccountSchema struct {
	UID     string  `firestore:"uid"`
	Address string  `firestore:"address"`
	ETH     float64 `firestore:"ETH"` // because of the 18 decimals, eth balance is stored as a string, converted in WEI (10**-18 ETH)
	USDC    float64 `firestore:"USDC"`
	USDT    float64 `firestore:"USDT"`
}

// BtcTransactionSchema firestore schema of a btc transaction
type BtcTransactionSchema struct {
	Amount      float64 `firestore:"amount"`
	To          string  `firestore:"to"`
	TxHash      string  `firestore:"txHash"`
	VoutIdx     int     `firestore:"vout_idx"`
	BlockHeight int     `firestore:"block_height"`
	Confirmed   bool    `firestore:"confirmed"`
}

// ChainStateSchema firestore schema of a chain state
type ChainStateSchema struct {
	Hash        string    `firestore:"hash"`
	Time        int       `firestore:"time"`
	LastUpdated time.Time `firestore:"last_updated"`
	BlockIndex  int       `firestore:"block_index"`
	Height      int       `firestore:"height"`
	TxIndexes   []int     `firestore:"txIndexes"`
}

// EthTransactionSchema firestore schema of an ethereum transaction
type EthTransactionSchema struct {
	From        string `firestore:"from"`
	Amount      string `firestore:"amount"`
	To          string `firestore:"to"`
	TxHash      string `firestore:"txHash"`
	LogIdx      string `firestore:"log_idx"`
	BlockHeight int    `firestore:"block_height"`
	Confirmed   bool   `firestore:"confirmed"`
	Currency    string `firestore:"currency"`
	Receiver    string `firestore:"receiver"`
}
