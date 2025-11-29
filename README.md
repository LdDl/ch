[![GoDoc](https://godoc.org/github.com/LdDl/ch?status.svg)](https://godoc.org/github.com/LdDl/ch)
[![Build Status](https://travis-ci.com/LdDl/ch.svg?branch=master)](https://travis-ci.com/LdDl/ch)
[![Sourcegraph](https://sourcegraph.com/github.com/LdDl/ch/-/badge.svg)](https://sourcegraph.com/github.com/LdDl/ch?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/LdDl/ch)](https://goreportcard.com/report/github.com/LdDl/ch)
[![GitHub tag](https://img.shields.io/github/tag/LdDl/ch.svg)](https://github.com/LdDl/ch/releases)

# ch - Contraction Hierarchies
## Contraction Hierarchies - technique for speed up of computing shortest path in graph.

This library provides [Contraction Hierarchies](https://en.wikipedia.org/wiki/Contraction_hierarchies) preprocessing graph technique for [Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm). Classic implementation of Dijkstra's algorithm, maneuver restrictions extension and [isochrones](https://en.wikipedia.org/wiki/Isochrone_map) estimation are included also.

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
* Contraction hierarchies
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
go: finding github.com/LdDl/ch v1.8.0
go: downloading github.com/LdDl/ch v1.8.0
```
And then you are good to go 

## Usage

* Shortest path (single-threaded)

    Please see this [test file](bidirectional_ch_test.go#L17)

    I hope it's pretty clear, but here is little explanation:
    ```go
    g := Graph{} // Prepare variable for storing graph
    graphFromCSV(&g, "data/pgrouting_osm.csv") // Import CSV-file file into programm
    g.PrepareContractionHierarchies() // Compute contraction hierarchies
    u := 144031 // Define source vertex
    v := 452090 // Define target vertex
    ans, path := g.ShortestPath(u, v) // Get shortest path and it's cost between source and target vertex
    ```

* Shortest path (thread-safe for concurrent use)

    Please see this [test file](query_threadsafe_test.go#L11)

    If you need to execute shortest path queries from multiple goroutines concurrently, use the `QueryPool` API:
    ```go
    g := Graph{} // Prepare variable for storing graph
    graphFromCSV(&g, "data/pgrouting_osm.csv") // Import CSV-file file into programm
    g.PrepareContractionHierarchies() // Compute contraction hierarchies

    pool := g.NewQueryPool() // Create a query pool for concurrent access

    // Now you can safely call from multiple goroutines:
    ans, path := pool.ShortestPath(u, v)

    // One-to-many queries are also supported:
    costs, paths := pool.ShortestPathOneToMany(source, targets)
    ```

    **Important**: The default `Graph.ShortestPath()` method is NOT thread-safe. If you call it from multiple goroutines without synchronization, you may get incorrect results. Use `QueryPool` for concurrent scenarios.

* Isochrones

    Please see this [test file](isochrones_test.go#L7)
    ```go
    g := Graph{} // Prepare variable for storing graph
    // ...
    // Fill graph with data (vertices and edges)
    // ...
    isochrones, err := graph.Isochrones(sourceVertex, maxCost) // Evaluate isochrones via bread-first search
	if err != nil {
		t.Error(err)
		return
    }
    ```
    
### If you want to import OSM (Open Street Map) file then follow instructions for [osm2ch](https://github.com/LdDl/osm2ch#osm2ch)

## Benchmark

You can check benchmarks [here](https://github.com/LdDl/ch/blob/master/BENCHMARK.md)

## Support
If you have troubles or questions please [open an issue](https://github.com/LdDl/ch/issues/new).

## ToDo

Please see [ROADMAP.md](ROADMAP.md)

## Theory
[Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)

[Bidirectional search](https://en.wikipedia.org/wiki/Bidirectional_search)

[Bidirectional Dijkstra's algorithm's stop condition](http://www.cs.princeton.edu/courses/archive/spr06/cos423/Handouts/EPP%20shortest%20path%20algorithms.pdf)

[Contraction hierarchies](https://en.wikipedia.org/wiki/Contraction_hierarchies)

[Video Lectures](https://ad-wiki.informatik.uni-freiburg.de/teaching/EfficientRoutePlanningSS2012)


## Thanks
Thanks to [this visual explanation](https://jlazarsfeld.github.io/ch.150.project/contents/)
Thanks to [this](https://github.com/navjindervirdee/Advanced-Shortest-Paths-Algorithms) Java implementation of mentioned algorithms

## Dependencies
Thanks to [paulmach](https://github.com/paulmach) for his [OSM-parser](https://github.com/paulmach/osm) written in Go.

Paulmach's license is [here](https://github.com/paulmach/osm/blob/master/LICENSE.md) (it's MIT)

## License
You can check it [here](https://github.com/LdDl/ch/blob/master/LICENSE)

[osm2ch]: (https://github.com/LdDl/osm2ch#osm2ch)
[open an issue]: (https://github.com/LdDl/ch/issues/new)
[BENCHMARK.md]: (https://github.com/LdDl/ch/blob/master/BENCHMARK.md)
