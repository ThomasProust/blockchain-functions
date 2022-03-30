package helpers

import (
	"math"
	"math/big"
	"time"

	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/eth"
	"github.com/SoteriaTech/blockchain-functions/store"
)

// FormatBtcChainState format the block header for database persistence
func FormatBtcChainState(head *btc.HeadBlock) map[string]interface{} {
	state := make(map[string]interface{})
	state["height"] = head.Height
	state["time"] = head.Time
	state["last_updated"] = time.Now()
	state["block_index"] = head.BlockIndex
	state["tx_indexes"] = head.TxIndexes
	return state
}

// FormatEthChainState format the block header for database persistence
func FormatEthChainState(head *eth.Header) map[string]interface{} {
	state := make(map[string]interface{})
	state["height"] = head.Height
	state["hash"] = head.Hash
	state["time"] = head.Time
	state["last_updated"] = time.Now()
	state["nonce"] = head.Nonce
	return state
}

// SetCurrencyAmount set the amount of a specific eth currency for a given account
func SetCurrencyAmount(acc *store.EthAccountSchema, t *store.EthTransactionSchema) *store.EthAccountSchema {

	v, _ := new(big.Float).SetString(t.Amount)
	a, _ := v.Float64()
	switch t.Currency {
	case "ETH":
		acc.ETH = a
	case "USDC":
		acc.USDC = a
	case "USDT":
		acc.USDT = a
	default:
		acc.ETH = a
	}
	return acc
}

func updateEthDecimal(i *big.Int, dec int) *big.Float {
	return intToFloat(i, dec)
}

// FromSatoshiToBtc convert a value in satoshi (int) to a value in btc (float)
func FromSatoshiToBtc(i *big.Int) (f float64) {
	f, _ = intToFloat(i, 9).Float64()
	return
}

// FromBtcToSatoshi convert a value in btc (float)  to a value in satoshi (int)
func FromBtcToSatoshi(f *big.Float) *big.Int {
	return floatToInt(f, 9)
}

func floatToInt(v *big.Float, d int) *big.Int {
	i, _ := v.Mul(v, big.NewFloat(math.Pow10(d))).Int64()
	return big.NewInt(i)
}

func intToFloat(v *big.Int, d int) *big.Float {
	flVal := new(big.Float).SetInt(v)
	return flVal.Mul(flVal, big.NewFloat(math.Pow10(-d)))
}
