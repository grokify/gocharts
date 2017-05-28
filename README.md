Go Rickshaw
===========

[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

[Rickshaw](http://code.shutterstock.com/rickshaw/) is a JavaScript toolkit for creating interactive time series graphs.

Go Rickshaw is a Go library that prepares data to be represented in Rickshaw. It uses [`quicktemplate`](https://github.com/valyala/quicktemplate) for rendering. The initial goal is to provide an easy way to format data for the Rickshaw extensions example:

![](images/graph_example_2.png)

## Usage

See the example here:

[reports/extensions_by_month/examples/example1/report.go](reports/extensions_by_month/examples/example1/report.go)

 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-rickshaw
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/go-rickshaw
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/go-rickshaw
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/go-rickshaw/blob/master/LICENSE.md