package ch

import (
	"container/heap"
	"fmt"
	"log"
)

// Graph Graph object
//
// mapping Internal map for 1:1 relation of internal IDs to user's IDs
// Vertices Slice of vertices of graph
// shortcuts Found and stored shortcuts based on contraction hierarchies
type Graph struct {
	shortcuts    map[int64]map[int64]*ShortcutPath
	restrictions map[int64]map[int64]int64
	mapping      map[int64]int64

	Vertices     []Vertex
	edgesNum     int64
	shortcutsNum int64

	frozen  bool
	verbose bool

	Reporter Reporter
}

type Reporter interface {
	VertexSettled(direction, endpointIndex int, vertexID int64, heapSize int)
	EdgeRelaxed(direction, endpointIndex int, vertexID, toVertexID int64, ch bool, heapSize int)
	FoundBetterPath(direction, sourceEndpointIndex, targetEndpointIndex int, vertexID int64, estimate float64)
}

// NewGraph returns pointer to created Graph and does preallocations for processing purposes
func NewGraph() *Graph {
	return &Graph{
		mapping:      make(map[int64]int64),
		Vertices:     make([]Vertex, 0),
		edgesNum:     0,
		shortcutsNum: 0,
		shortcuts:    make(map[int64]map[int64]*ShortcutPath),
		restrictions: make(map[int64]map[int64]int64),
		frozen:       false,
		verbose:      false,
	}
}

// CreateVertex Creates new vertex and assign internal ID to it
//
// label User's definied ID of vertex
func (graph *Graph) CreateVertex(label int64) error {
	if graph.frozen {
		return ErrGraphIsFrozen
	}
	v := Vertex{
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
func (graph *Graph) AddEdge(from, to int64, weight float64) error {
	if graph.frozen {
		return ErrGraphIsFrozen
	}
	graph.edgesNum++
	from = graph.mapping[from]
	to = graph.mapping[to]

	graph.addEdge(from, to, weight)
	return nil
}

func (graph *Graph) addEdge(from, to int64, weight float64) {
	graph.Vertices[from].outIncidentEdges = append(graph.Vertices[from].outIncidentEdges, incidentEdge{vertexID: to, weight: weight})
	graph.Vertices[to].inIncidentEdges = append(graph.Vertices[to].inIncidentEdges, incidentEdge{vertexID: from, weight: weight})
}

// AddShortcut Adds new shortcut between two vertices
//
// from - User's definied ID of first vertex of shortcut
// to - User's definied ID of last vertex of shortcut
// via - User's defined ID of vertex through which the shortcut exists
// weight - User's definied weight of shortcut
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
	graph.shortcutsNum++
	return nil
}

// PrepareContractionHierarchies Compute contraction hierarchies
func (graph *Graph) PrepareContractionHierarchies() {
	pqImportance := graph.computeImportance()
	graph.Preprocess(pqImportance)
	graph.Freeze()
}

// SetVerbose sets verbose parameter for debugging purposes
func (graph *Graph) SetVerbose(flag bool) {
	graph.verbose = flag
}

// computeImportance Returns heap to store computed importance of each vertex
func (graph *Graph) computeImportance() *importanceHeap {
	pqImportance := &importanceHeap{}
	heap.Init(pqImportance)
	for i := range graph.Vertices {
		graph.Vertices[i].computeImportance()
		heap.Push(pqImportance, &graph.Vertices[i])
	}
	graph.Freeze()
	return pqImportance
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

// GetVerticesNum Returns number of vertices in graph
func (graph *Graph) GetVerticesNum() int64 {
	return int64(len(graph.Vertices))
}

// GetShortcutsNum Returns number of shortcuts in graph
func (graph *Graph) GetShortcutsNum() int64 {
	return int64(graph.shortcutsNum)
}

// GetEdgesNum Returns number of edges in graph
func (graph *Graph) GetEdgesNum() int64 {
	return graph.edgesNum
}

// AddTurnRestriction Adds new turn restriction between two vertices via some other vertex
//
// from User's definied ID of source vertex
// via User's definied ID of prohibited vertex (between source and target)
// to User's definied ID of target vertex
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
