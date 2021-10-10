package ch

const (
	Infinity = float64(^uint(0) >> 1)
	// Infinity = Infinity
)

// shortestPathsWithMaxCost Internal implementation of Dijkstra's algorithm to compute witness paths
func (graph *Graph) shortestPathsWithMaxCost(source int64, maxcost float64, previousOrderPos int64) {
	// Heap to store traveled distance
	pqComparator := &distanceHeap{}
	pqComparator.Push(graph.Vertices[source])

	// Instead of inializing distances to Infinity every single shortestPathsWithMaxCost(...) call we can do following
	// Set dist[source] -> 0 (as usual)
	graph.Vertices[source].distance.distance = 0
	// Set order position to previously contracted (excluded from graph) vertex
	graph.Vertices[source].distance.previousOrderPos = previousOrderPos
	// Set source to identifier of vertex for which shortestPathsWithMaxCost(...) has been called
	graph.Vertices[source].distance.previousSourceID = source

	for pqComparator.Len() != 0 {
		vertex := pqComparator.Pop()
		// Do not consider any vertex has been excluded earlier
		if vertex.contracted {
			continue
		}
		// Once a vertex is settled with a shortest path score greater than max cost, search stops.
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
			if tempPtr.distance.distance > alt ||
				vertex.distance.previousOrderPos != tempPtr.distance.previousOrderPos ||
				vertex.distance.previousSourceID != tempPtr.distance.previousSourceID {
				// Update new shortest distance
				tempPtr.distance.distance = vertex.distance.distance + cost
				tempPtr.distance.previousOrderPos = previousOrderPos
				tempPtr.distance.previousSourceID = source
				pqComparator.Push(tempPtr)
			}
		}
	}
}
