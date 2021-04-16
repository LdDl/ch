package ch

import (
	"math"
)

// Vertex All information about vertex
//
//
// shortcutCover Number of shortcuts that would have to be added if vertex were contracted
// incidentEdges Number of edges incident to vertex
//
type Vertex struct {
	vertexNum int
	Label     int

	inEdges  []int
	inECost  []float64
	outEdges []int
	outECost []float64

	orderPos         int
	contracted       bool
	distance         *Distance
	edgeDiff         int
	deletedNeighbors int
	shortcutCover    int
	incidentEdges    int
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
func MakeVertex(label int) *Vertex {
	return &Vertex{
		Label:            label,
		deletedNeighbors: 0,
		distance:         NewDistance(),
		contracted:       false,
	}
}

// NewVertex Create vertex with vertex number
func NewVertex(vertexNum int) *Vertex {
	return &Vertex{
		vertexNum:        vertexNum,
		deletedNeighbors: 0,
		distance:         NewDistance(),
		contracted:       false,
	}
}

// computeImportance Update importance of vertex
func (vertex *Vertex) computeImportance() {
	vertex.shortcutCover = len(vertex.inEdges) * len(vertex.outEdges)
	vertex.incidentEdges = len(vertex.inEdges) + len(vertex.outEdges)
	vertex.edgeDiff = vertex.shortcutCover - vertex.incidentEdges
	// Spatial diversity heuristic: for each node maintain a count of the number of neighbors that have already been contracted [vertex.deletedNeighbors], and add this to the Importance
	// note: the more neighbours have already been contracted, the later this node will be contracted
	vertex.importance = vertex.edgeDiff + vertex.deletedNeighbors
}

// Distance Information about contraction between source vertex and contraction vertex
type Distance struct {
	contractID  int
	sourceID    int
	distance    float64
	queryDist   float64
	revDistance float64
}

// NewDistance Constructor for Distance
func NewDistance() *Distance {
	return &Distance{
		contractID:  -1,
		sourceID:    -1,
		distance:    math.MaxFloat64,
		revDistance: math.MaxFloat64,
		queryDist:   math.MaxFloat64,
	}
}
