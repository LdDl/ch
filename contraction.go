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
	for graph.pqImportance.Len() != 0 {
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
		graph.contractNode(vertex)
		if DEBUG_PREPROCESSING {
			if extractNum > 0 && graph.pqImportance.Len()%1000 == 0 {
				fmt.Printf("Contraction Order: %d / %d, Remain vertices in heap: %d. Currect shortcuts num: %d Time: %v\n", extractNum, len(graph.Vertices), graph.pqImportance.Len(), graph.shortcutsNum(), time.Now())
			}
		}
		extractNum++
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
// ViaVertex - ID of vertex through which the shortcut exists
// Cost - summary cost of path between two vertices
//
type ContractionPath struct {
	ViaVertex int64
	Cost      float64
}

// contractNode
//
// vertex Vertex to be contracted
//
func (graph *Graph) contractNode(vertex *Vertex) {
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

	contractionID := int64(vertex.orderPos - 1)
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inEdges[i].cost
		graph.dijkstra(inVertex, max, contractionID, int64(i)) // Finds the shortest distances from the inVertex to all outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j].vertexID
			outcost := outEdges[j].cost
			outVertexPtr := graph.Vertices[outVertex]

			if outVertexPtr.contracted {
				continue
			}
			summaryCost := incost + outcost
			if outVertexPtr.distance.contractID != contractionID || outVertexPtr.distance.sourceID != int64(i) || outVertexPtr.distance.distance > summaryCost {
				graph.createOrUpdateShortcut(inVertex, outVertex, vertex.vertexNum, summaryCost)
			}
		}
	}
}

// createOrUpdateShortcut Creates (or updates: it depends on conditions) shortcut
//
// fromVertex - Library defined ID of source vertex where shortcut starts from
// fromVertex - Library defined ID of target vertex where shortcut leads to
// viaVertex - Library defined ID of vertex through which the shortcut exists
// summaryCost - Travel path of a shortcut
//
func (graph *Graph) createOrUpdateShortcut(fromVertex, toVertex, viaVertex int64, summaryCost float64) {
	if _, ok := graph.shortcuts[fromVertex]; !ok {
		// If there is no such shortcut then add one.
		graph.shortcuts[fromVertex] = make(map[int64]*ContractionPath)
		graph.shortcuts[fromVertex][toVertex] = &ContractionPath{
			ViaVertex: viaVertex,
			Cost:      summaryCost,
		}
		graph.Vertices[fromVertex].addOutIncidentEdge(toVertex, summaryCost)
		graph.Vertices[toVertex].addInIncidentEdge(fromVertex, summaryCost)
	} else {
		if v, ok := graph.shortcuts[fromVertex][toVertex]; ok {
			// If shortcut already exists
			// we should check if the middle vertex is still the same
			if v.ViaVertex == viaVertex {
				// If middle vertex is still the same then change cost of shortcut only [Additional conditional: previous estimated cost is less than current one]
				if summaryCost < graph.shortcuts[fromVertex][toVertex].Cost {
					graph.shortcuts[fromVertex][toVertex].Cost = summaryCost
					updatedOutSuccess := graph.Vertices[fromVertex].updateOutIncidentEdge(toVertex, summaryCost)
					if !updatedOutSuccess {
						panic(fmt.Sprintf("Should not happen [1]. Can't update outcoming incident edge. %d has no common edge with %d", fromVertex, toVertex))
					}
					updatedInSuccess := graph.Vertices[toVertex].updateInIncidentEdge(fromVertex, summaryCost)
					if !updatedInSuccess {
						panic(fmt.Sprintf("Should not happen [2]. Can't update incoming incident edge. %d has no common edge with %d", toVertex, fromVertex))
					}
					graph.Vertices[fromVertex].addOutIncidentEdge(toVertex, summaryCost)
					graph.Vertices[toVertex].addInIncidentEdge(fromVertex, summaryCost)
				}
			} else {
				// If middle vertex is not optimal for shortcut then change both vertex ID and cost [Additional conditional: previous estimated cost is less than current one]
				if summaryCost < graph.shortcuts[fromVertex][toVertex].Cost {
					graph.shortcuts[fromVertex][toVertex].ViaVertex = viaVertex
					graph.shortcuts[fromVertex][toVertex].Cost = summaryCost
					deletedOutSuccess := graph.Vertices[fromVertex].deleteOutIncidentEdge(toVertex)
					if !deletedOutSuccess {
						panic(fmt.Sprintf("Should not happen [3]. Can't delete outcoming incident edge. %d has no common edge with %d", fromVertex, toVertex))
					}
					deletedInSuccess := graph.Vertices[toVertex].deleteInIncidentEdge(fromVertex)
					if !deletedInSuccess {
						panic(fmt.Sprintf("Should not happen [4]. Can't delete incoming incident edge. %d has no common edge with %d", toVertex, fromVertex))
					}
					graph.Vertices[fromVertex].addOutIncidentEdge(toVertex, summaryCost)
					graph.Vertices[toVertex].addInIncidentEdge(fromVertex, summaryCost)
				}
			}
		} else {
			graph.shortcuts[fromVertex][toVertex] = &ContractionPath{
				ViaVertex: viaVertex,
				Cost:      summaryCost,
			}
			graph.Vertices[fromVertex].addOutIncidentEdge(toVertex, summaryCost)
			graph.Vertices[toVertex].addInIncidentEdge(fromVertex, summaryCost)
		}
	}
}
