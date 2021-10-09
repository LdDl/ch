package ch

import (
	"container/heap"
)

// relaxEdgesBiForward Edge relaxation in a forward propagation
func (graph *Graph) relaxEdgesBiForward(vertex *simpleNode, forwQ *forwardHeap, prev map[int64]int64, queryDist []float64) {
	vertexList := graph.Vertices[vertex.id].outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := queryDist[vertex.id] + cost
			if queryDist[temp] > alt {
				queryDist[temp] = alt
				prev[temp] = vertex.id
				node := &simpleNode{
					id:        temp,
					queryDist: alt,
				}
				heap.Push(forwQ, node)
			}
		}
	}
}

// relaxEdgesBiForward Edge relaxation in a backward propagation
func (graph *Graph) relaxEdgesBiBackward(vertex *simpleNode, backwQ *backwardHeap, prev map[int64]int64, revDist []float64) {
	vertexList := graph.Vertices[vertex.id].inIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := revDist[vertex.id] + cost
			if revDist[temp] > alt {
				revDist[temp] = alt
				prev[temp] = vertex.id
				node := &simpleNode{
					id:          temp,
					revDistance: alt,
				}
				heap.Push(backwQ, node)
			}
		}
	}
}
