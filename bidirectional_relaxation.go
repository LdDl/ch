package ch

import (
	"container/heap"
)

// relaxEdgesBiForward Edge relaxation in a forward propagation
func (graph *Graph) relaxEdgesBiForward(vertex *bidirectionalVertex, forwQ *forwardHeap, prev map[int64]int64, queryDist []float64) {
	vertexList := graph.Vertices[vertex.id].outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := queryDist[vertex.id] + cost
			if queryDist[temp] > alt {
				queryDist[temp] = alt
				prev[temp] = vertex.id
				node := &bidirectionalVertex{
					id:        temp,
					queryDist: alt,
				}
				heap.Push(forwQ, node)
			}
		}
	}
}

// relaxEdgesBiForward Edge relaxation in a backward propagation
func (graph *Graph) relaxEdgesBiBackward(vertex *bidirectionalVertex, backwQ *backwardHeap, prev map[int64]int64, revQueryDist []float64) {
	vertexList := graph.Vertices[vertex.id].inIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := revQueryDist[vertex.id] + cost
			if revQueryDist[temp] > alt {
				revQueryDist[temp] = alt
				prev[temp] = vertex.id
				node := &bidirectionalVertex{
					id:               temp,
					revQueryDistance: alt,
				}
				heap.Push(backwQ, node)
			}
		}
	}
}
