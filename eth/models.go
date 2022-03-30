package eth

import (
	"math/big"
	"time"
)

// BlockData structure that contains the header of a block and its transactions
type BlockData struct {
	Txs  []*Transaction
	Meta Header
}

// Header structure with basic infos of a block
type Header struct {
	Hash        string    `json:"hash"`
	Height      int       `json:"height"`
	LastUpdated time.Time `json:"last_updated"`
	Time        int       `json:"time"`
	Nonce       string    `json:"nonce"`
}

// Transaction structure of an ethereum transaction
type Transaction struct {
	From        string
	To          string
	Hash        string
	Value       *big.Int
	BlockHeight int
	Currency    string
	LogIdx      string
	Receiver    string
}
