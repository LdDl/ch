package ch

import (
	"fmt"
)

var (
	// ErrGraphIsFrozen Graph is frozen, so it can not be modified.
	ErrGraphIsFrozen = fmt.Errorf("Graph has been frozen")
	// ErrCHNotPrepared Contraction hierarchies have not been prepared yet.
	ErrCHNotPrepared = fmt.Errorf("Contraction hierarchies have not been prepared")
	// ErrVertexNotFound Vertex with given label was not found in graph.
	ErrVertexNotFound = fmt.Errorf("Vertex not found")
	// ErrEdgeNotFound Edge between given vertices was not found.
	ErrEdgeNotFound = fmt.Errorf("Edge not found")
)
