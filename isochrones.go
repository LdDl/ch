package ch

import (
	"container/heap"
	"fmt"
)

// Isochrones Returns set of vertices and corresponding distances restricted by maximum travel cost for source vertex
// source - source vertex (user defined label)
// maxCost - restriction on travel cost for breadth search
// See ref. https://wiki.openstreetmap.org/wiki/Isochrone and https://en.wikipedia.org/wiki/Isochrone_map
// Note: implemented breadth-first searching path algorithm does not guarantee shortest pathes to reachable vertices (until all edges have cost 1.0). See ref: https://en.wikipedia.org/wiki/Breadth-first_search
// Note: result for estimated costs could be also inconsistent due nature of data structure
func (graph *Graph) Isochrones(source int64, maxCost float64) (map[int64]float64, error) {
	ok := true
	if source, ok = graph.mapping[source]; !ok {
		return nil, fmt.Errorf("No such source")
	}
	Q := &minheapSTD{}
	heap.Init(Q)
	distance := make(map[int64]float64, len(graph.Vertices))
	Q.Push(minheapNode{id: source, distance: 0})
	visit := make(map[int64]bool)
	for Q.Len() != 0 {
		next := heap.Pop(Q).(minheapNode)
		visit[next.id] = true
		if next.distance <= maxCost {
			distance[graph.Vertices[next.id].Label] = next.distance
			vertexList := graph.Vertices[next.id].outEdges
			costList := graph.Vertices[next.id].outECost
			for i := range vertexList {
				neighbor := vertexList[i]
				if v1, ok1 := graph.contracts[next.id]; ok1 {
					if _, ok2 := v1[neighbor]; ok2 {
						// Ignore contract
						continue
					}
				}
				target := vertexList[i]
				cost := costList[i]
				alt := distance[graph.Vertices[next.id].Label] + cost
				if visit[target] {
					continue
				}
				Q.Push(minheapNode{id: target, distance: alt})
			}
		}
	}
	return distance, nil
}
