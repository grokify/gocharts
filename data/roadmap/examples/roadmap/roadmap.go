package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/math/mathutil"

	"github.com/grokify/gocharts/v2/data/roadmap"
)

func main() {
	can, err := roadmap.GetCanvasQuarter(20174, 20182)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	item := roadmap.Item{Name: "Feature 1"}
	err = item.SetMinMaxQuarter(20174, 20174)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
	can.AddItem(item)

	item = roadmap.Item{Name: "Feature 2"}
	err = item.SetMinMaxQuarter(20174, 20181)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(3)
	}
	can.AddItem(item)

	can.Range = mathutil.RangeInt64{
		Min:   can.MinX,
		Max:   can.MaxX,
		Cells: 400,
	}

	item2 := roadmap.Item{Name: "Feature 3"}
	err = item2.SetMinMaxQuarter(20181, 20182)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(4)
	}

	can.AddItem(item2)
	err = can.InflateItems()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(5)
	}

	can.BuildRows()
	fmtutil.MustPrintJSON(can)

	fmt.Println("DONE")
}
