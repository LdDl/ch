package ch

// incidentEdge incident edge for certain vertex
type incidentEdge struct {
	vertexID int64
	weight   float64
}

// addInIncidentEdge Adds incident edge's to pool of "incoming" edges of given vertex.
// Just an alias to append() function with invalidation of bidirected cache
// incomingVertexID - Library defined ID of vertex
// weight - Travel cost of incoming edge
func (vertex *Vertex) addInIncidentEdge(incomingVertexID int64, weight float64) {
	vertex.inIncidentEdges = append(vertex.inIncidentEdges, incidentEdge{incomingVertexID, weight})
	vertex.bidirectedCached = false
}

// addOutIncidentEdge Adds incident edge's to pool of "outcoming" edges of given vertex.
// Just an alias to append() function with invalidation of bidirected cache
// outcomingVertexID - Library defined ID of vertex
// weight - Travel cost of outcoming edge
func (vertex *Vertex) addOutIncidentEdge(outcomingVertexID int64, weight float64) {
	vertex.outIncidentEdges = append(vertex.outIncidentEdges, incidentEdge{outcomingVertexID, weight})
	vertex.bidirectedCached = false
}

// findInIncidentEdge Returns index of incoming incident edge by vertex ID
// If incoming incident edge is not found then this function returns -1
func (vertex *Vertex) findInIncidentEdge(vertexID int64) int {
	for i := range vertex.inIncidentEdges {
		if vertex.inIncidentEdges[i].vertexID == vertexID {
			return i
		}
	}
	return -1
}

// findOutIncidentEdge Returns index of outcoming incident edge by vertex ID on the other side of that edge
// If outcoming incident edge is not found then this function returns -1
func (vertex *Vertex) findOutIncidentEdge(vertexID int64) int {
	for i := range vertex.outIncidentEdges {
		if vertex.outIncidentEdges[i].vertexID == vertexID {
			return i
		}
	}
	return -1
}

// updateInIncidentEdge Updates incoming incident edge's cost by vertex ID on the other side of that edge
// If operation is not successful then this function returns False
func (vertex *Vertex) updateInIncidentEdge(vertexID int64, weight float64) bool {
	idx := vertex.findInIncidentEdge(vertexID)
	if idx < 0 {
		return false
	}
	vertex.inIncidentEdges[idx].weight = weight
	return true
}

// updateOutIncidentEdge Updates outcoming incident edge's cost by vertex ID on the other side of that edge
// If operation is not successful then this function returns False
func (vertex *Vertex) updateOutIncidentEdge(vertexID int64, weight float64) bool {
	idx := vertex.findOutIncidentEdge(vertexID)
	if idx < 0 {
		return false
	}
	vertex.outIncidentEdges[idx].weight = weight
	return true
}
