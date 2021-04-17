package ch

import "container/heap"

// checkID Checks if both source's and target's contraction ID are not equal
func (graph *Graph) checkID(source, target int64) bool {
	if graph.Vertices[source].distance.contractID != graph.Vertices[target].distance.contractID || graph.Vertices[source].distance.sourceID != graph.Vertices[target].distance.sourceID {
		return true
	}
	return false
}

// relaxEdges Edge relaxation
func (graph *Graph) relaxEdges(vertex, contractID, sourceID int64) {
	vertexList := graph.Vertices[vertex].outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		if graph.Vertices[temp].contracted {
			continue
		}

		if graph.checkID(vertex, temp) || graph.Vertices[temp].distance.distance > graph.Vertices[vertex].distance.distance+cost {
			graph.Vertices[temp].distance.distance = graph.Vertices[vertex].distance.distance + cost
			graph.Vertices[temp].distance.contractID = contractID
			graph.Vertices[temp].distance.sourceID = sourceID
			heap.Push(graph.pqComparator, graph.Vertices[temp])
		}

	}
}

// dijkstra Internal dijkstra algorithm to compute contraction hierarchies
func (graph *Graph) dijkstra(source int64, maxcost float64, contractID, sourceID int64) {
	graph.pqComparator = &distanceHeap{}
	heap.Init(graph.pqComparator)
	heap.Push(graph.pqComparator, graph.Vertices[source])

	graph.Vertices[source].distance.distance = 0
	graph.Vertices[source].distance.contractID = contractID
	graph.Vertices[source].distance.sourceID = sourceID

	for graph.pqComparator.Len() != 0 {
		vertex := heap.Pop(graph.pqComparator).(*Vertex)
		if vertex.distance.distance > maxcost {
			return
		}
		graph.relaxEdges(vertex.vertexNum, contractID, sourceID)
	}

}
