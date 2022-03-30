package functions

import (
	"encoding/json"

	"github.com/SoteriaTech/blockchain-functions/env"
	"github.com/SoteriaTech/blockchain-functions/eth"
	"github.com/SoteriaTech/blockchain-functions/helpers"
	"github.com/SoteriaTech/blockchain-functions/store"
	"github.com/SoteriaTech/blockchain-functions/utils"
)

// ScanEthHead scan the head block of the ethereum blockchain for transactions
// also catches on missing blocks between two pings
func ScanEthHead(config *env.ChainConfig) ([]int, *utils.ErrorService) {

	state, err := store.Firestore.GetChainState(config.Chain)
	if err != nil {
		return nil, &utils.ErrorService{Code: 500, Err: err}
	}

	cs := &eth.Header{}
	jsonState, _ := json.Marshal(&state)
	json.Unmarshal(jsonState, &cs)

	headBlock, err := eth.GetHeadBlock()
	if err != nil {
		return nil, &utils.ErrorService{Code: 500, Err: err}
	}

	if cs.Height == int(headBlock) {
		return nil, nil
	}

	currHeight := cs.Height
	var blocks []int
	var newCS *eth.Header
	for {
		currHeight++

		bd, errScan := ScanEthBlock(uint64(currHeight), config)
		if errScan != nil {
			utils.ErrorReport.LogAndPrintError(errScan)
			return nil, &utils.ErrorService{Code: 500, Err: errScan}
		}
		blocks = append(blocks, currHeight)
		newCS = bd
		if currHeight == int(headBlock) {
			break
		}
	}
	errUpdate := store.Firestore.UpdateChainState(config.Chain, helpers.FormatEthChainState(newCS))
	if errUpdate != nil {
		return nil, &utils.ErrorService{Code: 500, Err: errUpdate}

	}

	return blocks, nil
}
