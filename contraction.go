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
		vertex := heap.Pop(graph.pqImportance).(*Vertex)
		vertex.computeImportance()
		if graph.pqImportance.Len() != 0 && vertex.importance > graph.pqImportance.Peek().(*Vertex).importance {
			graph.pqImportance.Push(vertex)
			continue
		}
		nodeOrdering[extractNum] = vertex.vertexNum
		vertex.orderPos = extractNum
		extractNum = extractNum + 1
		graph.contractNode(vertex, extractNum-1)
		// fmt.Printf(
		// 	"Contraction of vertex: %v (label %v, order %v, dist %v) | Contraction ID: %v (%v) | Done %v / %v | HeapLength: %v\n",
		// 	vertex.vertexNum, vertex.Label, vertex.orderPos, vertex.distance.distance, vertex.distance.contractID, vertex.distance.sourceID, iter, len(graph.Vertices), graph.pqImportance.Len(),
		// )
	}
	return nodeOrdering
}

// callNeighbors
//
// inEdges Incoming edges from vertex
// outEdges Outcoming edges from vertex
//
func (graph *Graph) callNeighbors(inEdges, outEdges []int) {
	for i := 0; i < len(inEdges); i++ {
		temp := inEdges[i]
		graph.Vertices[temp].delNeighbors++
	}
	for i := 0; i < len(outEdges); i++ {
		temp := outEdges[i]
		graph.Vertices[temp].delNeighbors++
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

	graph.callNeighbors(vertex.inEdges, vertex.outEdges)

	for i := 0; i < len(inECost); i++ {
		if graph.Vertices[inEdges[i]].contracted {
			continue
		}
		if inMax < inECost[i] {
			inMax = inECost[i]
		}
	}

	for i := 0; i < len(outECost); i++ {
		if graph.Vertices[outEdges[i]].contracted {
			continue
		}
		if outMax < outECost[i] {
			outMax = outECost[i]
		}
	}

	max := inMax + outMax

	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i]
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inECost[i]
		graph.dijkstra(inVertex, max, contractID, i) //finds the shortest distances from the inVertex to all the outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j]
			outcost := outECost[j]
			if graph.Vertices[outVertex].contracted {
				continue
			}
			if graph.Vertices[outVertex].distance.contractID != contractID || graph.Vertices[outVertex].distance.sourceID != i || graph.Vertices[outVertex].distance.distance > incost+outcost {
				if _, ok := graph.contracts[inVertex]; !ok {
					graph.contracts[inVertex] = make(map[int]int)
					graph.contracts[inVertex][outVertex] = vertex.vertexNum
				} else {
					graph.contracts[inVertex][outVertex] = vertex.vertexNum
				}
				graph.Vertices[inVertex].outEdges = append(graph.Vertices[inVertex].outEdges, outVertex)
				graph.Vertices[inVertex].outECost = append(graph.Vertices[inVertex].outECost, incost+outcost)
				graph.Vertices[outVertex].inEdges = append(graph.Vertices[outVertex].inEdges, inVertex)
				graph.Vertices[outVertex].inECost = append(graph.Vertices[outVertex].inECost, incost+outcost)
			}
		}
	}
}
