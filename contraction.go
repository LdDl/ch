package ch

import (
	"container/heap"
)

// Preprocess Computes contraction hierarchies and returns node ordering
func (graph *Graph) Preprocess() []int {
	nodeOrdering := make([]int, len(graph.Vertices))
	var extractNum int
	var iter int
	for graph.pqImportance.Len() != 0 {
		iter++
		// Lazy update heuristic:
		// update Importance of vertex "on demand" as follows:
		// Before contracting vertex with currently smallest Importance, recompute its Importance and see if it is still the smallest
		// If not pick next smallest one, recompute its Importance and see if that is the smallest now; If not, continue in same way ...
		vertex := heap.Pop(graph.pqImportance).(*Vertex)
		vertex.computeImportance()
		if graph.pqImportance.Len() != 0 && vertex.importance > graph.pqImportance.Peek().(*Vertex).importance {
			graph.pqImportance.Push(vertex)
			continue
		}

		nodeOrdering[extractNum] = vertex.vertexNum
		vertex.orderPos = extractNum
		extractNum = extractNum + 1
		graph.contractNode(vertex, int(extractNum-1))

		// Neighbours only heuristic:
		// After each contraction, recompute Importance, but only for the neighbours of the contracted node
		// for i := 0; i < len(vertex.inEdges); i++ {
		// 	inVertex := vertex.inEdges[i]
		// 	if !graph.Vertices[inVertex].contracted {
		// 		graph.Vertices[inVertex].computeImportance(graph)
		// 	}
		// }
		// for i := 0; i < len(vertex.outEdges); i++ {
		// 	outVertex := vertex.outEdges[i]
		// 	if !graph.Vertices[outVertex].contracted {
		// 		graph.Vertices[outVertex].computeImportance(graph)
		// 	}
		// }

		// Periodic update heuristic: Full recomputation every x rounds
		// if iter > 0 && graph.pqImportance.Len()%10000 == 0 {
		// 	for i := 0; i < len(graph.Vertices); i++ {
		// 		if !graph.Vertices[i].contracted {
		// 			graph.Vertices[i].computeImportance(graph)
		// 		}
		// 	}
		// }

		// if iter > 0 && graph.pqImportance.Len()%1000 == 0 {
		// 	fmt.Printf("Contraction Order: %d / %d, Remain vertices in heap: %d. Currect contractions: %d Time: %v\n", extractNum, len(graph.Vertices), graph.pqImportance.Len(), len(graph.shortcuts), time.Now())
		// }
	}
	return nodeOrdering
}

// markNeighbors Saves uniformity (affects importance of vertex)
//
// inEdges Incoming edges from vertex
// outEdges Outcoming edges from vertex
//
func (graph *Graph) markNeighbors(inEdges, outEdges []int) {
	for i := 0; i < len(inEdges); i++ {
		temp := inEdges[i]
		graph.Vertices[temp].deletedNeighbors++
	}
	for i := 0; i < len(outEdges); i++ {
		temp := outEdges[i]
		graph.Vertices[temp].deletedNeighbors++
	}
}

// contractNode
//
// vertex Vertex to be contracted
// contractID ID of contraction
//
func (graph *Graph) contractNode(vertex *Vertex, contractID int) {
	inEdges := vertex.inEdges
	inECost := vertex.inECost
	outEdges := vertex.outEdges
	outECost := vertex.outECost

	vertex.contracted = true

	inMax := 0.0
	outMax := 0.0

	inMaxVertex := -1
	outMaxVertex := -1

	graph.markNeighbors(vertex.inEdges, vertex.outEdges)

	for i := 0; i < len(inECost); i++ {
		if graph.Vertices[inEdges[i]].contracted {
			continue
		}
		if inMax < inECost[i] {
			inMax = inECost[i]
			inMaxVertex = inEdges[i]
		}
	}

	for i := 0; i < len(outECost); i++ {
		if graph.Vertices[outEdges[i]].contracted {
			continue
		}
		if outMax < outECost[i] {
			outMax = outECost[i]
			outMaxVertex = outEdges[i]
		}
	}

	max := inMax + outMax
	if inMaxVertex > 0 && outMaxVertex > 0 {
		if inMaxVertex == outMaxVertex {
			return
		}
	}
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i]
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inECost[i]
		graph.dijkstra(inVertex, max, contractID, int(i)) //finds the shortest distances from the inVertex to all the outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j]
			outcost := outECost[j]
			if graph.Vertices[outVertex].contracted {
				continue
			}
			summaryCost := incost + outcost
			if graph.Vertices[outVertex].distance.contractID != contractID || graph.Vertices[outVertex].distance.sourceID != int(i) || graph.Vertices[outVertex].distance.distance > summaryCost {
				if _, ok := graph.shortcuts[inVertex]; !ok {
					graph.shortcuts[inVertex] = make(map[int]*ShortcutInfo)
					graph.shortcuts[inVertex][outVertex] = &ShortcutInfo{
						ViaVertex: vertex.vertexNum,
						Cost:      summaryCost,
					}
				} else {
					graph.shortcuts[inVertex][outVertex] = &ShortcutInfo{
						ViaVertex: vertex.vertexNum,
						Cost:      summaryCost,
					}
				}
				graph.Vertices[inVertex].outEdges = append(graph.Vertices[inVertex].outEdges, outVertex)
				graph.Vertices[inVertex].outECost = append(graph.Vertices[inVertex].outECost, summaryCost)
				graph.Vertices[outVertex].inEdges = append(graph.Vertices[outVertex].inEdges, inVertex)
				graph.Vertices[outVertex].inECost = append(graph.Vertices[outVertex].inECost, summaryCost)
			}
		}
	}
}
