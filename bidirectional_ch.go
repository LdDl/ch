package ch

import (
	"container/heap"
	"container/list"
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
			if vertex1.queryDist <= estimate {
				forwProcessed[vertex1.id] = true
				graph.relaxEdgesBiForward(&vertex1, forwQ, prev, queryDist)
			}
			if revProcessed[vertex1.id] {
				if vertex1.queryDist+revDistance[vertex1.id] < estimate {
					middleID = vertex1.id
					estimate = vertex1.queryDist + revDistance[vertex1.id]
				}
			}
		}

		if backwQ.Len() != 0 {
			vertex2 := heap.Pop(backwQ).(simpleNode)
			if vertex2.revDistance <= estimate {
				revProcessed[vertex2.id] = true
				graph.relaxEdgesBiBackward(&vertex2, backwQ, prevReverse, revDistance)
			}

			if forwProcessed[vertex2.id] {
				if vertex2.revDistance+queryDist[vertex2.id] < estimate {
					middleID = vertex2.id
					estimate = vertex2.queryDist + queryDist[vertex2.id]
				}
			}
		}

	}
	// fmt.Println(middleID, estimate)
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

// ShortestPathProc -
func (graph *Graph) ShortestPathProc(source, target int) (float64, []int) {

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

	estimate1 := math.MaxFloat64
	var middleID1 int
	estimate2 := math.MaxFloat64
	var middleID2 int
	done := make(chan string)

	// forward
	go func(forwQ *forwardPropagationHeap, prev map[int]int, queryDist, revDistance []float64) {
		for forwQ.Len() != 0 {
			vertex1 := heap.Pop(forwQ).(simpleNode)
			if vertex1.queryDist <= estimate1 {
				forwProcessed[vertex1.id] = true
				graph.relaxEdgesBiForward(&vertex1, forwQ, prev, queryDist)
			}
			if revProcessed[vertex1.id] {
				if vertex1.queryDist+revDistance[vertex1.id] < estimate1 {
					middleID1 = vertex1.id
					estimate1 = vertex1.queryDist + revDistance[vertex1.id]
					// fmt.Println("1", middleID1, estimate1)
					done <- "forward"
					break
				}
			}
		}

	}(forwQ, prev, queryDist, revDistance)

	// backward
	go func(backwQ *backwardPropagationHeap, prevReverse map[int]int, queryDist, revDistance []float64) {
		for backwQ.Len() != 0 {
			vertex2 := heap.Pop(backwQ).(simpleNode)
			if vertex2.revDistance <= estimate2 {
				revProcessed[vertex2.id] = true
				graph.relaxEdgesBiBackward(&vertex2, backwQ, prevReverse, revDistance)
			}

			if forwProcessed[vertex2.id] {
				if vertex2.revDistance+queryDist[vertex2.id] < estimate2 {
					middleID2 = vertex2.id
					estimate2 = vertex2.revDistance + queryDist[vertex2.id]
					// fmt.Println("2", middleID2, estimate2)
					done <- "backward"
					break
				}
			}
		}
	}(backwQ, prevReverse, queryDist, revDistance)

	doneProcess := <-done
	// fmt.Println(doneProcess)

	// fmt.Println("end 1", middleID1, estimate1)
	// fmt.Println("end 2", middleID2, estimate2)

	if doneProcess == "forward" {
		if estimate1 == math.MaxFloat64 {
			return -1.0, nil
		}
		return estimate1, graph.ComputePath(middleID1, prev, prevReverse)
	} else if doneProcess == "backward" {
		if estimate2 == math.MaxFloat64 {
			return -1.0, nil
		}
		return estimate2, graph.ComputePath(middleID2, prev, prevReverse)
	}

	return -1.0, nil
	// 	// time.Sleep(5 * time.Second)

	// 	if estimate1 == math.MaxFloat64 {
	// 		return -1.0, nil
	// 	}
	// 	return estimate1, graph.ComputePath(middleID1, prev, prevReverse)
}
