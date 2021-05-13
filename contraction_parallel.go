package ch

import (
	"container/heap"
	"fmt"
	"runtime"
	"time"
)

// PreprocessParallel Same as Preprocess() but with parallelism
func (graph *Graph) PreprocessParallel() []int64 {
	numThreads := runtime.GOMAXPROCS(0)
	nodeOrdering := make([]int64, len(graph.Vertices))
	var extractNum int
	var iter int

	for i := range graph.Vertices {
		graph.Vertices[i].distance_v2 = make([]*Distance, numThreads)
		for threadID := 0; threadID < numThreads; threadID++ {
			graph.Vertices[i].distance_v2[threadID] = NewDistance()
		}
	}

	for graph.pqImportance.Len() != 0 {
		iter++
		vertex := heap.Pop(graph.pqImportance).(*Vertex)
		vertex.computeImportance()
		if graph.pqImportance.Len() != 0 && vertex.importance > graph.pqImportance.Peek().importance {
			graph.pqImportance.Push(vertex)
			continue
		}

		nodeOrdering[extractNum] = vertex.vertexNum
		vertex.orderPos = extractNum
		extractNum = extractNum + 1
		// graph.contractNodeParallel(vertex, int64(extractNum-1), distChan, numThreads)
		graph.contractNodeParallel_v2(vertex, int64(extractNum-1), numThreads)
		if DEBUG_PREPROCESSING {
			if iter > 0 && graph.pqImportance.Len()%1000 == 0 {
				fmt.Printf("Contraction Order: %d / %d, Remain vertices in heap: %d. Currect shortcuts num: %d Time: %v\n", extractNum, len(graph.Vertices), graph.pqImportance.Len(), graph.shortcutsNum(), time.Now())
			}
		}
	}
	return nodeOrdering
}

type chunk struct {
	fromIdx, toIdx int
}

// contractNodeParallel_v2 Same as contractNode() but with but with parallelism
func (graph *Graph) contractNodeParallel_v2(vertex *Vertex, contractID int64, threads int) {
	inEdges := vertex.inIncidentEdges
	outEdges := vertex.outIncidentEdges
	vertex.contracted = true
	inMax := 0.0
	outMax := 0.0
	graph.markNeighbors(inEdges, outEdges)
	for i := 0; i < len(inEdges); i++ {
		if graph.Vertices[inEdges[i].vertexID].contracted {
			continue
		}
		if inMax < inEdges[i].cost {
			inMax = inEdges[i].cost
		}
	}
	for i := 0; i < len(outEdges); i++ {
		if graph.Vertices[outEdges[i].vertexID].contracted {
			continue
		}
		if outMax < outEdges[i].cost {
			outMax = outEdges[i].cost
		}
	}
	max := inMax + outMax

	res := make(chan chunk, threads)
	limit := len(inEdges)
	lastIdx := 0

	chunkDone := chunk{}

	if limit < 2 {
		// When goroutines are not necessary;
		graph.workWithIncidentEdgesSingle(inEdges, outEdges, max, contractID, vertex.vertexNum)
	} else if limit >= 2 && threads < limit {
		// When number of goroutines is less then number of incoming edges
		// we should do batch processing
		graph.pqComparators = make([]*distanceHeap, threads)
		// fmt.Printf("Here threading for %d pathes\n", limit)
		for threadID := 0; threadID < threads; threadID++ {
			go func(thread int, r chan<- chunk) {
				start := (limit / threads) * thread
				end := start + (limit / threads)
				edgesSet := inEdges[start:end]
				lastIdx = end
				graph.callDijkstra(edgesSet, outEdges, max, contractID, vertex.vertexNum, thread)
				r <- chunk{start, end}
			}(threadID, res)
		}
		for threadID := 0; threadID < threads; threadID++ {
			chunkDone = <-res
			edgesSet := inEdges[chunkDone.fromIdx:chunkDone.toIdx]
			graph.callShortcuts(edgesSet, outEdges, contractID, vertex.vertexNum, threadID)
		}
		if lastIdx < len(inEdges) {
			graph.callDijkstra(inEdges[lastIdx:], outEdges, max, contractID, vertex.vertexNum, 0)
			graph.callShortcuts(inEdges[lastIdx:], outEdges, contractID, vertex.vertexNum, 0)
		}
	} else {
		// When number of goroutines is greater-or-equal to number of incoming edges
		// we should do batch processing with batch size = 1
		graph.pqComparators = make([]*distanceHeap, limit)
		for threadID := 0; threadID < limit; threadID++ {
			go func(thread int, r chan<- chunk) {
				start := thread
				end := start + 1
				edgesSet := inEdges[start:end]
				lastIdx = end
				graph.callDijkstra(edgesSet, outEdges, max, contractID, vertex.vertexNum, thread)
				r <- chunk{start, end}
			}(threadID, res)
		}
		for threadID := 0; threadID < limit; threadID++ {
			chunkDone = <-res
			edgesSet := inEdges[chunkDone.fromIdx:chunkDone.toIdx]
			graph.callShortcuts(edgesSet, outEdges, contractID, vertex.vertexNum, threadID)
		}
	}
}

func (graph *Graph) callDijkstra(inEdges []incidentEdge, outEdges []incidentEdge, max float64, contractID, vertexID int64, threadID int) {
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		graph.dijkstra_v2(inVertex, max, contractID, int64(i), threadID) // Finds the shortest distances from the inVertex to all outVertices.
	}
}

func (graph *Graph) callShortcuts(inEdges []incidentEdge, outEdges []incidentEdge, contractID, vertexID int64, threadID int) {
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inEdges[i].cost
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j].vertexID
			outcost := outEdges[j].cost
			if graph.Vertices[outVertex].contracted {
				continue
			}
			summaryCost := incost + outcost
			if graph.Vertices[outVertex].distance_v2[threadID].distance > summaryCost {
				graph.createOrUpdateShortcut(inVertex, outVertex, vertexID, summaryCost)
			}
		}
	}
}

func (graph *Graph) workWithIncidentEdgesSingle(inEdges []incidentEdge, outEdges []incidentEdge, max float64, contractID, vertexID int64) {
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inEdges[i].cost
		graph.dijkstra(inVertex, max, contractID, int64(i)) // Finds the shortest distances from the inVertex to all outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j].vertexID
			outcost := outEdges[j].cost
			if graph.Vertices[outVertex].contracted {
				continue
			}
			summaryCost := incost + outcost
			if graph.Vertices[outVertex].distance.contractID != contractID || graph.Vertices[outVertex].distance.sourceID != int64(i) || graph.Vertices[outVertex].distance.distance > summaryCost {
				graph.createOrUpdateShortcut(inVertex, outVertex, vertexID, summaryCost)
			}
		}
	}
}
