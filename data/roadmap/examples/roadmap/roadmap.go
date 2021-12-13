package main

import (
	"fmt"

	rm "github.com/grokify/gocharts/data/roadmap"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/math/mathutil"
)

func main() {
	can, err := rm.GetCanvasQuarter(int32(20174), int32(20182))
	if err != nil {
		panic(err)
	}

	item := rm.Item{Name: "Feature 1"}
	err = item.SetMinMaxQuarter(int32(20174), int32(20174))
	if err != nil {
		panic(err)
	}
	can.AddItem(item)

	item = rm.Item{Name: "Feature 2"}
	err = item.SetMinMaxQuarter(int32(20174), int32(20181))
	if err != nil {
		panic(err)
	}
	can.AddItem(item)

	can.Range = mathutil.RangeInt64{
		Min:   can.MinX,
		Max:   can.MaxX,
		Cells: 400,
	}

	item2 := rm.Item{Name: "Feature 3"}
	err = item2.SetMinMaxQuarter(int32(20181), int32(20182))
	if err != nil {
		panic(err)
	}

	can.AddItem(item2)
	can.InflateItems()

	can.BuildRows()
	fmtutil.PrintJSON(can)

	fmt.Println("DONE")
}
