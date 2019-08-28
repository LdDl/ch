package ch

import (
	"container/heap"
)

// relaxEdgesBiForward Edge relaxation in a forward propagation
func (graph *Graph) relaxEdgesBiForward(vertex *simpleNode, forwQ *forwardPropagationHeap, prev map[int]int, queryDist []float64, prevReverse []float64) {
	vertexList := graph.Vertices[vertex.id].outEdges
	costList := graph.Vertices[vertex.id].outECost
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i]
		cost := costList[i]
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			//if(graph[vertex].distance.forwqueryId != graph[temp].distance.forwqueryId || graph[temp].distance.queryDist > graph[vertex].distance.queryDist + cost){
			alt := queryDist[vertex.id] + cost
			if queryDist[temp] > alt {
				queryDist[temp] = alt
				prev[temp] = vertex.id
				node := simpleNode{
					id:          temp,
					queryDist:   alt,
					revDistance: prevReverse[temp],
				}
				heap.Push(forwQ, node)
			}
		}
	}
}

// relaxEdgesBiForward Edge relaxation in a backward propagation
func (graph *Graph) relaxEdgesBiBackward(vertex *simpleNode, backwQ *backwardPropagationHeap, prev map[int]int, queryDist []float64, prevReverse []float64) {
	vertexList := graph.Vertices[vertex.id].inEdges
	costList := graph.Vertices[vertex.id].inECost
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i]
		cost := costList[i]
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := prevReverse[vertex.id] + cost
			if prevReverse[temp] > alt {
				prevReverse[temp] = alt
				prev[temp] = vertex.id
				node := simpleNode{
					id:          temp,
					revDistance: alt,
					queryDist:   queryDist[temp],
				}
				heap.Push(backwQ, node)
			}
		}
	}
}
