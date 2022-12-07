package ch

import (
	"container/heap"
	"container/list"
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

func (graph *Graph) initShortestPath() (queryDist [directionsCount][]float64, processed [directionsCount][]bool, queues [directionsCount]*vertexDistHeap) {
	for d := forward; d < directionsCount; d++ {
		queryDist[d] = make([]float64, len(graph.Vertices))
		for i := range queryDist[d] {
			queryDist[d][i] = Infinity
		}
		processed[d] = make([]bool, len(graph.Vertices))
		queues[d] = &vertexDistHeap{}
		heap.Init(queues[d])
	}
	return
}

func (graph *Graph) shortestPath(endpoints [directionsCount]int64) (float64, []int64) {
	queryDist, processed, queues := graph.initShortestPath()
	for d := forward; d < directionsCount; d++ {
		processed[d][endpoints[d]] = true
		queryDist[d][endpoints[d]] = 0
		heapEndpoint := &vertexDist{
			id:   endpoints[d],
			dist: 0,
		}
		heap.Push(queues[d], heapEndpoint)
	}
	return graph.shortestPathCore(queryDist, processed, queues)
}

func (graph *Graph) shortestPathCore(queryDist [directionsCount][]float64, processed [directionsCount][]bool, queues [directionsCount]*vertexDistHeap) (float64, []int64) {
	var prev [directionsCount]map[int64]int64
	for d := forward; d < directionsCount; d++ {
		prev[d] = make(map[int64]int64)
	}
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
			graph.directionalSearch(d, queues[d], processed[d], processed[reverseDirection], queryDist[d], queryDist[reverseDirection], prev[d], &estimate, &middleID)
		}
		if !queuesProcessed {
			break
		}
	}
	if estimate == Infinity {
		return -1.0, nil
	}
	return estimate, graph.ComputePath(middleID, prev[forward], prev[backward])
}

func (graph *Graph) directionalSearch(d direction, q *vertexDistHeap, localProcessed, reverseProcessed []bool, localQueryDist, reverseQueryDist []float64, prev map[int64]int64, estimate *float64, middleID *int64) {
	vertex := heap.Pop(q).(*vertexDist)
	if graph.Reporter != nil {
		graph.Reporter.VertexSettled(int(d), 0, vertex.id, q.Len())
	}
	if vertex.dist <= *estimate {
		localProcessed[vertex.id] = true
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
				alt := localQueryDist[vertex.id] + cost
				if localQueryDist[temp] > alt {
					if graph.Reporter != nil {
						if graph.Reporter != nil {
							graph.Reporter.EdgeRelaxed(int(d), 0, vertex.id, temp, true, q.Len())
						}
					}
					localQueryDist[temp] = alt
					prev[temp] = vertex.id
					node := &vertexDist{
						id:   temp,
						dist: alt,
					}
					heap.Push(q, node)
				}
			}
		}
	}
	if reverseProcessed[vertex.id] {
		if vertex.dist+reverseQueryDist[vertex.id] < *estimate {
			*middleID = vertex.id
			*estimate = vertex.dist + reverseQueryDist[vertex.id]
			if graph.Reporter != nil {
				graph.Reporter.FoundBetterPath(int(d), 0, 0, vertex.id, *estimate)
			}
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
	queryDist, processed, queues := graph.initShortestPath()
	for d := forward; d < directionsCount; d++ {
		for _, endpoint := range endpoints[d] {
			if endpoint.vertexNum == vertexNotFound {
				continue
			}
			processed[d][endpoint.vertexNum] = true
			queryDist[d][endpoint.vertexNum] = endpoint.additionalDistance
			heapEndpoint := &vertexDist{
				id:   endpoint.vertexNum,
				dist: endpoint.additionalDistance,
			}
			heap.Push(queues[d], heapEndpoint)
		}
	}
	return graph.shortestPathCore(queryDist, processed, queues)
}

// ComputePath Returns slice of IDs (user defined) of computed path
func (graph *Graph) ComputePath(middleID int64, forwardPrev, backwardPrev map[int64]int64) []int64 {
	l := list.New()
	l.PushBack(middleID)
	u := middleID
	var ok bool
	for {
		if u, ok = forwardPrev[u]; ok {
			l.PushFront(u)
		} else {
			break
		}
	}
	u = middleID
	for {
		if u, ok = backwardPrev[u]; ok {
			l.PushBack(u)
		} else {
			break
		}
	}
	ok = true
	for ok {
		ok = false
		for e := l.Front(); e.Next() != nil; e = e.Next() {
			if contractedNode, ok2 := graph.shortcuts[e.Value.(int64)][e.Next().Value.(int64)]; ok2 {
				ok = true
				l.InsertAfter(contractedNode.Via, e)
			}
		}
	}

	var path = make([]int64, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		path = append(path, graph.Vertices[e.Value.(int64)].Label)
	}

	return path
}
