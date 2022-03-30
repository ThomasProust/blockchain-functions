package functions

import (
	"encoding/json"

	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/env"
	"github.com/SoteriaTech/blockchain-functions/helpers"
	"github.com/SoteriaTech/blockchain-functions/store"
	"github.com/SoteriaTech/blockchain-functions/utils"
)

// ScanBtcHead scan the head block of the btc blockchain for transactions
// also catches on missing blocks between two pings
func ScanBtcHead(config *env.ChainConfig) ([]int, *utils.ErrorService) {
	state, err := store.Firestore.GetChainState(config.Chain)
	if err != nil {
		return nil, &utils.ErrorService{Code: 500, Err: err}
	}

	jsonState, _ := json.Marshal(&state)
	cs := &btc.HeadBlock{}
	json.Unmarshal(jsonState, cs)

	headBlock, err := btc.BtcService.GetHeadInfo()
	if err != nil {
		return nil, &utils.ErrorService{Code: 500, Err: err}
	}

	if cs.Height == headBlock.Height {
		return nil, nil
	}

	currHeight := cs.Height
	var blocks []int
	for {
		currHeight++

		_, errScan := ScanBtcBlock(currHeight, config)
		if errScan != nil {
			utils.ErrorReport.LogAndPrintError(errScan)
			return nil, &utils.ErrorService{Code: 500, Err: errScan}
		}
		blocks = append(blocks, currHeight)
		if currHeight == headBlock.Height {
			break
		}
	}

	errUpdate := store.Firestore.UpdateChainState(config.Chain, helpers.FormatBtcChainState(headBlock))
	if errUpdate != nil {
		return nil, &utils.ErrorService{Code: 500, Err: errUpdate}

	}

	return blocks, nil
}
