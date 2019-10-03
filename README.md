[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/LdDl/ch)
[![Build Status](https://travis-ci.com/LdDl/ch.svg?branch=master)](https://travis-ci.com/LdDl/ch)

# ch - Contraction Hierarchies
Contraction Hierarchies (with bidirectional version of Dijkstra's algorithm) for computing shortest path in graph.

This library provides classic implementation of Dijkstra's algorithm and turn restrictions extension also.

## Table of Contents

- [About](#about)
- [Installation](#usage)
    - [Old](#old-way)
    - [New](#new-way)
- [Usage](#usage)
- [Benchmark] (#benchmark)
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

My PC is:

    Processor: Intel(R) Core(TM) i9-7900X CPU @ 3.30GHz x 10
    Memory: 46.8GiB
    Linux Kernel: 4.15.0-20-generic
    OS: Linux Mint 19.1 Cinnamon

I have used graph with ~187k vertices for benchmark.

```bash
goos: linux
goarch: amd64
pkg: github.com/LdDl/ch
BenchmarkShortestPath/CH_shortest_path/1/vertices-187853-20         	     500	   2801587 ns/op	 3532241 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/2/vertices-187853-20         	    1000	   2639499 ns/op	 3532225 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/4/vertices-187853-20         	    1000	   2730468 ns/op	 3532239 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/8/vertices-187853-20         	     500	   2887250 ns/op	 3532254 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/16/vertices-187853-20        	     500	   2292956 ns/op	 3532251 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/32/vertices-187853-20        	     500	   2837590 ns/op	 3532247 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/64/vertices-187853-20        	     500	   2649959 ns/op	 3532233 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/128/vertices-187853-20       	     500	   2790797 ns/op	 3532215 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/256/vertices-187853-20       	     500	   2640733 ns/op	 3532231 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/512/vertices-187853-20       	     500	   2381726 ns/op	 3532224 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/1024/vertices-187853-20      	     500	   2810581 ns/op	 3532223 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/2048/vertices-187853-20      	     500	   2770308 ns/op	 3532203 B/op	    2225 allocs/op
BenchmarkShortestPath/CH_shortest_path/4096/vertices-187853-20      	     500	   2592263 ns/op	 3532234 B/op	    2225 allocs/op
PASS
ok  	github.com/LdDl/ch	34.153s
```
## Support

If you have troubles or questions please [open an issue](https://github.com/LdDl/ch/issues/new).

## ToDo

* ~~Import file of specific format~~ **Done as CSV**
* ~~Export file of specific format~~ **Done as CSV**
* Turn Restricted Shortest Path extension for CH-algorithm
* Turn restrictions import (probably as CSV again)
* Thoughts and discussions about OSM graph and extensions **Need some ideas about parsing and preparing**
* Map matcher as another project **WIP**

## Theory
[Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)

[Bidirectional search](https://en.wikipedia.org/wiki/Bidirectional_search)

[Bidirectional Dijkstra's algorithm's stop condition](http://www.cs.princeton.edu/courses/archive/spr06/cos423/Handouts/EPP%20shortest%20path%20algorithms.pdf)

[Contraction hierarchies](https://en.wikipedia.org/wiki/Contraction_hierarchies)

[Video Lectures](https://ad-wiki.informatik.uni-freiburg.de/teaching/EfficientRoutePlanningSS2012)


## Thanks
Thanks to [this](https://github.com/navjindervirdee/Advanced-Shortest-Paths-Algorithms) Java implementation of mentioned algorithms
