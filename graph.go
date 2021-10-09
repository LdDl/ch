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

	shortcuts    map[int64]map[int64]*ShortcutPath
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
		graph.shortcuts = make(map[int64]map[int64]*ShortcutPath)
	}

	if _, ok := graph.mapping[label]; !ok {
		v.vertexNum = int64(len(graph.Vertices))
		graph.mapping[label] = v.vertexNum
		graph.Vertices = append(graph.Vertices, v)
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

// AddShortcut Adds new shortcut between two vertices
//
// from - User's definied ID of first vertex of shortcut
// to - User's definied ID of last vertex of shortcut
// via - User's defined ID of vertex through which the shortcut exists
// weight - User's definied weight of shortcut
//
func (graph *Graph) AddShortcut(from, to, via int64, weight float64) error {
	if graph.frozen {
		return ErrGraphIsFrozen
	}
	fromInternal := graph.mapping[from]
	toInternal := graph.mapping[to]
	viaInternal := graph.mapping[via]
	if _, ok := graph.shortcuts[fromInternal]; !ok {
		graph.shortcuts[fromInternal] = make(map[int64]*ShortcutPath)
		graph.shortcuts[fromInternal][toInternal] = &ShortcutPath{
			From: fromInternal,
			To:   toInternal,
			Via:  viaInternal,
			Cost: weight,
		}
	}
	graph.shortcuts[fromInternal][toInternal] = &ShortcutPath{
		From: fromInternal,
		To:   toInternal,
		Via:  viaInternal,
		Cost: weight,
	}
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

// PrepareContractionHierarchies Compute contraction hierarchies
func (graph *Graph) PrepareContractionHierarchies() {
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
	fmt.Println("Warning: You will need to call PrepareContractionHierarchies() or even refresh graph again if you want to modify graph data")
	graph.frozen = false
}

// shortcutsNum Calculate number of shortcuts (useful for debugging purposes)
func (graph *Graph) shortcutsNum() int {
	ans := 0
	for _, i := range graph.shortcuts {
		ans += len(i)
	}
	return ans
}

// IsShortcut Returns (vertex_id; true) if edge is a shortcut (edge defined as two vertices)
//
// If source or taget vertex is not found then returns (-1; false)
// If edge is not a shortcut then returns (-1; false)
//
func (graph *Graph) IsShortcut(labelFromVertex, labelToVertex int64) (int64, bool) {
	source, ok := graph.mapping[labelFromVertex]
	if !ok {
		return -1, ok
	}
	target, ok := graph.mapping[labelToVertex]
	if !ok {
		return -1, ok
	}
	shortcut, ok := graph.shortcuts[source][target]
	if !ok {
		return -1, ok
	}
	return graph.Vertices[shortcut.Via].Label, ok
}
