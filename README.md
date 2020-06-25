[![GoDoc](https://godoc.org/github.com/LdDl/ch?status.svg)](https://godoc.org/github.com/LdDl/ch)
[![Build Status](https://travis-ci.com/LdDl/ch.svg?branch=master)](https://travis-ci.com/LdDl/ch)
[![Sourcegraph](https://sourcegraph.com/github.com/LdDl/ch/-/badge.svg)](https://sourcegraph.com/github.com/LdDl/ch?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/LdDl/ch)](https://goreportcard.com/report/github.com/LdDl/ch)
[![GitHub tag](https://img.shields.io/github/tag/LdDl/ch.svg)](https://github.com/LdDl/ch/releases)

# ch - Contraction Hierarchies
## Contraction Hierarchies (with bidirectional version of Dijkstra's algorithm) for computing shortest path in graph.

This library provides classic implementation of Dijkstra's algorithm and maneuver restrictions extension also.

## Table of Contents

- [About](#about)
- [Installation](#installation)
    - [Go get](#go-get)
    - [Go mod](#go-mod)
- [Usage](#usage)
- [Benchmark](#benchmark)
- [Support](#support)
- [ToDo](#todo)
- [Thanks](#thanks)
- [Theory](#theory)
- [Dependencies](#dependencies)
- [License](#license)

## About
This package provides implemented next techniques and algorithms:
* Dijkstra's algorithm
* Contracntion hierarchies
* Bidirectional extension of Dijkstra's algorithm with contracted nodes

## Installation

### Go get
```go
go get github.com/LdDl/ch
```


### Go mod 
In your project folder execute next command (assuming you have GO111MODULE=on):
```go
go mod init mod
```
Then import library into your code:
```go
package main

import "github.com/LdDl/ch"

func main() {
	x := ch.Graph{}
	_ = x
}
```
and build
```go
go build
```
You will see next output:
```shell
go: finding github.com/LdDl/ch v1.3.4
go: downloading github.com/LdDl/ch v1.3.4
```
And then you are good to go 

## Usage

Please see this [test file](bidirectional_ch_test.go#L17)

I hope it's pretty clear, but here is little explanation:
```go
g := Graph{} // Prepare variable for storing graph
graphFromCSV(&g, "data/pgrouting_osm.csv") // Import CSV-file file into programm
g.PrepareContracts() // Compute contraction hierarchies
u := 144031 // Define source vertex
v := 452090 // Define target vertex
ans, path := g.ShortestPath(u, v) // Get shortest path and it's cost between source and target vertex
```

### If you want to import OSM (Open Street Map) file then follow instructions for [osm2ch](https://github.com/LdDl/ch/tree/master/cmd/osm2ch#osm2ch)

## Benchmark

You can check benchmarks [here](https://github.com/LdDl/ch/blob/master/BENCHMARK.md)

## Support
If you have troubles or questions please [open an issue](https://github.com/LdDl/ch/issues/new).

## ToDo

* ~~Import file of specific format~~ **Done as CSV**
* ~~Export file of specific format~~ **Done as CSV**
* ~~Turn Restricted Shortest Path extension for CH-algorithm~~ **Propably not modify algorithm, but graph. UPD: done with osm2ch**
* ~~Thoughts and discussions about OSM graph and extensions~~ **Need some ideas about parsing and preparing. UPD: osm2ch**
* ~~Map matcher as another project **WIP (now it is in local git storage)**~~ Here is the one https://github.com/LdDl/horizon
* Bring interfaces{} **Thoughts (do we really need it?)**
* ~~Bring OSM parser~~ **WIP It's done in poor way. PRs are welcome**
* ~~Bring OSM restrictions~~ **WIP It's done in poor way. PRs are welcome**
* ~~OneTwoMany function (contraction hierarchies)~~ **Done, ~~may be some bench comparisons~~**
* ManyToMany function (contraction hierarchies) **Thoughts**
* ~~Replace int with int64 (OSM purposes)~~ **Done**
* ~~Separate benchmarks to BENCHMARK.md~~ **Done**

## Theory
[Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)

[Bidirectional search](https://en.wikipedia.org/wiki/Bidirectional_search)

[Bidirectional Dijkstra's algorithm's stop condition](http://www.cs.princeton.edu/courses/archive/spr06/cos423/Handouts/EPP%20shortest%20path%20algorithms.pdf)

[Contraction hierarchies](https://en.wikipedia.org/wiki/Contraction_hierarchies)

[Video Lectures](https://ad-wiki.informatik.uni-freiburg.de/teaching/EfficientRoutePlanningSS2012)


## Thanks
Thanks to [this](https://github.com/navjindervirdee/Advanced-Shortest-Paths-Algorithms) Java implementation of mentioned algorithms

## Dependencies
Thanks to [paulmach](https://github.com/paulmach) for his [OSM-parser](https://github.com/paulmach/osm) written in Go.

Paulmach's license is [here](https://github.com/paulmach/osm/blob/master/LICENSE.md) (it's MIT)

## License
You can check it [here](https://github.com/LdDl/ch/blob/master/LICENSE)

[osm2ch]: (https://github.com/LdDl/ch/tree/master/cmd/osm2ch#osm2ch)
[open an issue]: (https://github.com/LdDl/ch/issues/new)
[BENCHMARK.md]: (https://github.com/LdDl/ch/blob/master/BENCHMARK.md)
