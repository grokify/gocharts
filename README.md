GoCharts
========

[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

GoCharts is a library to build SVG charts using various D3 libraries including C3 and [Rickshaw](https://github.com/shutterstock/rickshaw).

[`quicktemplate`](https://github.com/valyala/quicktemplate) is used for rendering.

An example chart is the Rickshaw chart shown below:

![](images/graph_example_2.png)

## Installation

```bash
$ go get github.com/grokify/gocharts/...
```

## Usage

See the example here:

[charts/rickshaw/examples/report.go](charts/rickshaw/examples/report.go)

 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gocharts
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/gocharts
 [docs-godoc-svg]: https://img.shields.io/badge/reference-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/gocharts
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/gocharts/blob/master/LICENSE.md