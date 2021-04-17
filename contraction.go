package ch

import (
	"container/heap"
)

// Preprocess Computes contraction hierarchies and returns node ordering
func (graph *Graph) Preprocess() []int64 {
	nodeOrdering := make([]int64, len(graph.Vertices))
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
		graph.contractNode(vertex, int64(extractNum-1))
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
func (graph *Graph) callNeighbors(inEdges []incidentEdge, outEdges []int64) {
	for i := 0; i < len(inEdges); i++ {
		temp := inEdges[i]
		graph.Vertices[temp.vertexID].delNeighbors++
	}
	for i := 0; i < len(outEdges); i++ {
		temp := outEdges[i]
		graph.Vertices[temp].delNeighbors++
	}
}

// ContractionPath
//
// ViaVertex - ID of vertex through which the contraction exists
// Cost - summary cost of path between two vertices
//
type ContractionPath struct {
	ViaVertex int64
	Cost      float64
}

// contractNode
//
// vertex Vertex to be contracted
// contractID ID of contraction
//
func (graph *Graph) contractNode(vertex *Vertex, contractID int64) {
	inEdges := vertex.inIncidentEdges
	outEdges := vertex.outEdges
	outECost := vertex.outECost

	vertex.contracted = true

	inMax := 0.0
	outMax := 0.0

	graph.callNeighbors(inEdges, vertex.outEdges)

	for i := 0; i < len(inEdges); i++ {
		if graph.Vertices[inEdges[i].vertexID].contracted {
			continue
		}
		if inMax < inEdges[i].cost {
			inMax = inEdges[i].cost
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
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inEdges[i].cost
		graph.dijkstra(inVertex, max, contractID, int64(i)) //finds the shortest distances from the inVertex to all the outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j]
			outcost := outECost[j]
			if graph.Vertices[outVertex].contracted {
				continue
			}
			summaryCost := incost + outcost
			if graph.Vertices[outVertex].distance.contractID != contractID || graph.Vertices[outVertex].distance.sourceID != int64(i) || graph.Vertices[outVertex].distance.distance > summaryCost {
				if _, ok := graph.contracts[inVertex]; !ok {
					graph.contracts[inVertex] = make(map[int64]*ContractionPath)
					graph.contracts[inVertex][outVertex] = &ContractionPath{
						ViaVertex: vertex.vertexNum,
						Cost:      summaryCost,
					}
				} else {
					graph.contracts[inVertex][outVertex] = &ContractionPath{
						ViaVertex: vertex.vertexNum,
						Cost:      summaryCost,
					}
				}
				graph.Vertices[inVertex].outEdges = append(graph.Vertices[inVertex].outEdges, outVertex)
				graph.Vertices[inVertex].outECost = append(graph.Vertices[inVertex].outECost, summaryCost)
				graph.Vertices[outVertex].inIncidentEdges = append(graph.Vertices[outVertex].inIncidentEdges, incidentEdge{inVertex, summaryCost})
			}
		}
	}
}
