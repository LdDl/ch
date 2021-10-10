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
//
func (graph *Graph) ShortestPath(source, target int64) (float64, []int64) {
	if source == target {
		return 0, []int64{source}
	}
	ok := true
	if source, ok = graph.mapping[source]; !ok {
		return -1.0, nil
	}
	if target, ok = graph.mapping[target]; !ok {
		return -1.0, nil
	}
	return graph.shortestPath(source, target)
}

func (graph *Graph) shortestPath(source, target int64) (float64, []int64) {
	forwardPrev := make(map[int64]int64)
	backwardPrev := make(map[int64]int64)

	queryDist := make([]float64, len(graph.Vertices))
	revQueryDist := make([]float64, len(graph.Vertices))

	forwProcessed := make([]bool, len(graph.Vertices))
	revProcessed := make([]bool, len(graph.Vertices))
	forwProcessed[source] = true
	revProcessed[target] = true

	for i := range queryDist {
		queryDist[i] = Infinity
		revQueryDist[i] = Infinity
	}
	queryDist[source] = 0
	revQueryDist[target] = 0

	forwQ := &forwardHeap{}
	backwQ := &backwardHeap{}

	heap.Init(forwQ)
	heap.Init(backwQ)

	heapSource := &bidirectionalVertex{
		id:               source,
		queryDist:        0,
		revQueryDistance: Infinity,
	}
	heapTarget := &bidirectionalVertex{
		id:               target,
		queryDist:        Infinity,
		revQueryDistance: 0,
	}

	heap.Push(forwQ, heapSource)
	heap.Push(backwQ, heapTarget)

	estimate := Infinity

	middleID := int64(-1)

	for forwQ.Len() != 0 || backwQ.Len() != 0 {
		// Upward search
		if forwQ.Len() != 0 {
			forwardVertex := heap.Pop(forwQ).(*bidirectionalVertex)
			if forwardVertex.queryDist <= estimate {
				forwProcessed[forwardVertex.id] = true
				// Edge relaxation in a forward propagation
				neighborsUpward := graph.Vertices[forwardVertex.id].outIncidentEdges
				for i := 0; i < len(neighborsUpward); i++ {
					temp := neighborsUpward[i].vertexID
					cost := neighborsUpward[i].weight
					if graph.Vertices[forwardVertex.id].orderPos < graph.Vertices[temp].orderPos {
						alt := queryDist[forwardVertex.id] + cost
						if queryDist[temp] > alt {
							queryDist[temp] = alt
							forwardPrev[temp] = forwardVertex.id
							node := &bidirectionalVertex{
								id:        temp,
								queryDist: alt,
							}
							heap.Push(forwQ, node)
						}
					}
				}
			}
			if revProcessed[forwardVertex.id] {
				if forwardVertex.queryDist+revQueryDist[forwardVertex.id] < estimate {
					middleID = forwardVertex.id
					estimate = forwardVertex.queryDist + revQueryDist[forwardVertex.id]
				}
			}
		}
		// Backward search
		if backwQ.Len() != 0 {
			backwardVertex := heap.Pop(backwQ).(*bidirectionalVertex)
			if backwardVertex.revQueryDistance <= estimate {
				revProcessed[backwardVertex.id] = true
				// Edge relaxation in a backward propagation
				vertexList := graph.Vertices[backwardVertex.id].inIncidentEdges
				for i := 0; i < len(vertexList); i++ {
					temp := vertexList[i].vertexID
					cost := vertexList[i].weight
					if graph.Vertices[backwardVertex.id].orderPos < graph.Vertices[temp].orderPos {
						alt := revQueryDist[backwardVertex.id] + cost
						if revQueryDist[temp] > alt {
							revQueryDist[temp] = alt
							backwardPrev[temp] = backwardVertex.id
							node := &bidirectionalVertex{
								id:               temp,
								revQueryDistance: alt,
							}
							heap.Push(backwQ, node)
						}
					}
				}

			}
			if forwProcessed[backwardVertex.id] {
				if backwardVertex.revQueryDistance+queryDist[backwardVertex.id] < estimate {
					middleID = backwardVertex.id
					estimate = backwardVertex.revQueryDistance + queryDist[backwardVertex.id]
				}
			}
		}

	}
	if estimate == Infinity {
		return -1.0, nil
	}
	return estimate, graph.ComputePath(middleID, forwardPrev, backwardPrev)
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
