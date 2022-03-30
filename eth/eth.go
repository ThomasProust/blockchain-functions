package eth

import (
	"math/big"

	ethinfura "github.com/INFURA/go-ethlibs/eth"
	"github.com/SoteriaTech/blockchain-functions/env"
	"github.com/SoteriaTech/blockchain-functions/utils"
)

// EthereumAPI interface that Eth service implements
type EthereumAPI interface {
	GetBlockHeader() (uint64, error)
	GetTransactionsFromBlock(h *big.Int) (*BlockData, error)
	GetTransactionByHash(hash string, b int) (*Transaction, error)
	GetReceipt(hash string) (*ethinfura.TransactionReceipt, error)
}

// Eth stucture of the Eth service
type Eth struct {
	api    EthereumAPI
	config *env.ChainConfig
}

// EthService instance of the EthService
var ethService *Eth

// InitEthService initialize the instance of EthService
func InitEthService(api EthereumAPI, config *env.ChainConfig) {
	ethService = &Eth{
		api:    api,
		config: config,
	}
}

// GetHeadBlock get the head block number of the chain
func GetHeadBlock() (uint64, error) {
	return ethService.api.GetBlockHeader()
}

// ScanBlock scan a block to retrieve its transactions
func ScanBlock(h uint64) (*BlockData, error) {

	bd, errT := ethService.api.GetTransactionsFromBlock(new(big.Int).SetUint64(h))
	if errT != nil {
		return nil, errT
	}
	for idx, tx := range bd.Txs {
		currency := env.FindCurrency(tx.To, ethService.config.Currencies)
		tx.Currency = currency.Name

		if currency.Address != "" {
			tf, err := parseTokenTransfer(tx, currency)
			if err != nil {
				utils.ErrorReport.LogAndPrintError(err)
				continue
			}
			bd.Txs[idx] = tf
		}
	}
	return bd, nil
}

// ConfirmTransactions ask the blockchain for confirmed transactions
func ConfirmTransactions(hashes []string, b int) (confirmed []string, errs []error) {
	for _, h := range hashes {
		tx, err := ethService.api.GetTransactionByHash(h, b)

		if err != nil {
			errs = append(errs, err)
			continue
		}
		confirmed = append(confirmed, tx.Hash)
	}
	return
}
