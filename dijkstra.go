package ch

// checkID Checks if both source's and target's contraction ID are not equal
func (graph *Graph) checkID(source, target int64) bool {
	s := graph.Vertices[source].distance
	t := graph.Vertices[target].distance
	return s.contracttionID != t.contracttionID || s.sourceID != t.sourceID
}

// relaxEdges Edge relaxation
func (graph *Graph) relaxEdges(vertexInfo *Vertex, contracttionID, sourceID int64) {
	vertexList := vertexInfo.outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		tempPtr := graph.Vertices[temp]
		// Skip shortcuts
		if tempPtr.contracted {
			continue
		}
		if graph.checkID(vertexInfo.vertexNum, temp) || tempPtr.distance.distance > vertexInfo.distance.distance+cost {
			tempPtr.distance.distance = vertexInfo.distance.distance + cost
			tempPtr.distance.contracttionID = contracttionID
			tempPtr.distance.sourceID = sourceID
			graph.pqComparator.Push(tempPtr)
			// graph.pqComparator.Push(tempPtr)
		}
	}
}

// dijkstra Internal dijkstra algorithm to compute contraction hierarchies
func (graph *Graph) dijkstra(source int64, maxcost float64, contracttionID, sourceID int64) {
	graph.pqComparator = &distanceHeap{}
	graph.pqComparator.Push(graph.Vertices[source])

	graph.Vertices[source].distance.distance = 0
	graph.Vertices[source].distance.contracttionID = contracttionID
	graph.Vertices[source].distance.sourceID = sourceID

	for graph.pqComparator.Len() != 0 {
		vertex := graph.pqComparator.Pop()
		if vertex.distance.distance > maxcost {
			return
		}
		graph.relaxEdges(vertex, contracttionID, sourceID)
	}
}
