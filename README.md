# uuidslice

[![Build Status](https://cloud.drone.io/api/badges/pfmt/uuidslice/status.svg)](https://cloud.drone.io/pfmt/uuidslice)
[![Go Reference](https://pkg.go.dev/badge/github.com/pfmt/uuidslice.svg)](https://pkg.go.dev/github.com/pfmt/uuidslice)

UUID slice utils for Go.  
Source files are distributed under the BSD-style license.

## About

The software is considered to be at a alpha level of readiness,
its extremely slow and allocates a lots of memory.

## Benchmark

```sh
$ go test -count=1 -race -bench ./... 
goos: linux
goarch: amd64
pkg: github.com/pfmt/uuidslice
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkUniqueCopy/slice_test.go:31-8         	 1628802	       746.8 ns/op
PASS
ok  	github.com/pfmt/uuidslice	1.995s
```
