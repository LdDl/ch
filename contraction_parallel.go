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
		graph.contractNodeParallel_v2(vertex, numThreads)
		if DEBUG_PREPROCESSING {
			if iter > 0 && graph.pqImportance.Len()%1000 == 0 {
				fmt.Printf("Contraction Order: %d / %d, Remain vertices in heap: %d. Currect shortcuts num: %d Time: %v\n", extractNum, len(graph.Vertices), graph.pqImportance.Len(), graph.shortcutsNum(), time.Now())
			}
		}
	}
	return nodeOrdering
}

// contractNodeParallel_v2 Same as contractNode() but with but with parallelism
func (graph *Graph) contractNodeParallel_v2(vertexInfo *Vertex, threads int) {
	inEdges := vertexInfo.inIncidentEdges
	outEdges := vertexInfo.outIncidentEdges
	vertexInfo.contracted = true
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

	res := make(chan *ShortcutPathChannel, threads)

	inEdgesForProcess := []incidentEdge{}
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		inEdgesForProcess = append(inEdgesForProcess, inEdges[i])
	}

	limit := len(inEdgesForProcess)
	lastIdx := 0

	chunkDone := &ShortcutPathChannel{data: []*ShortcutPath{}, processed: 0}
	finalShortcuts := []*ShortcutPath{}

	contractionID := int64(vertexInfo.orderPos - 1)
	if limit < 2 {
		// When goroutines are not necessary;
		graph.workWithIncidentEdgesSingle(inEdgesForProcess, outEdges, max, contractionID, vertexInfo.vertexNum)
	} else if limit >= 2 && threads < limit {
		// When number of goroutines is less then number of incoming edges
		// we should do batch processing
		graph.pqComparators = make([]*distanceHeapParallel, threads)
		// fmt.Printf("Here threading for %d pathes\n", limit)
		for threadID := 0; threadID < threads; threadID++ {
			go func(thread int, r chan<- *ShortcutPathChannel) {
				start := (limit / threads) * thread
				end := start + (limit / threads)
				edgesSet := inEdgesForProcess[start:end]
				shortcuts := graph.callDijkstra(edgesSet, outEdges, max, contractionID, vertexInfo.vertexNum, thread)
				r <- &ShortcutPathChannel{shortcuts, len(edgesSet)}
			}(threadID, res)
		}
		for threadID := 0; threadID < threads; threadID++ {
			chunkDone = <-res
			lastIdx += chunkDone.processed
			finalShortcuts = append(finalShortcuts, chunkDone.data...)
		}
		if lastIdx < len(inEdgesForProcess) {
			// When there is last batch with size < number of threads
			// we should do batch processing with batch size = 1
			remainingInEdges := inEdgesForProcess[lastIdx:]
			limit = len(remainingInEdges)
			graph.pqComparators = make([]*distanceHeapParallel, limit)
			for threadID := 0; threadID < limit; threadID++ {
				go func(thread int, r chan<- *ShortcutPathChannel) {
					start := thread
					end := start + 1
					edgesSet := remainingInEdges[start:end]
					lastIdx = end
					shortcuts := graph.callDijkstra(edgesSet, outEdges, max, contractionID, vertexInfo.vertexNum, thread)
					r <- &ShortcutPathChannel{shortcuts, len(edgesSet)}
				}(threadID, res)
			}
			for threadID := 0; threadID < limit; threadID++ {
				chunkDone = <-res
				finalShortcuts = append(finalShortcuts, chunkDone.data...)
			}
		}
	} else {
		// When number of goroutines is greater-or-equal to number of incoming edges
		// we should do batch processing with batch size = 1
		graph.pqComparators = make([]*distanceHeapParallel, limit)
		for threadID := 0; threadID < limit; threadID++ {
			go func(thread int, r chan<- *ShortcutPathChannel) {
				start := thread
				end := start + 1
				edgesSet := inEdgesForProcess[start:end]
				lastIdx = end
				shortcuts := graph.callDijkstra(edgesSet, outEdges, max, contractionID, vertexInfo.vertexNum, thread)
				r <- &ShortcutPathChannel{shortcuts, len(edgesSet)}
			}(threadID, res)
		}
		for threadID := 0; threadID < limit; threadID++ {
			chunkDone = <-res
			finalShortcuts = append(finalShortcuts, chunkDone.data...)
		}
	}
	for i := range finalShortcuts {
		d := finalShortcuts[i]
		graph.createOrUpdateShortcut(d.From, d.To, d.Via, d.Cost)
	}
}

type ShortcutPathChannel struct {
	data      []*ShortcutPath
	processed int
}

func (graph *Graph) callDijkstra(inEdges []incidentEdge, outEdges []incidentEdge, max float64, contractionID, vertexID int64, threadID int) []*ShortcutPath {
	shortcuts := []*ShortcutPath{}
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inEdges[i].cost
		graph.dijkstra_v2(inVertex, max, contractionID, int64(i), threadID) // Finds the shortest distances from the inVertex to all outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j].vertexID
			outcost := outEdges[j].cost
			outVertexPtr := graph.Vertices[outVertex]
			if outVertexPtr.contracted {
				continue
			}
			summaryCost := incost + outcost
			if outVertexPtr.distance_v2[threadID].contractID != contractionID || outVertexPtr.distance_v2[threadID].sourceID != int64(i) || outVertexPtr.distance_v2[threadID].distance > summaryCost {
				shortcuts = append(shortcuts, &ShortcutPath{From: inVertex, To: outVertex, Via: vertexID, Cost: summaryCost})
			}
		}
	}
	return shortcuts
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
