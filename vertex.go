package ch

import (
	"math"
)

// Vertex All information about vertex
type Vertex struct {
	vertexNum int64
	Label     int64

	inIncidentEdges  []incidentEdge
	outIncidentEdges []incidentEdge

	orderPos         int
	contracted       bool
	distance         *Distance
	incidentEdgesNum int
	edgeDiff         int
	delNeighbors     int
	shortcutCover    int
	importance       int
}

// OrderPos Returns order position (in terms of contraction hierarchies) of vertex
func (vertex *Vertex) OrderPos() int {
	return vertex.orderPos
}

// SetOrderPos Sets order position (in terms of contraction hierarchies) for vertex
func (vertex *Vertex) SetOrderPos(orderPos int) {
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

// NewVertex Create vertex with vertex number
func NewVertex(vertexNum int64) *Vertex {
	return &Vertex{
		vertexNum:    vertexNum,
		delNeighbors: 0,
		distance:     NewDistance(),
		contracted:   false,
	}
}

// computeImportance Update importance of vertex
func (vertex *Vertex) computeImportance() {
	// Worst possible shortcuts number throught the vertex is: NumWorstShortcuts = NumIncomingEdges*NumOutcomingEdges
	vertex.shortcutCover = len(vertex.inIncidentEdges) * len(vertex.outIncidentEdges)
	// Number of total incident edges is: NumIncomingEdges+NumOutcomingEdges
	vertex.incidentEdgesNum = len(vertex.inIncidentEdges) + len(vertex.outIncidentEdges)
	// Edge difference is between NumWorstShortcuts and TotalIncidentEdgesNum
	vertex.edgeDiff = vertex.shortcutCover - vertex.incidentEdgesNum
	// [+] Spatial diversity heuristic: for each vertex maintain a count of the number of neighbors that have already been contracted [vertex.delNeighbors], and add this to the summary importance
	// note: the more neighbours have already been contracted, the later this vertex will be contracted in further.
	// [+] Bidirection edges heuristic: for each vertex check how many bidirected incident edges vertex has.
	// note: the more bidirected incident edges == less important vertex is.
	vertex.importance = vertex.edgeDiff*14 + vertex.incidentEdgesNum*25 + vertex.delNeighbors*10 - vertex.bidirectedEdges()
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
type Distance struct {
	contractionID int64
	sourceID      int64
	distance      float64
	queryDist     float64
	revQueryDist  float64
}

// NewDistance Constructor for Distance
func NewDistance() *Distance {
	return &Distance{
		contractionID: -1,
		sourceID:      -1,
		distance:      math.MaxFloat64,
		revQueryDist:  math.MaxFloat64,
		queryDist:     math.MaxFloat64,
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
