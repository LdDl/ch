package ch

import (
	"math"
)

// dijkstra_v2 Same as dijkstra() but with but with parallelism
func (graph *Graph) dijkstra_v2(source int64, maxcost float64, contractID, sourceID int64, threadID int) {
	graph.pqComparators[threadID] = &potatoHeapParallel{}
	// heap.Init(graph.pqComparators[threadID])
	// heap.Push(graph.pqComparators[threadID], graph.Vertices[source])

	// graph.Vertices[source].distance_v2[threadID].distance = 0
	// graph.Vertices[source].distance_v2[threadID].contractID = contractID
	// graph.Vertices[source].distance_v2[threadID].sourceID = sourceID

	// for graph.pqComparators[threadID].Len() != 0 {
	// 	vertex := heap.Pop(graph.pqComparators[threadID]).(*Vertex)
	// 	if graph.Vertices[vertex.vertexNum].distance_v2[threadID].distance > maxcost {
	// 		return
	// 	}
	// 	graph.relaxEdges_v2(vertex, contractID, sourceID, threadID)
	// }
}

// relaxEdges_v2 Same as relaxEdges() but with but with parallelism
func (graph *Graph) relaxEdges_v2(vertexInfo *Vertex, contractID, sourceID int64, threadID int) {
	vertexList := vertexInfo.outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		tempPtr := graph.Vertices[temp]
		// Skip shortcuts
		if tempPtr.contracted {
			continue
		}
		if graph.checkID_v2(vertexInfo.vertexNum, temp, threadID) || tempPtr.distance_v2[threadID].distance > vertexInfo.distance_v2[threadID].distance+cost {
			tempPtr.distance_v2[threadID].distance = vertexInfo.distance_v2[threadID].distance + cost
			tempPtr.distance_v2[threadID].contractID = contractID
			tempPtr.distance_v2[threadID].sourceID = sourceID
			// heap.Push(graph.pqComparators[threadID], tempPtr)
		}
	}

}

// checkID_v2 Same as checkID() but with but with parallelism
func (graph *Graph) checkID_v2(source, target int64, threadID int) bool {
	s := graph.Vertices[source].distance_v2[threadID]
	t := graph.Vertices[target].distance_v2[threadID]
	return s.contractID != t.contractID || s.sourceID != t.sourceID
}

// dijkstra_v3 Same as dijkstra() but with but with parallelism
func (graph *Graph) dijkstra_v3(source int64, maxcost float64, contractID, sourceID int64, threadID int) {
	// graph.pqComparators[threadID] = &potatoHeapParallel{}
	// graph.pqComparators[threadID].Push(graph.Vertices[source], threadID)

	graph.pqComparators[threadID] = &potatoHeapParallel{}
	PushParallel(graph.pqComparators[threadID], graph.Vertices[source], threadID)

	graph.Vertices[source].distance_v2[threadID].distance = 0
	graph.Vertices[source].distance_v2[threadID].contractID = contractID
	graph.Vertices[source].distance_v2[threadID].sourceID = sourceID

	for graph.pqComparators[threadID].Len() != 0 {
		vertex := graph.pqComparators[threadID].Pop()
		if graph.Vertices[vertex.vertexNum].distance_v2[threadID].distance > maxcost {
			return
		}
		graph.relaxEdges_v3(vertex, contractID, sourceID, threadID)
	}
}

// relaxEdges_v3 Same as relaxEdges() but with but with parallelism
func (graph *Graph) relaxEdges_v3(vertexInfo *Vertex, contractID, sourceID int64, threadID int) {
	vertexList := vertexInfo.outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		tempPtr := graph.Vertices[temp]
		// Skip shortcuts
		if tempPtr.contracted {
			continue
		}
		// if graph.checkID_v2(vertexInfo.vertexNum, temp, threadID) || tempPtr.distance_v2[threadID].distance > vertexInfo.distance_v2[threadID].distance+cost {
		// 	tempPtr.distance_v2[threadID].distance = vertexInfo.distance_v2[threadID].distance + cost
		// 	tempPtr.distance_v2[threadID].contractID = contractID
		// 	tempPtr.distance_v2[threadID].sourceID = sourceID
		// 	graph.pqComparators[threadID].Push(tempPtr, threadID)
		// }
		if graph.checkID_v2(vertexInfo.vertexNum, temp, threadID) || tempPtr.distance_v2[threadID].distance > vertexInfo.distance_v2[threadID].distance+cost {
			tempPtr.distance_v2[threadID].distance = vertexInfo.distance_v2[threadID].distance + cost
			tempPtr.distance_v2[threadID].contractID = contractID
			tempPtr.distance_v2[threadID].sourceID = sourceID
			// graph.pqComparators[threadID].Push(tempPtr, threadID)
			PushParallel(graph.pqComparators[threadID], tempPtr, threadID)
		}
	}

}

func (g *Graph) newCostTable(startNode int64) map[int64]float64 {
	costTable := make(map[int64]float64)
	costTable[startNode] = 0
	for _, i := range g.Vertices {
		if i.vertexNum != startNode {
			costTable[i.vertexNum] = math.MaxFloat64
		}
	}
	return costTable
}
