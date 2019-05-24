[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/LdDl/ch)

# ch
Contraction Hierarchies (with bidirectional version of Dijkstra's algorithm) for computing shortest path in graph

## Table of Contents

- [About](#about)
- [Installation](#usage)
- [Usage](#usage)
- [Support](#support)
- [ToDo](#todo)
- [Thanks](#thanks)
- [Theory](#thanks)

## About
This package provides implemented next techniques and algorithms:
* Dijkstra's algorithm
* Contracntion hierarchies
* Bidirectional extension of Dijkstra's algorithm with contracted nodes

## Installation

Simple:
```go
go get github.com/LdDl/ch
```

## Usage

Please see this [benchmark](github.com/LdDl/ch/blob/master/bidirectional_ch_test.go#L13)

I hope it's pretty clear, but here is little explanation:
```go
    g := Graph{} // Prepare variable for storing graph
    graphFromCSV(&g, "benchmark_graph.csv") // Import CSV-file file into programm
    g.PrepareContracts() // Compute contraction hierarchies
    u := 144031 // Define source vertex
    v := 452090 // Define target vertex
    ans, path := g.ShortestPath(u, v) // Get shortest path and it's cost between source and target vertex
```

## Support

If you have troubles or questions please [open an issue](https://github.com/LdDl/ch/issues/new).

## ToDo

* Import file of specific format
* Export file of specific format
* Thoughts and discussions about OSM graph and extensions

## Theory
[Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)

[Bidirectional search](https://en.wikipedia.org/wiki/Bidirectional_search)

[Bidirectional Dijkstra's algorithm's stop condition](http://www.cs.princeton.edu/courses/archive/spr06/cos423/Handouts/EPP%20shortest%20path%20algorithms.pdf)

[Contraction hierarchies](https://en.wikipedia.org/wiki/Contraction_hierarchies)

[Video Lectures](https://ad-wiki.informatik.uni-freiburg.de/teaching/EfficientRoutePlanningSS2012)


## Thanks
Thanks to [this](https://github.com/navjindervirdee/Advanced-Shortest-Paths-Algorithms) Java implementation of mentioned algorithms