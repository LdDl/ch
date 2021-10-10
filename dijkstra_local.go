package ch

const (
	Infinity = float64(^uint(0) >> 1)
	// Infinity = Infinity
)

// shortestPathsWithMaxCost Internal implementation of Dijkstra's algorithm to compute witness paths
func (graph *Graph) shortestPathsWithMaxCost(source int64, maxcost float64, previousOrderPos, neighborIndex int64) {
	// Heap to store traveled distance
	pqComparator := &distanceHeap{}
	pqComparator.Push(graph.Vertices[source])

	graph.Vertices[source].distance.distance = 0
	graph.Vertices[source].distance.previousOrderPos = previousOrderPos
	graph.Vertices[source].distance.sourceID = neighborIndex

	for pqComparator.Len() != 0 {
		vertex := pqComparator.Pop()
		// Do not consider any vertex has been excluded earlier
		if vertex.contracted {
			continue
		}
		if vertex.distance.distance > maxcost {
			return
		}
		// Edge relaxation
		vertexList := vertex.outIncidentEdges
		for i := range vertexList {
			temp := vertexList[i].vertexID
			cost := vertexList[i].weight
			tempPtr := graph.Vertices[temp]
			// Do not consider any vertex has been excluded earlier
			if tempPtr.contracted {
				continue
			}
			alt := vertex.distance.distance + cost
			if graph.checkID(vertex.vertexNum, temp) || tempPtr.distance.distance > alt {
				tempPtr.distance.distance = vertex.distance.distance + cost
				tempPtr.distance.previousOrderPos = previousOrderPos
				tempPtr.distance.sourceID = neighborIndex
				pqComparator.Push(tempPtr)
			}
		}
	}
}

// checkID Checks if both source's and target's contraction ID are not equal
func (graph *Graph) checkID(source, target int64) bool {
	s := graph.Vertices[source].distance
	t := graph.Vertices[target].distance
	return s.previousOrderPos != t.previousOrderPos || s.sourceID != t.sourceID
}
