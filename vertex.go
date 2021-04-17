package ch

import (
	"math"
)

// Vertex All information about vertex
type Vertex struct {
	vertexNum int64
	Label     int64

	inIncidentEdges []incidentEdge
	outEdges        []int64
	outECost        []float64

	orderPos      int
	contracted    bool
	distance      *Distance
	edgeDiff      int
	delNeighbors  int
	shortcutCover int
	importance    int
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
	vertex.edgeDiff = len(vertex.inIncidentEdges)*len(vertex.outEdges) - len(vertex.inIncidentEdges) - len(vertex.outEdges)
	vertex.shortcutCover = len(vertex.inIncidentEdges) + len(vertex.outEdges)
	vertex.importance = vertex.edgeDiff*14 + vertex.shortcutCover*25 + vertex.delNeighbors*10
}

// incidentEdge incident edge to correspondence
type incidentEdge struct {
	vertexID int64
	cost     float64
}

// Distance Information about contraction between source vertex and contraction vertex
type Distance struct {
	contractID  int64
	sourceID    int64
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
