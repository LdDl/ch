package ch

// checkID Checks if both source's and target's contraction ID are not equal
func (graph *Graph) checkID(source, target int64) bool {
	s := graph.Vertices[source].distance
	t := graph.Vertices[target].distance
	return s.contractionID != t.contractionID || s.sourceID != t.sourceID
}

// dijkstra Internal dijkstra algorithm to compute contraction hierarchies
func (graph *Graph) dijkstra(source int64, maxcost float64, contractionID, sourceID int64) {
	graph.pqComparator = &distanceHeap{}
	graph.pqComparator.Push(graph.Vertices[source])

	graph.Vertices[source].distance.distance = 0
	graph.Vertices[source].distance.contractionID = contractionID
	graph.Vertices[source].distance.sourceID = sourceID

	for graph.pqComparator.Len() != 0 {
		vertex := graph.pqComparator.Pop()
		if vertex.distance.distance > maxcost {
			return
		}
		// Edge relaxation
		vertexList := vertex.outIncidentEdges
		for i := 0; i < len(vertexList); i++ {
			temp := vertexList[i].vertexID
			cost := vertexList[i].cost
			tempPtr := graph.Vertices[temp]
			// Skip contracted vertices
			if tempPtr.contracted {
				continue
			}
			if graph.checkID(vertex.vertexNum, temp) || tempPtr.distance.distance > vertex.distance.distance+cost {
				tempPtr.distance.distance = vertex.distance.distance + cost
				tempPtr.distance.contractionID = contractionID
				tempPtr.distance.sourceID = sourceID
				graph.pqComparator.Push(tempPtr)
			}
		}
	}
}
