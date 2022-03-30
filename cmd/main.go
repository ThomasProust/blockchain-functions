package main

import (
	"context"
	"log"
	"os"

	functions "github.com/SoteriaTech/blockchain-functions"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	ctx := context.Background()
	funcframework.RegisterHTTPFunctionContext(ctx, "/sync_btc_balance", functions.SyncBtcBalance)
	funcframework.RegisterHTTPFunctionContext(ctx, "/scan_btc_block", functions.ScanBtcBlock)
	funcframework.RegisterHTTPFunctionContext(ctx, "/scan_btc_head", functions.ScanBtcHead)

	funcframework.RegisterHTTPFunctionContext(ctx, "/scan_eth_block", functions.ScanEthBlock)
	funcframework.RegisterHTTPFunctionContext(ctx, "/scan_eth_head", functions.ScanEthHead)

	funcframework.RegisterHTTPFunctionContext(ctx, "/update_histories", functions.UpdateInterestHistories)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
