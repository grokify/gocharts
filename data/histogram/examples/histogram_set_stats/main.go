package main

import (
	"encoding/json"
	"log"

	"github.com/grokify/mogo/fmt/fmtutil"

	"github.com/grokify/gocharts/v2/data/histogram"
)

func main() {
	data := `{
	"Foo":{
	   "One":100,
	   "Two":2000,
	   "Three":300000
	},
	"Bar":{
	   "One":100,
	   "Two":2000,
	   "Three":300000
	},
	"Baz":{
	   "One":100,
	   "Two":2000,
	   "Three":300000
	},
	"Qux":{
	   "One":100,
	   "Three":300000
	},
	"Quux":{
	   "One":100,
	   "Two":2000,
	   "Three":300000
	},
	"Quuz":{
	   "One":100,
	   "Two":2000,
	   "Three":300000
	}
 }`
	msmsi := map[string]map[string]int{}

	err := json.Unmarshal([]byte(data), &msmsi)
	if err != nil {
		log.Fatal(err)
	}

	fs := histogram.NewHistogramSetWithData("FooBar", msmsi)
	stats := fs.LeafStats("CategoryNumber")
	stats.Inflate()
	fmtutil.MustPrintJSON(stats)
}
