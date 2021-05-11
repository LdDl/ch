package ch

import (
	"container/heap"
)

// checkID Checks if both source's and target's contraction ID are not equal
func (graph *Graph) checkID(source, target int64) bool {
	return graph.Vertices[source].distance.contractID != graph.Vertices[target].distance.contractID || graph.Vertices[source].distance.sourceID != graph.Vertices[target].distance.sourceID
}

// relaxEdges Edge relaxation
func (graph *Graph) relaxEdges(vertex, contractID, sourceID int64) {
	vertexList := graph.Vertices[vertex].outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		// Skip shortcuts
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

// dijkstra_v2 Same as dijkstra() but with but with parallelism
func (graph *Graph) dijkstra_v2(source int64, maxcost float64, contractID, sourceID int64, threadID int) {
	graph.pqComparators[threadID] = &distanceHeap{}
	heap.Init(graph.pqComparators[threadID])
	heap.Push(graph.pqComparators[threadID], graph.Vertices[source])

	distances := make(map[int64]*Distance)
	for i := range graph.Vertices {
		distances[graph.Vertices[i].vertexNum] = NewDistance()
	}

	distances[source].distance = 0
	distances[source].contractID = contractID
	distances[source].sourceID = sourceID

	for graph.pqComparators[threadID].Len() != 0 {
		vertex := heap.Pop(graph.pqComparators[threadID]).(*Vertex)
		if distances[vertex.vertexNum].distance > maxcost {
			return
		}
		graph.relaxEdges_v2(graph.pqComparators[threadID], distances, vertex.vertexNum, contractID, sourceID)
	}
}

// relaxEdges_v2 Same as relaxEdges() but with but with parallelism
func (graph *Graph) relaxEdges_v2(pqComparator *distanceHeap, distances map[int64]*Distance, vertex, contractID, sourceID int64) {
	vertexList := graph.Vertices[vertex].outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		// Skip shortcuts
		if graph.Vertices[temp].contracted {
			continue
		}
		if graph.checkID(vertex, temp) || distances[temp].distance > distances[vertex].distance+cost {
			distances[temp].distance = distances[vertex].distance + cost
			distances[temp].contractID = contractID
			distances[temp].sourceID = sourceID
			heap.Push(pqComparator, graph.Vertices[temp])
		}
	}
}
