package cloudfunctions

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/SoteriaTech/blockchain-functions/api"
	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/env"
	"github.com/SoteriaTech/blockchain-functions/eth"
	"github.com/SoteriaTech/blockchain-functions/functions"
	"github.com/SoteriaTech/blockchain-functions/store"
	"github.com/SoteriaTech/blockchain-functions/utils"
)

var config *env.Config

// init function is ran automatically by GCP prior to the rest
func init() {
	config = env.InitConfig()

	utils.InitErrorReporting(config.ProjectID)
	store.InitFirestoreStore(config.ProjectID, config.KeyPath)

	api.InitBlockInfoClient(config.Bitcoin.Endpoint)
	btc.InitBtcService(api.BlockInfo)

	api.InitInfuraClient(config.Ethereum.Endpoint)
	eth.InitEthService(api.Infura, &config.Ethereum)
}

/***********************************************
*
* HTTP functions
*
***********************************************/

// SyncBtcBalance function sync the btc balance of a given user's account
func SyncBtcBalance(w http.ResponseWriter, r *http.Request) {

	data, errReq := utils.RequestData(r)
	if errReq != nil {
		utils.ErrorReport.LogAndPrintError(errReq)
		utils.RespondJSONWithError(w, 400, errReq.Error())
	}

	btcAccount, err := functions.SyncBtcBalance(data["uid"])
	if err != nil {
		utils.ErrorReport.LogAndPrintError(err.Err)
		utils.RespondJSONWithError(w, err.Code, err.Err.Error())
	}
	utils.RespondJSON(w, 200, btcAccount)
}

// ScanBtcBlock scan a bitcoin blockchain block and parse it
func ScanBtcBlock(w http.ResponseWriter, r *http.Request) {
	data, errReq := utils.RequestData(r)
	if errReq != nil {
		utils.RespondJSONWithError(w, 400, errReq.Error())
	}

	height, errConv := strconv.Atoi(data["height"])
	if errConv != nil {
		utils.ErrorReport.LogAndPrintError(errConv)
		utils.RespondJSONWithError(w, 400, "error height format is incorrect")
		return
	}

	rsp, err := functions.ScanBtcBlock(height, &config.Bitcoin)
	if err != nil {
		utils.ErrorReport.LogAndPrintError(err)
		utils.RespondJSONWithError(w, 500, err.Error())
	}

	utils.RespondJSON(w, 200, rsp)
}

// ScanBtcHead scan this is a replica of the pub/sub to test on the local server
func ScanBtcHead(w http.ResponseWriter, r *http.Request) {
	blocks, err := functions.ScanBtcHead(&config.Bitcoin)
	if err != nil {
		utils.RespondJSONWithError(w, err.Code, err.Err.Error())
		return
	}

	utils.RespondJSON(w, 200, blocks)
}

// ScanEthBlock scan an ethereum block for transactions
func ScanEthBlock(w http.ResponseWriter, r *http.Request) {
	data, errReq := utils.RequestData(r)
	if errReq != nil {
		utils.RespondJSONWithError(w, 400, errReq.Error())
	}

	h, errConv := strconv.Atoi(data["height"])
	if errConv != nil {
		utils.ErrorReport.LogAndPrintError(errConv)
		utils.RespondJSONWithError(w, 400, "error height format is incorrect")
		return
	}

	b, err := functions.ScanEthBlock(uint64(h), &config.Ethereum)
	if err != nil {
		utils.RespondJSONWithError(w, 400, err.Error())
		return
	}
	utils.RespondJSON(w, 200, b)
}

// ScanEthHead scan this is a replica of the pub/sub to test on the local server
func ScanEthHead(w http.ResponseWriter, r *http.Request) {
	blocks, err := functions.ScanEthHead(&config.Ethereum)
	if err != nil {
		utils.RespondJSONWithError(w, err.Code, err.Err.Error())
		return
	}

	utils.RespondJSON(w, 200, blocks)
}

/***********************************************
*
* Pub/Sub functions
*
***********************************************/

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ScanBtcPubSub ping the btc blockchain for new block and scan them for transactions
func ScanBtcPubSub(ctx context.Context, m PubSubMessage) error {
	blocks, err := functions.ScanBtcHead(&config.Bitcoin)
	if err != nil {
		utils.NotifySlack(err.Err.Error(), config.ProjectID)
		return err.Err
	}

	log.Printf("Blocks  aggregated: %v", blocks)
	return nil
}

// ScanEthPubSub ping the ethereum blockchain for new block and scan them for transactions
func ScanEthPubSub(ctx context.Context, m PubSubMessage) error {
	blocks, err := functions.ScanEthHead(&config.Ethereum)
	if err != nil {
		utils.NotifySlack(err.Err.Error(), config.ProjectID)
		return err.Err
	}

	log.Printf("Ethereum Blocks aggregated: %v", blocks)
	return nil
}

/***********************************************
*
* DO NOT CALL
* DO NOT DEPLOY
*
***********************************************/

// UpdateInterestHistories update the timestamp of the lastest history for each user by the given time (string)
func UpdateInterestHistories(w http.ResponseWriter, r *http.Request) {
	data, errReq := utils.RequestData(r)
	if errReq != nil {
		utils.RespondJSONWithError(w, 400, errReq.Error())
	}

	target, errC := strconv.Atoi(data["target"])
	if errC != nil {
		utils.RespondJSONWithError(w, 400, errC.Error())
		return
	}

	res, err := store.Firestore.FindAndUpdateLatestHistories(data["time"], target)
	if err != nil {
		utils.RespondJSONWithError(w, 500, err.Error())
		return
	}

	utils.RespondJSON(w, 200, res)

}
