package ch

import (
	"container/heap"
	"container/list"
)

// ShortestPath Computes and returns shortest path and it's cost (extended Dijkstra's algorithm)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source User's definied ID of source vertex
// target User's definied ID of target vertex
func (graph *Graph) ShortestPath(source, target int64) (float64, []int64) {
	if source == target {
		return 0, []int64{source}
	}
	var ok bool
	if source, ok = graph.mapping[source]; !ok {
		return -1.0, nil
	}
	if target, ok = graph.mapping[target]; !ok {
		return -1.0, nil
	}
	return graph.shortestPath(source, target)
}

func (graph *Graph) initShortestPath() (
	queryDist, revQueryDist []float64,
	forwProcessed, revProcessed []bool,
	forwQ *vertexDistHeap,
	backwQ *vertexDistHeap,
) {
	queryDist = make([]float64, len(graph.Vertices))
	revQueryDist = make([]float64, len(graph.Vertices))

	for i := range queryDist {
		queryDist[i] = Infinity
		revQueryDist[i] = Infinity
	}

	forwProcessed = make([]bool, len(graph.Vertices))
	revProcessed = make([]bool, len(graph.Vertices))

	forwQ = &vertexDistHeap{}
	backwQ = &vertexDistHeap{}

	heap.Init(forwQ)
	heap.Init(backwQ)

	return
}

func (graph *Graph) shortestPath(source, target int64) (float64, []int64) {
	queryDist, revQueryDist, forwProcessed, revProcessed, forwQ, backwQ := graph.initShortestPath()

	forwProcessed[source] = true
	revProcessed[target] = true

	queryDist[source] = 0
	revQueryDist[target] = 0

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

	return graph.shortestPathCore(queryDist, revQueryDist, forwProcessed, revProcessed, forwQ, backwQ)
}

func (graph *Graph) shortestPathCore(
	queryDist, revQueryDist []float64,
	forwProcessed, revProcessed []bool,
	forwQ *vertexDistHeap,
	backwQ *vertexDistHeap,
) (float64, []int64) {
	forwardPrev := make(map[int64]int64)
	backwardPrev := make(map[int64]int64)

	estimate := Infinity

	middleID := int64(-1)

	for forwQ.Len() != 0 || backwQ.Len() != 0 {
		// Upward search
		if forwQ.Len() != 0 {
			forwardVertex := heap.Pop(forwQ).(*vertexDist)
			if forwardVertex.dist <= estimate {
				forwProcessed[forwardVertex.id] = true
				// Edge relaxation in a forward propagation
				neighborsUpward := graph.Vertices[forwardVertex.id].outIncidentEdges
				for i := range neighborsUpward {
					temp := neighborsUpward[i].vertexID
					cost := neighborsUpward[i].weight
					if graph.Vertices[forwardVertex.id].orderPos < graph.Vertices[temp].orderPos {
						alt := queryDist[forwardVertex.id] + cost
						if queryDist[temp] > alt {
							queryDist[temp] = alt
							forwardPrev[temp] = forwardVertex.id
							node := &vertexDist{
								id:   temp,
								dist: alt,
							}
							heap.Push(forwQ, node)
						}
					}
				}
			}
			if revProcessed[forwardVertex.id] {
				if forwardVertex.dist+revQueryDist[forwardVertex.id] < estimate {
					middleID = forwardVertex.id
					estimate = forwardVertex.dist + revQueryDist[forwardVertex.id]
				}
			}
		}
		// Backward search
		if backwQ.Len() != 0 {
			backwardVertex := heap.Pop(backwQ).(*vertexDist)
			if backwardVertex.dist <= estimate {
				revProcessed[backwardVertex.id] = true
				// Edge relaxation in a backward propagation
				vertexList := graph.Vertices[backwardVertex.id].inIncidentEdges
				for i := range vertexList {
					temp := vertexList[i].vertexID
					cost := vertexList[i].weight
					if graph.Vertices[backwardVertex.id].orderPos < graph.Vertices[temp].orderPos {
						alt := revQueryDist[backwardVertex.id] + cost
						if revQueryDist[temp] > alt {
							revQueryDist[temp] = alt
							backwardPrev[temp] = backwardVertex.id
							node := &vertexDist{
								id:   temp,
								dist: alt,
							}
							heap.Push(backwQ, node)
						}
					}
				}

			}
			if forwProcessed[backwardVertex.id] {
				if backwardVertex.dist+queryDist[backwardVertex.id] < estimate {
					middleID = backwardVertex.id
					estimate = backwardVertex.dist + queryDist[backwardVertex.id]
				}
			}
		}

	}
	if estimate == Infinity {
		return -1.0, nil
	}
	return estimate, graph.ComputePath(middleID, forwardPrev, backwardPrev)
}

type VertexAlternative struct {
	Label              int64
	AdditionalDistance float64
}

// ShortestPathWithAlternatives Computes and returns shortest path and it's cost (extended Dijkstra's algorithm),
// with multiple alternatives for source and target vertices with additional distances to reach the vertices
// (useful if source and target are outside of the graph)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// sources Source vertex alternatives
// targets Target vertex alternatives
func (graph *Graph) ShortestPathWithAlternatives(sources, targets []VertexAlternative) (float64, []int64) {
	sourcesInternal := make([]vertexAlternativeInternal, 0, len(sources))
	targetsInternal := make([]vertexAlternativeInternal, 0, len(targets))
	for _, source := range sources {
		sourceInternal := vertexAlternativeInternal{additionalDistance: source.AdditionalDistance}
		var ok bool
		if sourceInternal.vertexNum, ok = graph.mapping[source.Label]; !ok {
			return -1.0, nil
		}
		sourcesInternal = append(sourcesInternal, sourceInternal)
	}
	for _, target := range targets {
		targetInternal := vertexAlternativeInternal{additionalDistance: target.AdditionalDistance}
		var ok bool
		if targetInternal.vertexNum, ok = graph.mapping[target.Label]; !ok {
			return -1.0, nil
		}
		targetsInternal = append(targetsInternal, targetInternal)
	}
	return graph.shortestPathWithAlternatives(sourcesInternal, targetsInternal)
}

type vertexAlternativeInternal struct {
	vertexNum          int64
	additionalDistance float64
}

func (graph *Graph) shortestPathWithAlternatives(sources, targets []vertexAlternativeInternal) (float64, []int64) {
	queryDist, revQueryDist, forwProcessed, revProcessed, forwQ, backwQ := graph.initShortestPath()

	for _, source := range sources {
		forwProcessed[source.vertexNum] = true
		queryDist[source.vertexNum] = source.additionalDistance
		heapSource := &vertexDist{
			id:   source.vertexNum,
			dist: source.additionalDistance,
		}
		heap.Push(forwQ, heapSource)
	}
	for _, target := range targets {
		revProcessed[target.vertexNum] = true
		revQueryDist[target.vertexNum] = target.additionalDistance
		heapTarget := &vertexDist{
			id:   target.vertexNum,
			dist: target.additionalDistance,
		}
		heap.Push(backwQ, heapTarget)
	}

	return graph.shortestPathCore(queryDist, revQueryDist, forwProcessed, revProcessed, forwQ, backwQ)
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

	var path []int64
	for e := l.Front(); e != nil; e = e.Next() {
		path = append(path, graph.Vertices[e.Value.(int64)].Label)
	}

	return path
}
