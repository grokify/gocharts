GoCharts
========

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

GoCharts is a library to assist with building charts, by directly working with charting libraries, generating tabular data for Excel XLSX files and CSV files, or to transfer data to/from analytics solutions like [Metabase](https://pkg.go.dev/github.com/grokify/go-metabase/metabaseutil) and [SimpleKPI](https://pkg.go.dev/github.com/grokify/go-simplekpi/simplekpiutil).

It includes two sets of packages:

1. data structures to generically hold and manipulate different types of data
1. chart library helpers to make generating charts eaiser, often times using data structures mentioned above

## Data Structures

Commonly used data structures include:

* [Table](https://pkg.go.dev/github.com/grokify/gocharts/data/table) - Easy manipulation including [writing to CSV and XLSX](data/table/write.go).
* [Time Series](https://pkg.go.dev/github.com/grokify/gocharts/data/timeseries) - for building time-based line charts and bar charts.
* [Histogram](https://pkg.go.dev/github.com/grokify/gocharts/data/histogram) - for building histograms and bar charts.

A full list is available in the [`data`](data) folder.

## Chart Helpers

* [C3](https://c3js.org/) - [code](charts/c3)
* [D3](https://d3js.org/) - [code](charts/d3)
* [Rickshaw](https://github.com/shutterstock/rickshaw) - [code](charts/rickshaw)
* [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart) - [code](charts/wchart)

[`quicktemplate`](https://github.com/valyala/quicktemplate) is used for rendering some of the charts.

An example chart is the Rickshaw chart shown below:

![](charts/rickshaw/graph_example_2.png)

## Installation

```bash
$ go get github.com/grokify/gocharts/...
```

## Usage

See the example here:

[charts/rickshaw/examples/report.go](charts/rickshaw/examples/report.go)

 [build-status-svg]: https://github.com/grokify/gocharts/workflows/go%20build/badge.svg
 [build-status-url]: https://github.com/grokify/gocharts/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gocharts
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gocharts
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gocharts
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gocharts
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gocharts/blob/master/LICENSE