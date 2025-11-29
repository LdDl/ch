package ch

import (
	"container/heap"
)

type direction int

const (
	forward direction = iota
	backward
	directionsCount
)

// ShortestPath Computes and returns shortest path and it's cost (extended Dijkstra's algorithm)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source - user's definied ID of source vertex
// target - user's definied ID of target vertex
func (graph *Graph) ShortestPath(source, target int64) (float64, []int64) {
	if source == target {
		return 0, []int64{source}
	}
	endpoints := [directionsCount]int64{source, target}
	for d, endpoint := range endpoints {
		var ok bool
		if endpoints[d], ok = graph.mapping[endpoint]; !ok {
			return -1.0, nil
		}
	}
	return graph.shortestPath(endpoints)
}

func (graph *Graph) initShortestPath() (queues [directionsCount]*vertexDistHeap) {
	n := len(graph.Vertices)

	// Increment query epoch to clears all previous distances
	graph.queryEpoch++

	// Lazy initialization of query buffers (only on first query)
	if graph.queryDist[forward] == nil {
		for d := forward; d < directionsCount; d++ {
			graph.queryDist[d] = make([]float64, n)
			graph.queryEpochs[d] = make([]int64, n)
			graph.queryPrev[d] = make(map[int64]int64)
		}
	}

	// Clear prev maps by recreating them (reuses memory in Go's map implementation)
	for d := forward; d < directionsCount; d++ {
		graph.queryPrev[d] = make(map[int64]int64)
		queues[d] = &vertexDistHeap{}
		heap.Init(queues[d])
	}

	return
}

func (graph *Graph) shortestPath(endpoints [directionsCount]int64) (float64, []int64) {
	queues := graph.initShortestPath()
	for d := forward; d < directionsCount; d++ {
		graph.queryEpochs[d][endpoints[d]] = graph.queryEpoch
		graph.queryDist[d][endpoints[d]] = 0
		heapEndpoint := &vertexDist{
			id:   endpoints[d],
			dist: 0,
		}
		heap.Push(queues[d], heapEndpoint)
	}
	return graph.shortestPathCore(queues)
}

func (graph *Graph) shortestPathCore(queues [directionsCount]*vertexDistHeap) (float64, []int64) {
	estimate := Infinity
	middleID := int64(-1)
	for {
		queuesProcessed := false
		for d := forward; d < directionsCount; d++ {
			if queues[d].Len() == 0 {
				continue
			}
			queuesProcessed = true
			reverseDirection := (d + 1) % directionsCount
			graph.directionalSearch(d, queues[d], reverseDirection, &estimate, &middleID)
		}
		if !queuesProcessed {
			break
		}
	}
	if estimate == Infinity {
		return -1.0, nil
	}
	return estimate, graph.ComputePath(middleID, graph.queryPrev[forward], graph.queryPrev[backward])
}

func (graph *Graph) directionalSearch(d direction, q *vertexDistHeap, reverseDirection direction, estimate *float64, middleID *int64) {
	vertex := heap.Pop(q).(*vertexDist)
	if vertex.dist <= *estimate {
		graph.queryEpochs[d][vertex.id] = graph.queryEpoch
		// Edge relaxation in a forward propagation
		var vertexList []incidentEdge
		if d == forward {
			vertexList = graph.Vertices[vertex.id].outIncidentEdges
		} else {
			vertexList = graph.Vertices[vertex.id].inIncidentEdges
		}
		for i := range vertexList {
			temp := vertexList[i].vertexID
			cost := vertexList[i].weight
			if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
				alt := graph.queryDist[d][vertex.id] + cost
				// Check if temp was visited this epoch, if not treat as Infinity
				if graph.queryEpochs[d][temp] != graph.queryEpoch || graph.queryDist[d][temp] > alt {
					graph.queryDist[d][temp] = alt
					graph.queryEpochs[d][temp] = graph.queryEpoch
					graph.queryPrev[d][temp] = vertex.id
					node := &vertexDist{
						id:   temp,
						dist: alt,
					}
					heap.Push(q, node)
				}
			}
		}
	}
	// Check if reverse direction has processed this vertex
	if graph.queryEpochs[reverseDirection][vertex.id] == graph.queryEpoch {
		if vertex.dist+graph.queryDist[reverseDirection][vertex.id] < *estimate {
			*middleID = vertex.id
			*estimate = vertex.dist + graph.queryDist[reverseDirection][vertex.id]
		}
	}
}

// ShortestPathWithAlternatives Computes and returns shortest path and it's cost (extended Dijkstra's algorithm),
// with multiple alternatives for source and target vertices with additional distances to reach the vertices
// (useful if source and target are outside of the graph)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// sources - user's definied ID of source vertex with additional penalty
// targets - user's definied ID of target vertex with additional penalty
func (graph *Graph) ShortestPathWithAlternatives(sources, targets []VertexAlternative) (float64, []int64) {
	endpoints := [directionsCount][]VertexAlternative{sources, targets}
	var endpointsInternal [directionsCount][]vertexAlternativeInternal
	for d, alternatives := range endpoints {
		endpointsInternal[d] = graph.vertexAlternativesToInternal(alternatives)
	}
	return graph.shortestPathWithAlternatives(endpointsInternal)
}

func (graph *Graph) shortestPathWithAlternatives(endpoints [directionsCount][]vertexAlternativeInternal) (float64, []int64) {
	queues := graph.initShortestPath()
	for d := forward; d < directionsCount; d++ {
		for _, endpoint := range endpoints[d] {
			if endpoint.vertexNum == vertexNotFound {
				continue
			}
			graph.queryEpochs[d][endpoint.vertexNum] = graph.queryEpoch
			graph.queryDist[d][endpoint.vertexNum] = endpoint.additionalDistance
			heapEndpoint := &vertexDist{
				id:   endpoint.vertexNum,
				dist: endpoint.additionalDistance,
			}
			heap.Push(queues[d], heapEndpoint)
		}
	}
	return graph.shortestPathCore(queues)
}

// ComputePath Returns slice of IDs (user defined) of computed path
func (graph *Graph) ComputePath(middleID int64, forwardPrev, backwardPrev map[int64]int64) []int64 {
	// Build forward path (reversed, from middle to source)
	forwardPath := make([]int64, 0, 16)
	u := middleID
	for {
		prev, ok := forwardPrev[u]
		if !ok {
			break
		}
		forwardPath = append(forwardPath, prev)
		u = prev
	}

	// Build backward path (from middle to target)
	backwardPath := make([]int64, 0, 16)
	u = middleID
	for {
		next, ok := backwardPrev[u]
		if !ok {
			break
		}
		backwardPath = append(backwardPath, next)
		u = next
	}

	// Combine: reverse(forwardPath) + middle + backwardPath
	pathLen := len(forwardPath) + 1 + len(backwardPath)
	path := make([]int64, 0, pathLen)

	// Append reversed forward path
	for i := len(forwardPath) - 1; i >= 0; i-- {
		path = append(path, forwardPath[i])
	}
	path = append(path, middleID)
	path = append(path, backwardPath...)

	// Expand shortcuts iteratively
	for {
		expanded := false
		newPath := make([]int64, 0, len(path)*2)
		for i := 0; i < len(path); i++ {
			newPath = append(newPath, path[i])
			if i+1 < len(path) {
				if shortcut, ok := graph.shortcuts[path[i]][path[i+1]]; ok {
					newPath = append(newPath, shortcut.Via)
					expanded = true
				}
			}
		}
		path = newPath
		if !expanded {
			break
		}
	}

	// Convert internal IDs to user labels
	for i := range path {
		path[i] = graph.Vertices[path[i]].Label
	}

	return path
}
