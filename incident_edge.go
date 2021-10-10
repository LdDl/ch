package ch

// incidentEdge incident edge for certain vertex
type incidentEdge struct {
	vertexID int64
	cost     float64
}

// addInIncidentEdge Adds incident edge's to pool of "incoming" edges of given vertex.
// Just an alias to append() function
// incomingVertexID - Library defined ID of vertex
// cost - Travel cost of incoming edge
func (vertex *Vertex) addInIncidentEdge(incomingVertexID int64, cost float64) {
	vertex.inIncidentEdges = append(vertex.inIncidentEdges, incidentEdge{incomingVertexID, cost})
}

// addOutIncidentEdge Adds incident edge's to pool of "outcoming" edges of given vertex.
// Just an alias to append() function
// outcomingVertexID - Library defined ID of vertex
// cost - Travel cost of outcoming edge
func (vertex *Vertex) addOutIncidentEdge(outcomingVertexID int64, cost float64) {
	vertex.outIncidentEdges = append(vertex.outIncidentEdges, incidentEdge{outcomingVertexID, cost})
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
func (vertex *Vertex) updateInIncidentEdge(vertexID int64, cost float64) bool {
	idx := vertex.findInIncidentEdge(vertexID)
	if idx < 0 {
		return false
	}
	vertex.inIncidentEdges[idx].cost = cost
	return true
}

// updateOutIncidentEdge Updates outcoming incident edge's cost by vertex ID on the other side of that edge
// If operation is not successful then this function returns False
func (vertex *Vertex) updateOutIncidentEdge(vertexID int64, cost float64) bool {
	idx := vertex.findOutIncidentEdge(vertexID)
	if idx < 0 {
		return false
	}
	vertex.outIncidentEdges[idx].cost = cost
	return true
}
