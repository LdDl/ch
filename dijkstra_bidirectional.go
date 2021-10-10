package ch

import (
	"container/heap"
	"container/list"
	"math"
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

	forwardPrev := make(map[int64]int64)
	backwardPrev := make(map[int64]int64)

	queryDist := make([]float64, len(graph.Vertices))
	revQueryDist := make([]float64, len(graph.Vertices))

	forwProcessed := make([]bool, len(graph.Vertices))
	revProcessed := make([]bool, len(graph.Vertices))
	forwProcessed[source] = true
	revProcessed[target] = true

	for i := range queryDist {
		queryDist[i] = math.MaxFloat64
		revQueryDist[i] = math.MaxFloat64
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
		revQueryDistance: math.MaxFloat64,
	}
	heapTarget := &bidirectionalVertex{
		id:               target,
		queryDist:        math.MaxFloat64,
		revQueryDistance: 0,
	}

	heap.Push(forwQ, heapSource)
	heap.Push(backwQ, heapTarget)

	estimate := math.MaxFloat64

	var iter int

	var middleID int64

	for forwQ.Len() != 0 || backwQ.Len() != 0 {
		iter++
		if forwQ.Len() != 0 {
			vertex1 := heap.Pop(forwQ).(*bidirectionalVertex)
			if vertex1.queryDist <= estimate {
				forwProcessed[vertex1.id] = true
				graph.relaxEdgesBiForward(vertex1, forwQ, forwardPrev, queryDist)
			}
			if revProcessed[vertex1.id] {
				if vertex1.queryDist+revQueryDist[vertex1.id] < estimate {
					middleID = vertex1.id
					estimate = vertex1.queryDist + revQueryDist[vertex1.id]
				}
			}
		}

		if backwQ.Len() != 0 {
			vertex2 := heap.Pop(backwQ).(*bidirectionalVertex)
			if vertex2.revQueryDistance <= estimate {
				revProcessed[vertex2.id] = true
				graph.relaxEdgesBiBackward(vertex2, backwQ, backwardPrev, revQueryDist)
			}

			if forwProcessed[vertex2.id] {
				if vertex2.revQueryDistance+queryDist[vertex2.id] < estimate {
					middleID = vertex2.id
					estimate = vertex2.revQueryDistance + queryDist[vertex2.id]
				}
			}
		}

	}
	if estimate == math.MaxFloat64 {
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
