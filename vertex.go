package ch

import (
	"fmt"
	"math"
)

type Vertex struct {
	vertexNum int
	Label     int

	inEdges  []int
	inECost  []float64
	outEdges []int
	outECost []float64

	orderPos      int
	contracted    bool
	distance      *Distance
	processed     *Processed
	edgeDiff      int
	delNeighbors  int
	shortcutCover int
	importance    int
}

func (v *Vertex) PrintInOut() {
	fmt.Println(v.outEdges, v.inEdges)
}

func MakeVertex(label int) *Vertex {
	return &Vertex{
		Label:        label,
		delNeighbors: 0,
		distance:     NewDistance(),
		processed:    NewProcessed(),
		contracted:   false,
	}
}

func NewVertex(vertexNum int) *Vertex {
	return &Vertex{
		vertexNum:    vertexNum,
		delNeighbors: 0,
		distance:     NewDistance(),
		processed:    NewProcessed(),
		contracted:   false,
	}
}

// func (vertex *Vertex) GetOut() map[int]float64 {
// 	return vertex.out
// }

// func (vertex *Vertex) GetIn() map[int]float64 {
// 	return vertex.in
// }

func (vertex *Vertex) computeImportance() {
	vertex.edgeDiff = len(vertex.inEdges)*len(vertex.outEdges) - len(vertex.inEdges) - len(vertex.outEdges)
	vertex.shortcutCover = len(vertex.inEdges) + len(vertex.outEdges)
	vertex.importance = vertex.edgeDiff*14 + vertex.shortcutCover*25 + vertex.delNeighbors*10
}

type Distance struct {
	contractID  int
	sourceID    int
	distance    float64
	queryDist   float64
	revDistance float64
}

func NewDistance() *Distance {
	return &Distance{
		contractID:  -1,
		sourceID:    -1,
		distance:    math.MaxFloat64,
		revDistance: math.MaxFloat64,
		queryDist:   math.MaxFloat64,
	}
}

type Processed struct {
	forwProcessed bool
	revProcessed  bool
}

func NewProcessed() *Processed {
	return &Processed{}
}
