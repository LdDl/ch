package ch

import (
	"container/heap"
	"log"
)

// Graph Graph object
//
// pqImportance Heap to store importance of each vertex
// pqComparator Heap to store traveled distance
// mapping Internal map for 1:1 relation of internal IDs to user's IDs
// Vertices Slice of vertices of graph
// nodeOrdering Ordering of vertices
// contracts found and stored contraction hierarchies
//
type Graph struct {
	pqImportance *importanceHeap
	pqComparator *distanceHeap

	mapping      map[int]int
	Vertices     []*Vertex
	nodeOrdering []int

	contracts    map[int]map[int]int
	restrictions map[int]map[int]int
}

// CreateVertex Creates new vertex and assign internal ID to it
//
// label User's definied ID of vertex
//
func (graph *Graph) CreateVertex(label int) {
	v := &Vertex{
		Label:        label,
		delNeighbors: 0,
		distance:     NewDistance(),
		processed:    NewProcessed(),
		contracted:   false,
	}
	if graph.mapping == nil {
		graph.mapping = make(map[int]int)
	}
	if graph.contracts == nil {
		graph.contracts = make(map[int]map[int]int)
	}

	if _, ok := graph.mapping[label]; !ok {
		v.vertexNum = len(graph.Vertices)
		graph.mapping[label] = v.vertexNum
		graph.Vertices = append(graph.Vertices, v)
	}
}

// AddVertex Adds vertex with provided internal ID
//
// labelExternal User's definied ID of vertex
// labelInternal internal ID of vertex
//
func (graph *Graph) AddVertex(labelExternal, labelInternal int) {
	v := &Vertex{
		Label:        labelExternal,
		delNeighbors: 0,
		distance:     NewDistance(),
		processed:    NewProcessed(),
		contracted:   true,
		vertexNum:    labelInternal,
	}
	if graph.mapping == nil {
		graph.mapping = make(map[int]int)
	}
	if graph.contracts == nil {
		graph.contracts = make(map[int]map[int]int)
	}

	if _, ok := graph.mapping[labelExternal]; !ok {
		graph.mapping[labelExternal] = labelInternal
		if labelInternal < len(graph.Vertices) {
			graph.Vertices[labelInternal] = v
		} else {
			diff := labelInternal - len(graph.Vertices) + 1
			empty := make([]*Vertex, diff)
			graph.Vertices = append(graph.Vertices, empty...)
			graph.Vertices[labelInternal] = v
		}
	}
}

// AddEdge Adds new edge between two vertices
//
// from User's definied ID of first vertex of edge
// to User's definied ID of last vertex of edge
// weight User's definied weight of edge
//
func (graph *Graph) AddEdge(from, to int, weight float64) {

	from = graph.mapping[from]
	to = graph.mapping[to]

	graph.Vertices[from].outEdges = append(graph.Vertices[from].outEdges, to)
	graph.Vertices[from].outECost = append(graph.Vertices[from].outECost, weight)

	graph.Vertices[to].inEdges = append(graph.Vertices[to].inEdges, from)
	graph.Vertices[to].inECost = append(graph.Vertices[to].inECost, weight)
}

// AddTurnRestriction Adds new turn restriction between two vertices via some other vertex
//
// from User's definied ID of source vertex
// via User's definied ID of prohibited vertex (between source and target)
// to User's definied ID of target vertex
//
func (graph *Graph) AddTurnRestriction(from, via, to int) {

	from = graph.mapping[from]
	via = graph.mapping[via]
	to = graph.mapping[to]

	if graph.restrictions == nil {
		graph.restrictions = make(map[int]map[int]int)
	}

	if _, ok := graph.restrictions[from]; !ok {
		graph.restrictions[from] = make(map[int]int)
		if _, ok := graph.restrictions[from][via]; ok {
			log.Printf("Warning: Please notice, library supports only one 'from-via' relation currently. From %d Via %d\n", from, via)
		}
		graph.restrictions[from][via] = to
	}

}

// computeImportance Compute vertices' importance
func (graph *Graph) computeImportance() {
	graph.pqImportance = &importanceHeap{}
	heap.Init(graph.pqImportance)
	for i := 0; i < len(graph.Vertices); i++ {
		graph.Vertices[i].edgeDiff = len(graph.Vertices[i].inEdges)*len(graph.Vertices[i].outEdges) - len(graph.Vertices[i].inEdges) - len(graph.Vertices[i].outEdges)
		graph.Vertices[i].shortcutCover = len(graph.Vertices[i].inEdges) + len(graph.Vertices[i].outEdges)
		graph.Vertices[i].importance = graph.Vertices[i].edgeDiff*14 + graph.Vertices[i].shortcutCover*25 + graph.Vertices[i].delNeighbors*10
		heap.Push(graph.pqImportance, graph.Vertices[i])
	}
}

// PrepareContracts Compute contraction hierarchies
func (graph *Graph) PrepareContracts() {
	graph.computeImportance()
	graph.nodeOrdering = graph.Preprocess()
}
