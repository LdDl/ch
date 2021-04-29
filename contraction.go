package ch

import (
	"container/heap"
	"fmt"
	"time"
)

const DEBUG_PREPROCESSING = false

// Preprocess Computes contraction hierarchies and returns node ordering
func (graph *Graph) Preprocess() []int64 {
	nodeOrdering := make([]int64, len(graph.Vertices))
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
		if graph.pqImportance.Len() != 0 && vertex.importance > graph.pqImportance.Peek().importance {
			graph.pqImportance.Push(vertex)
			continue
		}

		nodeOrdering[extractNum] = vertex.vertexNum
		vertex.orderPos = extractNum
		extractNum = extractNum + 1
		graph.contractNode(vertex, int64(extractNum-1))

		if DEBUG_PREPROCESSING {
			if iter > 0 && graph.pqImportance.Len()%1000 == 0 {
				fmt.Printf("Contraction Order: %d / %d, Remain vertices in heap: %d. Currect shortcuts num: %d Time: %v\n", extractNum, len(graph.Vertices), graph.pqImportance.Len(), graph.shortcutsNum(), time.Now())
			}
		}
	}
	return nodeOrdering
}

// markNeighbors
//
// inEdges Incoming edges from vertex
// outEdges Outcoming edges from vertex
//
func (graph *Graph) markNeighbors(inEdges, outEdges []incidentEdge) {
	for i := 0; i < len(inEdges); i++ {
		temp := inEdges[i]
		graph.Vertices[temp.vertexID].delNeighbors++
	}
	for i := 0; i < len(outEdges); i++ {
		temp := outEdges[i]
		graph.Vertices[temp.vertexID].delNeighbors++
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
				if _, ok := graph.shortcuts[inVertex]; !ok {
					// If there is no such shortcut add one.
					graph.shortcuts[inVertex] = make(map[int64]*ContractionPath)
					graph.shortcuts[inVertex][outVertex] = &ContractionPath{
						ViaVertex: vertex.vertexNum,
						Cost:      summaryCost,
					}
					graph.Vertices[inVertex].outIncidentEdges = append(graph.Vertices[inVertex].outIncidentEdges, incidentEdge{outVertex, summaryCost})
					graph.Vertices[outVertex].inIncidentEdges = append(graph.Vertices[outVertex].inIncidentEdges, incidentEdge{inVertex, summaryCost})
				} else {
					if v, ok := graph.shortcuts[inVertex][outVertex]; ok {
						// If shortcut already exists
						// we should check if the middle vertex is still the same
						if v.ViaVertex == vertex.vertexNum {
							// If middle vertex is still the same then change cost of shortcut only
							graph.shortcuts[inVertex][outVertex].Cost = summaryCost
							bk1 := graph.Vertices[inVertex].updateOutIncidentEdge(outVertex, summaryCost)
							if !bk1 {
								panic("Should not happen [1]. Can't update outcoming incident edge")
							}
							bk2 := graph.Vertices[outVertex].updateInIncidentEdge(inVertex, summaryCost)
							if !bk2 {
								panic("Should not happen [2]/ Can't update incoming incident edge")
							}
						} else {
							// If middle vertex is not optimal for shortcut then change both vertex ID and cost
							graph.shortcuts[inVertex][outVertex].ViaVertex = vertex.vertexNum
							graph.shortcuts[inVertex][outVertex].Cost = summaryCost

							dk1 := graph.Vertices[inVertex].deleteOutIncidentEdge(outVertex)
							if !dk1 {
								panic("Should not happen [3]. Can't delete outcoming incident edge")
							}
							dk2 := graph.Vertices[outVertex].deleteInIncidentEdge(inVertex)
							if !dk2 {
								panic("Should not happen [4]. Can't delete incoming incident edge")
							}
							graph.Vertices[inVertex].outIncidentEdges = append(graph.Vertices[inVertex].outIncidentEdges, incidentEdge{outVertex, summaryCost})
							graph.Vertices[outVertex].inIncidentEdges = append(graph.Vertices[outVertex].inIncidentEdges, incidentEdge{inVertex, summaryCost})
						}
					} else {
						graph.shortcuts[inVertex][outVertex] = &ContractionPath{
							ViaVertex: vertex.vertexNum,
							Cost:      summaryCost,
						}
						graph.Vertices[inVertex].outIncidentEdges = append(graph.Vertices[inVertex].outIncidentEdges, incidentEdge{outVertex, summaryCost})
						graph.Vertices[outVertex].inIncidentEdges = append(graph.Vertices[outVertex].inIncidentEdges, incidentEdge{inVertex, summaryCost})
					}
				}
			}
		}
	}
}
