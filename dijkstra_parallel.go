package ch

import (
	"container/heap"
	"sync"
)

type distanceHeapParallel struct {
	minheapSTD
}

type distanceParallel struct {
	contractID int
	sourceID   int64
	edgeID     int
	distance   []float64
	contract   []int
}

func createDistanceParallel(graph *Graph) distanceParallel {
	narr := distanceParallel{}
	narr.distance = make([]float64, len(graph.Vertices))
	narr.contract = make([]int, len(graph.Vertices))
	return narr
}

// relaxEdgesParallel Same as relaxEdges() but with parallelism
func (graph *Graph) relaxEdgesParallel(vertex int64, das *distanceParallel, pqComparator *distanceHeapParallel, iteration int) {
	vertexList := graph.Vertices[vertex].outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		if graph.Vertices[temp].contracted {
			continue
		}
		newPath := das.distance[vertex] + cost
		dist := das.distance[temp]
		if dist > newPath || das.contract[temp] != das.contractID {
			das.distance[temp] = newPath
			das.contract[temp] = das.contractID
			if iteration < 4 {
				heap.Push(pqComparator, minheapNode{temp, newPath, iteration + 1})
			}
		}

	}
}

// dijkstraParallel Same as dijkstra() but with parallelism
func (graph *Graph) dijkstraParallel(source int64, edgeID, contractID int, maxcost float64, distChan, outChan chan *distanceParallel, wg *sync.WaitGroup) {
	distance := <-distChan
	distance.contractID = contractID
	distance.edgeID = edgeID
	distance.sourceID = source
	pqComparator := &distanceHeapParallel{}
	heap.Init(pqComparator)
	sourceVertex := minheapNode{source, 0.0, 0}
	heap.Push(pqComparator, sourceVertex)
	distance.distance[source] = 0
	distance.contract[source] = contractID
	for pqComparator.Len() != 0 {
		vertex := heap.Pop(pqComparator).(minheapNode)
		if vertex.distance > maxcost {
			outChan <- distance
			wg.Done()
			return
		}
		graph.relaxEdgesParallel(vertex.id, distance, pqComparator, vertex.iteration)
	}
	outChan <- distance
	wg.Done()
}
