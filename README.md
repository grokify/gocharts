# GoCharts

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

GoCharts is a library to assist with building charts, by directly working with charting libraries, generating tabular data for Excel XLSX files and CSV files, or to transfer data to/from analytics solutions like [Metabase](https://pkg.go.dev/github.com/grokify/go-metabase/metabaseutil) and [SimpleKPI](https://pkg.go.dev/github.com/grokify/go-simplekpi/simplekpiutil).

## Features

- **Data Structures** - Table, TimeSeries, Histogram, Roadmap, and more
- **Multiple Chart Libraries** - C3, D3, ECharts, Google Charts, Rickshaw, wchart
- **Excel Integration** - Read/write XLSX files with formatting
- **Markdown Output** - Generate tables for documentation
- **Analytics Integration** - Metabase and SimpleKPI support
- **Text Charts** - Terminal-friendly progress bars and funnel charts

## Contents

It includes two sets of packages:

1. data structures to generically hold and manipulate different types of data
1. chart library helpers to make generating charts easier, often times using data structures mentioned above

### Data Structures

Commonly used data structures include:

* [Table](https://pkg.go.dev/github.com/grokify/gocharts/v2/data/table) - easy manipulation of tabular data including [writing to CSV and XLSX](data/table/write.go).
* [Time Series](https://pkg.go.dev/github.com/grokify/gocharts/v2/data/timeseries) - for building time-based line charts and bar charts.
* [Histogram](https://pkg.go.dev/github.com/grokify/gocharts/v2/data/histogram) - for building histograms and bar charts.

A full list is available in the [`data`](data) folder.

### Chart Helpers

* [C3](https://pkg.go.dev/github.com/grokify/gocharts/v2/charts/c3) - [code](charts/c3), [project](https://c3js.org/)
* [D3](https://pkg.go.dev/github.com/grokify/gocharts/v2/charts/d3) - [code](charts/d3), [project](https://d3js.org/)
* [ECharts](https://pkg.go.dev/github.com/grokify/gocharts/v2/charts/echarts) - [code](charts/echarts), [project](https://echarts.apache.org/)
* [Google Charts](https://pkg.go.dev/github.com/grokify/gocharts/v2/charts/google) - [code](charts/google), [project](https://developers.google.com/chart/interactive/docs)
* [Rickshaw](https://pkg.go.dev/github.com/grokify/gocharts/v2/charts/rickshaw) - [code](charts/rickshaw), [project](https://github.com/shutterstock/rickshaw)
* [go-analyze/charts](https://pkg.go.dev/github.com/grokify/gocharts/v2/charts/wchart) - [code](charts/wchart), [project](https://github.com/go-analyze/charts)
* [Text Charts](https://pkg.go.dev/github.com/grokify/gocharts/v2/charts/text) - [code](charts/text) - Text-based progress bars and funnel charts

[`quicktemplate`](https://github.com/valyala/quicktemplate) is used for rendering some of the charts.

An example chart is the Rickshaw chart shown below:

![](charts/rickshaw/graph_example_2.png)

### Collections

Data collections are provided in the [`collections`](collections) folder for the primary purpose of providing example data to run in the examples. Currently, cryptocurrency data from Yahoo! Finance is included.

### Applications

Various helpers to use applications are located in the [`apps`](apps) folder for the primary purpose of providing reusable and example code.

## Installation

```bash
$ go get github.com/grokify/gocharts/v2/...
```

## Quick Start

### Creating a Histogram

```go
import "github.com/grokify/gocharts/v2/data/histogram"

h := histogram.NewHistogram("Response Codes")
h.Add("200", 150)
h.Add("404", 25)
h.Add("500", 10)

// Output as Markdown table
md := h.Markdown()
```

### Creating a Table and Exporting to XLSX

```go
import "github.com/grokify/gocharts/v2/data/table"

tbl := table.NewTable("Sales Report")
tbl.Columns = []string{"Region", "Q1", "Q2", "Q3", "Q4"}
tbl.Rows = [][]string{
    {"North", "100", "120", "130", "150"},
    {"South", "90", "95", "100", "110"},
}

// Write to Excel
err := tbl.WriteXLSX("report.xlsx")
```

### Converting Histogram to Google Charts DataTable

```go
import (
    "github.com/grokify/gocharts/v2/data/histogram"
    "github.com/grokify/gocharts/v2/charts/google"
)

hs := histogram.NewHistogramSet("Traffic by Hour")
// ... populate histogram set ...

dt := google.DataTableFromHistogramSet(hs)
```

## Output Formats

GoCharts supports multiple output formats:

| Format | Package | Description |
|--------|---------|-------------|
| CSV | `data/table` | Comma-separated values |
| XLSX | `data/table` | Excel spreadsheets via [excelize](https://github.com/xuri/excelize) |
| Markdown | `data/histogram` | GitHub-flavored markdown tables |
| HTML | `charts/*` | Chart library-specific HTML via [quicktemplate](https://github.com/valyala/quicktemplate) |
| ASCII | `data/table` | Terminal-friendly tables via [tablewriter](https://github.com/olekukonko/tablewriter) |

## Examples

Examples are available in each chart package's `examples/` directory:

| Chart Library | Example Location |
|---------------|------------------|
| Google Charts | [charts/google/examples/](charts/google/examples/) |
| wchart | [charts/wchart/examples/](charts/wchart/examples/) |
| C3 | [charts/c3/examples/](charts/c3/examples/) |
| Rickshaw | [charts/rickshaw/examples/](charts/rickshaw/examples/) |
| D3 Bullet | [charts/d3/d3bullet/examples/](charts/d3/d3bullet/examples/) |

 [build-status-svg]: https://github.com/grokify/gocharts/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/gocharts/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/gocharts/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/gocharts/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gocharts
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gocharts
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gocharts
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gocharts/v2
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gocharts/blob/main/LICENSE
 [used-by-svg]: https://sourcegraph.com/github.com/grokify/gocharts/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/gocharts?badge
 [loc-svg]: https://tokei.rs/b1/github/grokify/gocharts
 [repo-url]: https://github.com/grokify/gocharts

 ## Mentions

 1. [Philip Gardner's GitHub stars: `github.com/gaahrdner/starred`](https://github.com/gaahrdner/starred)
