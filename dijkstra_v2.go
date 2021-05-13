package ch

import (
	"container/heap"
)

// dijkstra_v2 Same as dijkstra() but with but with parallelism
func (graph *Graph) dijkstra_v2(source int64, maxcost float64, contractID, sourceID int64, threadID int) {
	graph.pqComparators[threadID] = &distanceHeap{}
	heap.Init(graph.pqComparators[threadID])
	heap.Push(graph.pqComparators[threadID], graph.Vertices[source])

	graph.Vertices[source].distance_v2[threadID].distance = 0
	graph.Vertices[source].distance_v2[threadID].contractID = contractID
	graph.Vertices[source].distance_v2[threadID].sourceID = sourceID

	for graph.pqComparators[threadID].Len() != 0 {
		vertex := heap.Pop(graph.pqComparators[threadID]).(*Vertex)
		if graph.Vertices[vertex.vertexNum].distance_v2[threadID].distance > maxcost {
			return
		}
		graph.relaxEdges_v2(graph.pqComparators[threadID], vertex.vertexNum, contractID, sourceID, threadID)
	}
}

// relaxEdges_v2 Same as relaxEdges() but with but with parallelism
func (graph *Graph) relaxEdges_v2(pqComparator *distanceHeap, vertex, contractID, sourceID int64, threadID int) {
	vertexList := graph.Vertices[vertex].outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		// Skip shortcuts
		if graph.Vertices[temp].contracted {
			continue
		}
		if graph.Vertices[temp].distance_v2[threadID].distance > graph.Vertices[vertex].distance_v2[threadID].distance+cost {
			graph.Vertices[temp].distance_v2[threadID].distance = graph.Vertices[vertex].distance_v2[threadID].distance + cost
			graph.Vertices[temp].distance_v2[threadID].contractID = contractID
			graph.Vertices[temp].distance_v2[threadID].sourceID = sourceID
			heap.Push(pqComparator, graph.Vertices[temp])
		}
	}

}

// checkID_v2 Same as checkID() but with but with parallelism
func (graph *Graph) checkID_v2(source, target int64, threadID int) bool {
	return graph.Vertices[source].distance_v2[threadID].contractID != graph.Vertices[target].distance_v2[threadID].contractID || graph.Vertices[source].distance_v2[threadID].sourceID != graph.Vertices[target].distance_v2[threadID].sourceID
}
