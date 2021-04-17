package ch

import (
	"container/heap"
	"fmt"
	"log"
)

// Graph Graph object
//
// pqImportance Heap to store importance of each vertex
// pqComparator Heap to store traveled distance
// mapping Internal map for 1:1 relation of internal IDs to user's IDs
// Vertices Slice of vertices of graph
// nodeOrdering Ordering of vertices
// shortcuts Found and stored shortcuts based on contraction hierarchies
//
type Graph struct {
	pqImportance *importanceHeap
	pqComparator *distanceHeap

	mapping      map[int64]int64
	Vertices     []*Vertex
	nodeOrdering []int64

	shortcuts    map[int64]map[int64]*ContractionPath
	restrictions map[int64]map[int64]int64

	frozen bool
}

// CreateVertex Creates new vertex and assign internal ID to it
//
// label User's definied ID of vertex
//
func (graph *Graph) CreateVertex(label int64) error {
	if graph.frozen {
		return ErrGraphIsFrozen
	}
	v := &Vertex{
		Label:        label,
		delNeighbors: 0,
		distance:     NewDistance(),
		contracted:   false,
	}
	if graph.mapping == nil {
		graph.mapping = make(map[int64]int64)
	}
	if graph.shortcuts == nil {
		graph.shortcuts = make(map[int64]map[int64]*ContractionPath)
	}

	if _, ok := graph.mapping[label]; !ok {
		v.vertexNum = int64(len(graph.Vertices))
		graph.mapping[label] = v.vertexNum
		graph.Vertices = append(graph.Vertices, v)
	}
	return nil
}

// AddVertex Adds vertex with provided internal ID
//
// labelExternal User's definied ID of vertex
// labelInternal internal ID of vertex
//
func (graph *Graph) AddVertex(labelExternal, labelInternal int64) error {
	if graph.frozen {
		return ErrGraphIsFrozen
	}
	v := &Vertex{
		Label:        labelExternal,
		delNeighbors: 0,
		distance:     NewDistance(),
		contracted:   true,
		vertexNum:    labelInternal,
	}
	if graph.mapping == nil {
		graph.mapping = make(map[int64]int64)
	}
	if graph.shortcuts == nil {
		graph.shortcuts = make(map[int64]map[int64]*ContractionPath)
	}

	if _, ok := graph.mapping[labelExternal]; !ok {
		graph.mapping[labelExternal] = labelInternal
		if labelInternal < int64(len(graph.Vertices)) {
			graph.Vertices[labelInternal] = v
		} else {
			diff := labelInternal - int64(len(graph.Vertices)) + 1
			empty := make([]*Vertex, diff)
			graph.Vertices = append(graph.Vertices, empty...)
			graph.Vertices[labelInternal] = v
		}
	}
	return nil
}

// AddEdge Adds new edge between two vertices
//
// from User's definied ID of first vertex of edge
// to User's definied ID of last vertex of edge
// weight User's definied weight of edge
//
func (graph *Graph) AddEdge(from, to int64, weight float64) error {
	if graph.frozen {
		return ErrGraphIsFrozen
	}

	from = graph.mapping[from]
	to = graph.mapping[to]

	graph.Vertices[from].outIncidentEdges = append(graph.Vertices[from].outIncidentEdges, incidentEdge{to, weight})
	graph.Vertices[to].inIncidentEdges = append(graph.Vertices[to].inIncidentEdges, incidentEdge{from, weight})
	return nil
}

// AddTurnRestriction Adds new turn restriction between two vertices via some other vertex
//
// from User's definied ID of source vertex
// via User's definied ID of prohibited vertex (between source and target)
// to User's definied ID of target vertex
//
func (graph *Graph) AddTurnRestriction(from, via, to int64) error {
	if graph.frozen {
		return ErrGraphIsFrozen
	}

	from = graph.mapping[from]
	via = graph.mapping[via]
	to = graph.mapping[to]

	if graph.restrictions == nil {
		graph.restrictions = make(map[int64]map[int64]int64)
	}

	if _, ok := graph.restrictions[from]; !ok {
		graph.restrictions[from] = make(map[int64]int64)
		if _, ok := graph.restrictions[from][via]; ok {
			log.Printf("Warning: Please notice, library supports only one 'from-via' relation currently. From %d Via %d\n", from, via)
		}
		graph.restrictions[from][via] = to
	}
	return nil
}

// computeImportance Compute vertices' importance
func (graph *Graph) computeImportance() {
	graph.pqImportance = &importanceHeap{}
	heap.Init(graph.pqImportance)
	for i := 0; i < len(graph.Vertices); i++ {
		graph.Vertices[i].computeImportance()
		heap.Push(graph.pqImportance, graph.Vertices[i])
	}
	graph.Freeze()
}

// PrepareContracts Compute contraction hierarchies
func (graph *Graph) PrepareContracts() {
	graph.computeImportance()
	graph.nodeOrdering = graph.Preprocess()
	graph.Freeze()
}

// Freeze Freeze graph. Should be called after contraction hierarchies had been prepared.
func (graph *Graph) Freeze() {
	graph.frozen = true
}

// Unfreeze Freeze graph. Should be called if graph modification is needed.
func (graph *Graph) Unfreeze() {
	fmt.Println("Warning: You will need to call PrepareContracts() or even refresh graph again if you want to modify graph data")
	graph.frozen = false
}
