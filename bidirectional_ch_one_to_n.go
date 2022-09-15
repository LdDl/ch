package ch

import (
	"container/heap"
)

// ShortestPathOneToMany Computes and returns shortest path and it's cost (extended Dijkstra's algorithm) for one-to-many relation
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source User's definied ID of source vertex
// targets User's definied IDs for target vertetices
func (graph *Graph) ShortestPathOneToMany(source int64, targets []int64) ([]float64, [][]int64) {
	estimateAll := []float64{}
	pathAll := [][]int64{}

	prev := make(map[int64]int64)
	prevReverse := make(map[int64]int64)

	queryDist := make([]float64, len(graph.Vertices))
	revQueryDist := make([]float64, len(graph.Vertices))

	forwProcessed := make([]int64, len(graph.Vertices))
	revProcessed := make([]int64, len(graph.Vertices))

	var ok bool
	if source, ok = graph.mapping[source]; !ok {
		estimateAll = append(estimateAll, -1.0)
		pathAll = append(pathAll, nil)
		return estimateAll, pathAll
	}

	for idx, target := range targets {
		nextQueue := int64(idx) + 1
		if source == target {
			estimateAll = append(estimateAll, 0)
			pathAll = append(pathAll, []int64{source})
			continue
		}
		var ok bool
		if target, ok = graph.mapping[target]; !ok {
			estimateAll = append(estimateAll, -1.0)
			pathAll = append(pathAll, nil)
			continue
		}

		forwProcessed[source] = nextQueue
		revProcessed[target] = nextQueue

		queryDist[source] = 0
		revQueryDist[target] = 0

		forwQ := &vertexDistHeap{}
		backwQ := &vertexDistHeap{}

		heap.Init(forwQ)
		heap.Init(backwQ)

		heapSource := &vertexDist{
			id:   source,
			dist: 0,
		}
		heapTarget := &vertexDist{
			id:   target,
			dist: 0,
		}

		heap.Push(forwQ, heapSource)
		heap.Push(backwQ, heapTarget)

		estimate := Infinity

		var middleID int64

		for forwQ.Len() != 0 || backwQ.Len() != 0 {
			if forwQ.Len() != 0 {
				vertex1 := heap.Pop(forwQ).(*vertexDist)
				if vertex1.dist <= estimate {
					forwProcessed[vertex1.id] = nextQueue
					graph.relaxEdgesBiForwardOneToMany(vertex1, forwQ, prev, queryDist, nextQueue, forwProcessed)
				}
				if revProcessed[vertex1.id] == nextQueue {
					if vertex1.dist+revQueryDist[vertex1.id] < estimate {
						middleID = vertex1.id
						estimate = vertex1.dist + revQueryDist[vertex1.id]
					}
				}
			}

			if backwQ.Len() != 0 {
				vertex2 := heap.Pop(backwQ).(*vertexDist)
				if vertex2.dist <= estimate {
					revProcessed[vertex2.id] = nextQueue
					graph.relaxEdgesBiBackwardOneToMany(vertex2, backwQ, prevReverse, revQueryDist, nextQueue, revProcessed)
				}

				if forwProcessed[vertex2.id] == nextQueue {
					if vertex2.dist+queryDist[vertex2.id] < estimate {
						middleID = vertex2.id
						estimate = vertex2.dist + queryDist[vertex2.id]
					}
				}
			}

		}
		if estimate == Infinity {
			estimateAll = append(estimateAll, -1)
			pathAll = append(pathAll, nil)
			continue
		}

		es, pp := estimate, graph.ComputePath(middleID, prev, prevReverse)

		estimateAll = append(estimateAll, es)
		pathAll = append(pathAll, pp)
	}

	return estimateAll, pathAll
}

func (graph *Graph) relaxEdgesBiForwardOneToMany(vertex *vertexDist, forwQ *vertexDistHeap, prev map[int64]int64, queryDist []float64, cid int64, forwProcessed []int64) {
	vertexList := graph.Vertices[vertex.id].outIncidentEdges
	for i := range vertexList {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := queryDist[vertex.id] + cost
			if forwProcessed[temp] != cid || queryDist[temp] > alt {
				queryDist[temp] = alt
				prev[temp] = vertex.id
				forwProcessed[temp] = cid
				node := &vertexDist{
					id:   temp,
					dist: alt,
				}
				heap.Push(forwQ, node)
			}
		}
	}
}

func (graph *Graph) relaxEdgesBiBackwardOneToMany(vertex *vertexDist, backwQ *vertexDistHeap, prev map[int64]int64, revQueryDist []float64, cid int64, revProcessed []int64) {
	vertexList := graph.Vertices[vertex.id].inIncidentEdges
	for i := range vertexList {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := revQueryDist[vertex.id] + cost
			if revProcessed[temp] != cid || revQueryDist[temp] > alt {
				revQueryDist[temp] = alt
				prev[temp] = vertex.id
				revProcessed[temp] = cid
				node := &vertexDist{
					id:   temp,
					dist: alt,
				}
				heap.Push(backwQ, node)
			}
		}
	}
}
