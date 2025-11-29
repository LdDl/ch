package ch

import (
	"container/heap"
)

func (graph *Graph) initShortestPathOneToMany() {
	n := len(graph.Vertices)

	// Lazy initialization of query buffers (only on first query)
	if graph.oneToManyDist[forward] == nil {
		for d := forward; d < directionsCount; d++ {
			graph.oneToManyDist[d] = make([]float64, n)
			graph.oneToManyEpochs[d] = make([]int64, n)
			graph.oneToManyPrev[d] = make(map[int64]int64)
		}
	}

	// Clear prev maps by recreating them
	for d := forward; d < directionsCount; d++ {
		graph.oneToManyPrev[d] = make(map[int64]int64)
	}
}

// ShortestPathOneToMany computes and returns shortest paths and theirs's costs (extended Dijkstra's algorithm) between single source and multiple targets
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source - user's definied ID of source vertex
// targets - set of user's definied IDs of target vertices
func (graph *Graph) ShortestPathOneToMany(source int64, targets []int64) ([]float64, [][]int64) {
	graph.initShortestPathOneToMany()

	estimateAll := make([]float64, 0, len(targets))
	pathAll := make([][]int64, 0, len(targets))

	var ok bool
	if source, ok = graph.mapping[source]; !ok {
		estimateAll = append(estimateAll, -1.0)
		pathAll = append(pathAll, nil)
		return estimateAll, pathAll
	}

	for _, target := range targets {
		// Increment epoch for each target query
		graph.oneToManyEpoch++
		epoch := graph.oneToManyEpoch

		// Clear prev maps for each target to avoid stale path reconstruction
		graph.oneToManyPrev[forward] = make(map[int64]int64)
		graph.oneToManyPrev[backward] = make(map[int64]int64)

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

		graph.oneToManyEpochs[forward][source] = epoch
		graph.oneToManyEpochs[backward][target] = epoch

		graph.oneToManyDist[forward][source] = 0
		graph.oneToManyDist[backward][target] = 0

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

		estimate, path := graph.shortestPathOneToManyCore(epoch, forwQ, backwQ)

		estimateAll = append(estimateAll, estimate)
		pathAll = append(pathAll, path)
	}

	return estimateAll, pathAll
}

func (graph *Graph) shortestPathOneToManyCore(epoch int64, forwQ *vertexDistHeap, backwQ *vertexDistHeap) (float64, []int64) {
	estimate := Infinity

	var middleID int64

	for forwQ.Len() != 0 || backwQ.Len() != 0 {
		if forwQ.Len() != 0 {
			vertex1 := heap.Pop(forwQ).(*vertexDist)
			if vertex1.dist <= estimate {
				graph.oneToManyEpochs[forward][vertex1.id] = epoch
				graph.relaxEdgesBiForwardOneToMany(vertex1, forwQ, epoch)
			}
			if graph.oneToManyEpochs[backward][vertex1.id] == epoch {
				if vertex1.dist+graph.oneToManyDist[backward][vertex1.id] < estimate {
					middleID = vertex1.id
					estimate = vertex1.dist + graph.oneToManyDist[backward][vertex1.id]
				}
			}
		}

		if backwQ.Len() != 0 {
			vertex2 := heap.Pop(backwQ).(*vertexDist)
			if vertex2.dist <= estimate {
				graph.oneToManyEpochs[backward][vertex2.id] = epoch
				graph.relaxEdgesBiBackwardOneToMany(vertex2, backwQ, epoch)
			}

			if graph.oneToManyEpochs[forward][vertex2.id] == epoch {
				if vertex2.dist+graph.oneToManyDist[forward][vertex2.id] < estimate {
					middleID = vertex2.id
					estimate = vertex2.dist + graph.oneToManyDist[forward][vertex2.id]
				}
			}
		}

	}
	if estimate == Infinity {
		return -1, nil
	}

	return estimate, graph.ComputePath(middleID, graph.oneToManyPrev[forward], graph.oneToManyPrev[backward])
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
	graph.initShortestPathOneToMany()

	estimateAll := make([]float64, 0, len(targetsAlternatives))
	pathAll := make([][]int64, 0, len(targetsAlternatives))

	sourceAlternativesInternal := graph.vertexAlternativesToInternal(sourceAlternatives)

	for _, targetAlternatives := range targetsAlternatives {
		// Increment epoch for each target query
		graph.oneToManyEpoch++
		epoch := graph.oneToManyEpoch

		// Clear prev maps for each target to avoid stale path reconstruction
		graph.oneToManyPrev[forward] = make(map[int64]int64)
		graph.oneToManyPrev[backward] = make(map[int64]int64)

		targetAlternativesInternal := graph.vertexAlternativesToInternal(targetAlternatives)

		forwQ := &vertexDistHeap{}
		backwQ := &vertexDistHeap{}

		heap.Init(forwQ)
		heap.Init(backwQ)

		for _, sourceAlternative := range sourceAlternativesInternal {
			if sourceAlternative.vertexNum == vertexNotFound {
				continue
			}
			graph.oneToManyEpochs[forward][sourceAlternative.vertexNum] = epoch
			graph.oneToManyDist[forward][sourceAlternative.vertexNum] = sourceAlternative.additionalDistance

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
			graph.oneToManyEpochs[backward][targetAlternative.vertexNum] = epoch
			graph.oneToManyDist[backward][targetAlternative.vertexNum] = targetAlternative.additionalDistance

			heapTarget := &vertexDist{
				id:   targetAlternative.vertexNum,
				dist: targetAlternative.additionalDistance,
			}
			heap.Push(backwQ, heapTarget)
		}

		estimate, path := graph.shortestPathOneToManyCore(epoch, forwQ, backwQ)

		estimateAll = append(estimateAll, estimate)
		pathAll = append(pathAll, path)
	}

	return estimateAll, pathAll
}

func (graph *Graph) relaxEdgesBiForwardOneToMany(vertex *vertexDist, forwQ *vertexDistHeap, epoch int64) {
	vertexList := graph.Vertices[vertex.id].outIncidentEdges
	for i := range vertexList {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := graph.oneToManyDist[forward][vertex.id] + cost
			if graph.oneToManyEpochs[forward][temp] != epoch || graph.oneToManyDist[forward][temp] > alt {
				graph.oneToManyDist[forward][temp] = alt
				graph.oneToManyPrev[forward][temp] = vertex.id
				graph.oneToManyEpochs[forward][temp] = epoch
				node := &vertexDist{
					id:   temp,
					dist: alt,
				}
				heap.Push(forwQ, node)
			}
		}
	}
}

func (graph *Graph) relaxEdgesBiBackwardOneToMany(vertex *vertexDist, backwQ *vertexDistHeap, epoch int64) {
	vertexList := graph.Vertices[vertex.id].inIncidentEdges
	for i := range vertexList {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := graph.oneToManyDist[backward][vertex.id] + cost
			if graph.oneToManyEpochs[backward][temp] != epoch || graph.oneToManyDist[backward][temp] > alt {
				graph.oneToManyDist[backward][temp] = alt
				graph.oneToManyPrev[backward][temp] = vertex.id
				graph.oneToManyEpochs[backward][temp] = epoch
				node := &vertexDist{
					id:   temp,
					dist: alt,
				}
				heap.Push(backwQ, node)
			}
		}
	}
}
