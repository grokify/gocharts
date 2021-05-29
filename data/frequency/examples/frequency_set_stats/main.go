package main

import (
	"encoding/json"
	"log"

	"github.com/grokify/gocharts/data/frequency"
	"github.com/grokify/simplego/fmt/fmtutil"
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

	fs := frequency.NewHistogramSetWithData("FooBar", msmsi)
	stats := fs.LeafStats("CategoryNumber")
	stats.Inflate()
	fmtutil.PrintJSON(stats)
}
