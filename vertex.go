package ch

// Vertex All information about vertex
type Vertex struct {
	distance         *Distance
	inIncidentEdges  []*incidentEdge
	outIncidentEdges []*incidentEdge

	vertexNum int64
	Label     int64

	orderPos     int64
	delNeighbors int
	importance   int
	contracted   bool
}

// OrderPos Returns order position (in terms of contraction hierarchies) of vertex
func (vertex *Vertex) OrderPos() int64 {
	return vertex.orderPos
}

// SetOrderPos Sets order position (in terms of contraction hierarchies) for vertex
func (vertex *Vertex) SetOrderPos(orderPos int64) {
	vertex.orderPos = orderPos
}

// Importance Returns importance (in terms of contraction hierarchies) of vertex
func (vertex *Vertex) Importance() int {
	return vertex.importance
}

// SetImportance Sets order position (in terms of contraction hierarchies) for vertex
func (vertex *Vertex) SetImportance(importance int) {
	vertex.importance = importance
}

// MakeVertex Create vertex with label
func MakeVertex(label int64) *Vertex {
	return &Vertex{
		Label:        label,
		delNeighbors: 0,
		distance:     NewDistance(),
		contracted:   false,
	}
}

// computeImportance Update importance of vertex
func (vertex *Vertex) computeImportance() {
	// Worst possible shortcuts number through the vertex is: NumWorstShortcuts = NumIncomingEdges*NumOutcomingEdges
	shortcutCover := len(vertex.inIncidentEdges) * len(vertex.outIncidentEdges)
	// Number of total incident edges is: NumIncomingEdges+NumOutcomingEdges
	incidentEdgesNum := len(vertex.inIncidentEdges) + len(vertex.outIncidentEdges)
	// Edge difference is between NumWorstShortcuts and TotalIncidentEdgesNum
	edgeDiff := shortcutCover - incidentEdgesNum
	// [+] Spatial diversity heuristic: for each vertex maintain a count of the number of neighbors that have already been contracted [vertex.delNeighbors], and add this to the summary importance
	// note: the more neighbours have already been contracted, the later this vertex will be contracted in further.
	// [+] Bidirection edges heuristic: for each vertex check how many bidirected incident edges vertex has.
	// note: the more bidirected incident edges == less important vertex is.
	vertex.importance = edgeDiff + incidentEdgesNum + vertex.delNeighbors - vertex.bidirectedEdges()
}

// bidirectedEdges Number of bidirected edges
func (vertex *Vertex) bidirectedEdges() int {
	hash := make(map[int64]bool)
	for _, e := range vertex.inIncidentEdges {
		hash[e.vertexID] = true
	}
	ans := 0
	for _, e := range vertex.outIncidentEdges {
		if hash[e.vertexID] {
			ans++
		}
	}
	return ans
}

// Distance Information about contraction between source vertex and contraction vertex
// distance - used in Dijkstra local searches
// previousOrderPos - previous contraction order (shortest path tree)
// previousSourceID - previously found source vertex (shortest path tree)
type Distance struct {
	previousOrderPos int64
	previousSourceID int64
	distance         float64
}

// NewDistance Constructor for Distance
func NewDistance() *Distance {
	return &Distance{
		previousOrderPos: -1,
		previousSourceID: -1,
		distance:         Infinity,
	}
}

// FindVertex Returns index of vertex in graph
//
// labelExternal - User defined ID of vertex
// If vertex is not found then returns (-1; false)
//
func (graph *Graph) FindVertex(labelExternal int64) (idx int64, ok bool) {
	idx, ok = graph.mapping[labelExternal]
	if !ok {
		return -1, ok
	}
	return
}
