GoCharts
========

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Used By][used-by-svg]][used-by-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![SLOC][loc-svg]][repo-url]
[![License][license-svg]][license-url]

GoCharts is a library to assist with building charts, by directly working with charting libraries, generating tabular data for Excel XLSX files and CSV files, or to transfer data to/from analytics solutions like [Metabase](https://pkg.go.dev/github.com/grokify/go-metabase/metabaseutil) and [SimpleKPI](https://pkg.go.dev/github.com/grokify/go-simplekpi/simplekpiutil).

It includes two sets of packages:

1. data structures to generically hold and manipulate different types of data
1. chart library helpers to make generating charts eaiser, often times using data structures mentioned above

## Data Structures

Commonly used data structures include:

* [Table](https://pkg.go.dev/github.com/grokify/gocharts/v2/data/table) - easy manipulation of tabular data including [writing to CSV and XLSX](data/table/write.go).
* [Time Series](https://pkg.go.dev/github.com/grokify/gocharts/v2/data/timeseries) - for building time-based line charts and bar charts.
* [Histogram](https://pkg.go.dev/github.com/grokify/gocharts/v2/data/histogram) - for building histograms and bar charts.

A full list is available in the [`data`](data) folder.

## Chart Helpers

* [C3](https://c3js.org/) - [code](charts/c3)
* [D3](https://d3js.org/) - [code](charts/d3)
* [Rickshaw](https://github.com/shutterstock/rickshaw) - [code](charts/rickshaw)
* [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart) - [code](charts/wchart)

[`quicktemplate`](https://github.com/valyala/quicktemplate) is used for rendering some of the charts.

An example chart is the Rickshaw chart shown below:

![](charts/rickshaw/graph_example_2.png)

## Collections

Data collections are provided in the [`collections`](collections) folder for the primary purpose of providing example data to run in the examples. Currently, cryptocurrency data from Yahoo! Finance is included.

## Installation

```bash
$ go get github.com/grokify/gocharts/v2/...
```

## Usage

See the example here:

[charts/rickshaw/examples/report.go](charts/rickshaw/examples/report.go)

 [build-status-svg]: https://github.com/grokify/gocharts/workflows/go%20build/badge.svg
 [build-status-url]: https://github.com/grokify/gocharts/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gocharts
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gocharts
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gocharts
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gocharts/v2
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gocharts/blob/master/LICENSE
 [used-by-svg]: https://sourcegraph.com/github.com/grokify/gocharts/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/gocharts?badge
 [loc-svg]: https://tokei.rs/b1/github/grokify/gocharts
 [repo-url]: https://github.com/grokify/gocharts