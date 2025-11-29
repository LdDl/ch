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
* Dynamic edge weight updates (lightweight recustomization)

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
go: finding github.com/LdDl/ch v1.10.0
go: downloading github.com/LdDl/ch v1.10.0
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

* Dynamic edge weight updates (Recustomization)

    Please see this [test file](recustomize_test.go)

    This feature allows you to update edge weights without rebuilding the entire contraction hierarchy. Useful for:
    - Traffic updates (congestion, accidents and other events)
    - Time-dependent routing

    ```go
    g := Graph{}
    graphFromCSV(&g, "data/pgrouting_osm.csv")
    g.PrepareContractionHierarchies()

    // Single update with immediate recustomization
    err := g.UpdateEdgeWeight(fromVertex, toVertex, newWeight, true)

    // Batch updates (more efficient for multiple changes)
    g.UpdateEdgeWeight(edge1From, edge1To, weight1, false)
    g.UpdateEdgeWeight(edge2From, edge2To, weight2, false)
    g.UpdateEdgeWeight(edge3From, edge3To, weight3, false)
    g.Recustomize() // Apply all changes at once
    ```

    **When to use single vs batch updates:**
    | Scenario | Method | Why |
    |----------|--------|-----|
    | One edge changed | `UpdateEdgeWeight(..., true)` | Simple, immediate |
    | Multiple edges changed | Batch + `Recustomize()` | Faster, single pass |
    | Real-time traffic feed | Batch + periodic `Recustomize()` | Amortize cost |

    **How it works:**

    ```mermaid
    flowchart TB
    subgraph Preprocessing["Preprocessing (one-time)"]
        P1[Build CH with importance ordering] --> P2[Store contractionOrder array]
        P2 --> P3[Index shortcuts by Via vertex<br/>shortcutsByVia map]
    end
    
    subgraph Update["UpdateEdgeWeight call"]
        U1[Convert user labels to internal IDs] --> U2[Update edge weight in<br/>outIncidentEdges & inIncidentEdges]
        U2 --> U3{needRecustom?}
        U3 -->|Yes| R1
        U3 -->|No| U4[Return - batch mode]
    end
    
    subgraph Recustomize["Recustomize call"]
        R1[For each vertex V in contractionOrder] --> R2[Get shortcuts via V<br/>from shortcutsByVia]
        R2 --> R3[For each shortcut A => C via V]
        R3 --> R4["newCost = cost(A => V) + cost(V => C)"]
        R4 --> R5[Update shortcut.Cost]
        R5 --> R6[Update incident edges]
        R6 --> R3
    end
    
    P3 --> U1
    R6 -.->|next vertex| R1
    ```

    Processing in contraction order ensures that when updating shortcut `A => C via V`, the edges `A => V` and `V => C` (which might themselves be shortcuts) have already been updated.

    **Note:** This is a lightweight recustomization inspired by [Customizable Contraction Hierarchies](https://arxiv.org/abs/1402.0402) (Dibbelt, Strasser, Wagner), but uses the existing importance-based ordering instead of nested dissection. It's simpler and requires no external dependencies, while still providing efficient metric updates.

    In future may be added full CCH support with nested dissection ordering (need to investigate METIS or similar libraries for graph partitioning).
    
### If you want to import OSM (Open Street Map) file then follow instructions for [osm2ch](https://github.com/LdDl/osm2ch#osm2ch)

### Custom import with pre-computed CH

If you have your own import logic (e.g., reading additional data like GeoJSON coordinates alongside the graph), you need to call `FinalizeImport()` after loading all vertices, edges, and shortcuts:

```go
graph := ch.NewGraph()

// Your custom import logic:
// - CreateVertex() for each vertex
// - AddEdge() for each edge
// - SetOrderPos() and SetImportance() for each vertex
// - AddShortcut() for each shortcut

graph.FinalizeImport() // Required for recustomization support

// Now graph is ready for queries and UpdateEdgeWeight/Recustomize
```

This is required because `FinalizeImport()` builds internal data structures (`contractionOrder`, `shortcutsByVia`) needed for [dynamic edge weight updates](#dynamic-edge-weight-updates-recustomization).

If you use the built-in `ImportFromFile()` function, this is called automatically.

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

[Customizable Contraction Hierarchies](https://arxiv.org/abs/1402.0402) - Dibbelt, Strasser, Wagner (2014). The recustomization feature in this library is inspired by CCH concepts.

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
