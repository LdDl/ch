package ch

import (
	"container/heap"
	"log"
)

// VanillaShortestPath Computes and returns shortest path and it's cost (vanilla Dijkstra's algorithm)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source User's definied ID of source vertex
// target User's definied ID of target vertex
//
// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
func (graph *Graph) VanillaShortestPath(source, target int64) (float64, []int64) {

	if source == target {
		return 0, []int64{source}
	}
	ok := true

	if source, ok = graph.mapping[source]; !ok {
		log.Println("No such source")
		return -1.0, nil
	}
	if target, ok = graph.mapping[target]; !ok {
		log.Println("No such target")
		return -1.0, nil
	}

	// create vertex set Q
	Q := &minHeap{}

	// dist[source] ← 0
	distance := make(map[int64]float64, len(graph.Vertices))
	distance[source] = 0

	// prev[v] ← UNDEFINED
	prev := make(map[int64]int64, len(graph.Vertices))
	// for each vertex v in Graph:
	for i := range graph.Vertices {
		// if v ≠ source:
		if graph.Vertices[i].vertexNum != source {
			// dist[v] = INFINITY
			distance[graph.Vertices[i].vertexNum] = Infinity
		}
		// prev[v] ← UNDEFINED
		// nothing here
	}
	Q.add_with_priority(graph.Vertices[source].vertexNum, distance[graph.Vertices[source].vertexNum])
	heap.Init(Q)
	// while Q is not empty:
	for Q.Len() != 0 {
		// u ← Q.extract_min()
		u := heap.Pop(Q).(*minHeapVertex)

		// if u == target:
		if u.id == target {
			// break
			break
		}

		vertexList := graph.Vertices[u.id].outIncidentEdges

		// for each neighbor v of u:
		for v := range vertexList {
			neighbor := vertexList[v].vertexID
			if v1, ok1 := graph.shortcuts[u.id]; ok1 {
				if _, ok2 := v1[neighbor]; ok2 {
					// Ignore shortcut
					continue
				}
			}
			cost := vertexList[v].weight
			// alt ← dist[u] + length(u, v)
			alt := distance[u.id] + cost
			// if alt < dist[v]
			if distance[neighbor] > alt {
				// dist[v] ← alt
				distance[neighbor] = alt
				// prev[v] ← u
				prev[neighbor] = u.id
				// Q.decrease_priority(v, alt)
				// Q.decrease_priority(v, alt)
				Q.add_with_priority(neighbor, alt)
			}
		}
		// heap.Init(Q)
	}

	if distance[target] == Infinity {
		return -1, []int64{}
	}
	// path = []
	var path []int64
	// u = target
	u := target

	// while prev[u] is defined:
	for {
		if _, ok := prev[u]; !ok {
			break
		}
		// path.push_front(u)
		temp := make([]int64, len(path)+1)
		temp[0] = u
		copy(temp[1:], path)
		path = temp

		// u = prev[u]
		u = prev[u]
	}

	temp := make([]int64, len(path)+1)
	temp[0] = source
	copy(temp[1:], path)
	path = temp

	usersLabelsPath := make([]int64, len(path))
	for e := 0; e < len(usersLabelsPath); e++ {
		usersLabelsPath[e] = graph.Vertices[path[e]].Label
	}

	// return path, prev
	return distance[target], usersLabelsPath
}
