package ch

import (
	"container/heap"
	"container/list"
	"fmt"
	"log"
	"math"
)

// ShortestPath Computes and returns shortest path and it's cost (extended Dijkstra's algorithm)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source User's definied ID of source vertex
// target User's definied ID of target vertex
//
func (graph *Graph) ShortestPath(source, target int) (float64, []int) {

	if source == target {
		return 0, []int{source}
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

	prev := make(map[int]int)
	prevReverse := make(map[int]int)

	queryDist := make([]float64, len(graph.Vertices), len(graph.Vertices))
	revDistance := make([]float64, len(graph.Vertices), len(graph.Vertices))

	forwProcessed := make([]bool, len(graph.Vertices), len(graph.Vertices))
	revProcessed := make([]bool, len(graph.Vertices), len(graph.Vertices))
	forwProcessed[source] = true
	revProcessed[target] = true

	for i := range queryDist {
		queryDist[i] = math.MaxFloat64
		revDistance[i] = math.MaxFloat64
	}
	queryDist[source] = 0
	revDistance[target] = 0

	forwQ := &forwardPropagationHeap{}
	backwQ := &backwardPropagationHeap{}

	heap.Init(forwQ)
	heap.Init(backwQ)

	heapSource := simpleNode{
		id:          source,
		queryDist:   0,
		revDistance: math.MaxFloat64,
	}
	heapTarget := simpleNode{
		id:          target,
		queryDist:   math.MaxFloat64,
		revDistance: 0,
	}

	heap.Push(forwQ, heapSource)
	heap.Push(backwQ, heapTarget)

	estimate := math.MaxFloat64

	var iter int

	var middleID int

	for forwQ.Len() != 0 || backwQ.Len() != 0 {
		iter++
		if forwQ.Len() != 0 {
			vertex1 := heap.Pop(forwQ).(simpleNode)
			forwProcessed[vertex1.id] = true
			if vertex1.queryDist <= estimate {
				graph.relaxEdgesBiForward(&vertex1, forwQ, prev, queryDist, revDistance)
			}
			fmt.Println("Feed", graph.Vertices[vertex1.id].Label)
			if revProcessed[vertex1.id] {
				if vertex1.queryDist+revDistance[vertex1.id] < estimate {
					middleID = vertex1.id
					estimate = vertex1.queryDist + revDistance[vertex1.id]
				}
			}
		}

		if backwQ.Len() != 0 {
			vertex2 := heap.Pop(backwQ).(simpleNode)
			revProcessed[vertex2.id] = true

			if vertex2.revDistance <= estimate {
				graph.relaxEdgesBiBackward(&vertex2, backwQ, prevReverse, queryDist, revDistance)
			}
			fmt.Println("Back", graph.Vertices[vertex2.id].Label)

			if forwProcessed[vertex2.id] {
				if vertex2.revDistance+queryDist[vertex2.id] < estimate {
					middleID = vertex2.id
					estimate = vertex2.queryDist + queryDist[vertex2.id]
				}
			}
		}

	}
	if estimate == math.MaxFloat64 {
		return -1.0, nil
	}
	return estimate, graph.ComputePath(middleID, prev, prevReverse)
}

// ComputePath Returns slice of IDs (user defined) of computed path
func (graph *Graph) ComputePath(middleID int, prevF, prevR map[int]int) []int {
	l := list.New()
	l.PushBack(middleID)
	u := middleID
	var ok bool
	for {
		if u, ok = prevF[u]; ok {
			l.PushFront(u)
		} else {
			break
		}
	}
	u = middleID
	for {
		if u, ok = prevR[u]; ok {
			l.PushBack(u)
		} else {
			break
		}
	}
	ok = true
	for ok {
		ok = false
		for e := l.Front(); e.Next() != nil; e = e.Next() {
			if contractedNode, ok2 := graph.contracts[e.Value.(int)][e.Next().Value.(int)]; ok2 {
				ok = true
				l.InsertAfter(contractedNode, e)
			}
		}
	}

	var path []int
	for e := l.Front(); e != nil; e = e.Next() {
		path = append(path, graph.Vertices[e.Value.(int)].Label)
	}

	return path
}
