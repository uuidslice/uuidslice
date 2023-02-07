# uuidslice

[![Build Status](https://cloud.drone.io/api/badges/pfmt/uuidslice/status.svg)](https://cloud.drone.io/pfmt/uuidslice)
[![Go Reference](https://pkg.go.dev/badge/github.com/pfmt/uuidslice.svg)](https://pkg.go.dev/github.com/pfmt/uuidslice)

Copying of parsed/unique UUIDs to slice for Go.  
Source files are distributed under the BSD-style license.

## About

The software is considered to be at a alpha level of readiness,
its extremely slow and allocates a lots of memory.

## Benchmark

```sh
$ go test -count=1 -race -bench=. -benchmem ./...
goos: linux
goarch: amd64
pkg: github.com/pfmt/uuidslice
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkStringCopy/slice_test.go:34-8         	 2978049	       422.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkUniqueCopy/slice_test.go:168-8        	 1729981	       683.7 ns/op	      57 B/op	       0 allocs/op
BenchmarkExceptCopy/slice_test.go:316-8        	 1922455	       627.1 ns/op	      56 B/op	       0 allocs/op
PASS
ok  	github.com/pfmt/uuidslice	5.419s
```
