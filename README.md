[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/LdDl/ch)
[![Build Status](https://travis-ci.com/LdDl/ch.svg?branch=master)](https://travis-ci.com/LdDl/ch)

# ch - Contraction Hierarchies
Contraction Hierarchies (with bidirectional version of Dijkstra's algorithm) for computing shortest path in graph.

This library provides classic implementation of Dijkstra's algorithm and turn restrictions extension also.

## Table of Contents

- [About](#about)
- [Installation](#installation)
    - [Old](#old-way)
    - [New](#new-way)
- [Usage](#usage)
- [Benchmark](#benchmark)
- [Support](#support)
- [ToDo](#todo)
- [Thanks](#thanks)
- [Theory](#theory)

## About
This package provides implemented next techniques and algorithms:
* Dijkstra's algorithm
* Contracntion hierarchies
* Bidirectional extension of Dijkstra's algorithm with contracted nodes

## Installation

### Old way
```go
go get github.com/LdDl/ch
```


### New way 
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
go: finding github.com/LdDl/ch v1.2.0
go: downloading github.com/LdDl/ch v1.2.0
```
And then you are good to go 

## Usage

Please see this [benchmark](bidirectional_ch_test.go#L59)

I hope it's pretty clear, but here is little explanation:
```go
    g := Graph{} // Prepare variable for storing graph
    graphFromCSV(&g, "data/pgrouting_osm.csv") // Import CSV-file file into programm
    g.PrepareContracts() // Compute contraction hierarchies
    u := 144031 // Define source vertex
    v := 452090 // Define target vertex
    ans, path := g.ShortestPath(u, v) // Get shortest path and it's cost between source and target vertex
```

## Benchmark

You can check benchmarks [here](BENCHMARK.md)

## Support

If you have troubles or questions please [open an issue](https://github.com/LdDl/ch/issues/new).

## ToDo

* ~~Import file of specific format~~ **Done as CSV**
* ~~Export file of specific format~~ **Done as CSV**
* Turn Restricted Shortest Path extension for CH-algorithm **Propably not modify algorithm, but graph**
* Thoughts and discussions about OSM graph and extensions **Need some ideas about parsing and preparing**
* Map matcher as another project **WIP (now it is in local git storage)**
* Bring interfaces{} **Thoughts**
* Bring OSM parser **WIP It exists, now need restrictions handle**
* Bring OSM restrictions **WIP PRs are welcome**
* ~~OneTwoMany function (contraction hierarchies)~~ **Done, ~~may be some bench comparisons~~**
* ManyToMany function (contraction hierarchies) **Thoughts**
* Replace int with int64 (OSM purposes) **Done**
* Separate benchmarks to BENCHMARK.md

## Theory
[Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)

[Bidirectional search](https://en.wikipedia.org/wiki/Bidirectional_search)

[Bidirectional Dijkstra's algorithm's stop condition](http://www.cs.princeton.edu/courses/archive/spr06/cos423/Handouts/EPP%20shortest%20path%20algorithms.pdf)

[Contraction hierarchies](https://en.wikipedia.org/wiki/Contraction_hierarchies)

[Video Lectures](https://ad-wiki.informatik.uni-freiburg.de/teaching/EfficientRoutePlanningSS2012)


## Thanks
Thanks to [this](https://github.com/navjindervirdee/Advanced-Shortest-Paths-Algorithms) Java implementation of mentioned algorithms
