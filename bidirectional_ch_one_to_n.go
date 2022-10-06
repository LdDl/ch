package ch

import (
	"container/heap"
)

func (graph *Graph) initShortestPathOneToMany() (estimateAll []float64, pathAll [][]int64, prev map[int64]int64, prevReverse map[int64]int64, queryDist, revQueryDist []float64, forwProcessed, revProcessed []int64) {
	estimateAll = []float64{}
	pathAll = [][]int64{}

	prev = make(map[int64]int64)
	prevReverse = make(map[int64]int64)

	queryDist = make([]float64, len(graph.Vertices))
	revQueryDist = make([]float64, len(graph.Vertices))

	forwProcessed = make([]int64, len(graph.Vertices))
	revProcessed = make([]int64, len(graph.Vertices))

	return
}

// ShortestPathOneToMany computes and returns shortest paths and theirs's costs (extended Dijkstra's algorithm) between single source and multiple targets
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source - user's definied ID of source vertex
// targets - set of user's definied IDs of target vertices
func (graph *Graph) ShortestPathOneToMany(source int64, targets []int64) ([]float64, [][]int64) {
	estimateAll, pathAll, prev, prevReverse, queryDist, revQueryDist, forwProcessed, revProcessed := graph.initShortestPathOneToMany()

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

		estimate, path := graph.shortestPathOneToManyCore(nextQueue, prev, prevReverse, queryDist, revQueryDist, forwProcessed, revProcessed, forwQ, backwQ)

		estimateAll = append(estimateAll, estimate)
		pathAll = append(pathAll, path)
	}

	return estimateAll, pathAll
}

func (graph *Graph) shortestPathOneToManyCore(nextQueue int64, prev map[int64]int64, prevReverse map[int64]int64, queryDist, revQueryDist []float64, forwProcessed, revProcessed []int64, forwQ *vertexDistHeap, backwQ *vertexDistHeap) (float64, []int64) {
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
		return -1, nil
	}

	return estimate, graph.ComputePath(middleID, prev, prevReverse)
}

// ShortestPathOneToManyWithAlternatives Computes and returns shortest path and it's cost (extended Dijkstra's algorithm) between single source and multiple targets
// with multiple alternatives for source and target vertices with additional distances to reach the vertices
// (useful if source and target are outside of the graph)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// sourceAlternatives - user's definied ID of source vertex with additional penalty
// targetsAlternatives - set of user's definied IDs of target vertices  with additional penalty
func (graph *Graph) ShortestPathOneToManyWithAlternatives(sourceAlternatives []VertexAlternative, targetsAlternatives [][]VertexAlternative) ([]float64, [][]int64) {
	estimateAll, pathAll, prev, prevReverse, queryDist, revQueryDist, forwProcessed, revProcessed := graph.initShortestPathOneToMany()

	sourceAlternativesInternal := graph.vertexAlternativesToInternal(sourceAlternatives)

	for idx, targetAlternatives := range targetsAlternatives {
		nextQueue := int64(idx) + 1

		targetAlternativesInternal := graph.vertexAlternativesToInternal(targetAlternatives)

		forwQ := &vertexDistHeap{}
		backwQ := &vertexDistHeap{}

		heap.Init(forwQ)
		heap.Init(backwQ)

		for _, sourceAlternative := range sourceAlternativesInternal {
			if sourceAlternative.vertexNum == vertexNotFound {
				continue
			}
			forwProcessed[sourceAlternative.vertexNum] = nextQueue
			queryDist[sourceAlternative.vertexNum] = sourceAlternative.additionalDistance

			heapSource := &vertexDist{
				id:   sourceAlternative.vertexNum,
				dist: sourceAlternative.additionalDistance,
			}
			heap.Push(forwQ, heapSource)
		}
		for _, targetAlternative := range targetAlternativesInternal {
			if targetAlternative.vertexNum == vertexNotFound {
				continue
			}
			revProcessed[targetAlternative.vertexNum] = nextQueue
			revQueryDist[targetAlternative.vertexNum] = targetAlternative.additionalDistance

			heapTarget := &vertexDist{
				id:   targetAlternative.vertexNum,
				dist: targetAlternative.additionalDistance,
			}
			heap.Push(backwQ, heapTarget)
		}

		estimate, path := graph.shortestPathOneToManyCore(nextQueue, prev, prevReverse, queryDist, revQueryDist, forwProcessed, revProcessed, forwQ, backwQ)

		estimateAll = append(estimateAll, estimate)
		pathAll = append(pathAll, path)
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
