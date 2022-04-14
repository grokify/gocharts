package main

import (
	"fmt"

	"github.com/grokify/gocharts/v2/collections"
	"github.com/grokify/gocharts/v2/collections/cryptocurrency"
	"github.com/grokify/mogo/log/logutil"
)

func main() {
	hdBTC := cryptocurrency.HistoricalDataBTCUSDMonthly()
	err := collections.WriteFilesHistoricalData("data_btc-usd", hdBTC, true)
	logutil.FatalErr(err)

	hdETH := cryptocurrency.HistoricalDataETHUSDMonthly()
	err2 := collections.WriteFilesHistoricalData("data_eth-usd", hdETH, true)
	logutil.FatalErr(err2)

	fmt.Println("DONE")
}
