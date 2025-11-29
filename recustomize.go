package ch

// UpdateEdgeWeight Updates the weight of an existing edge in the graph.
// This function works with user-defined vertex labels.
//
// from - User's defined ID of source vertex
// to - User's defined ID of target vertex
// weight - New weight for the edge
// needRecustom - If true, Recustomize() is called automatically after update
//
// Returns error if edge is not found or if recustomization fails.
func (graph *Graph) UpdateEdgeWeight(from, to int64, weight float64, needRecustom bool) error {
	fromInternal, ok := graph.mapping[from]
	if !ok {
		return ErrVertexNotFound
	}
	toInternal, ok := graph.mapping[to]
	if !ok {
		return ErrVertexNotFound
	}

	// Update outgoing edge weight
	updatedOut := graph.Vertices[fromInternal].updateOutIncidentEdge(toInternal, weight)
	if !updatedOut {
		return ErrEdgeNotFound
	}

	// Update incoming edge weight
	updatedIn := graph.Vertices[toInternal].updateInIncidentEdge(fromInternal, weight)
	if !updatedIn {
		return ErrEdgeNotFound
	}

	if needRecustom {
		return graph.Recustomize()
	}
	return nil
}

// Recustomize Recomputes all shortcut costs based on current edge weights.
// This is useful after edge weights have been modified (e.g., traffic updates).
// The contraction order is preserved - only the metric (costs) is updated.
//
// Returns error if CH has not been prepared yet.
func (graph *Graph) Recustomize() error {
	if !graph.chPrepared {
		return ErrCHNotPrepared
	}

	// Process vertices in contraction order
	// Shortcuts via vertex V are processed after all shortcuts via vertices
	// with lower orderPos have been updated.
	for _, viaVertex := range graph.contractionOrder {
		shortcuts := graph.shortcutsByVia[viaVertex]
		for _, shortcut := range shortcuts {
			// Recompute shortcut cost: cost(From->Via) + cost(Via->To)
			fromViaCost := graph.getEdgeCost(shortcut.From, shortcut.Via)
			viaToaCost := graph.getEdgeCost(shortcut.Via, shortcut.To)

			if fromViaCost < 0 || viaToaCost < 0 {
				// Edge not found - should not happen in a valid CH
				continue
			}

			newCost := fromViaCost + viaToaCost
			if newCost != shortcut.Cost {
				shortcut.Cost = newCost

				// Update incident edges
				graph.Vertices[shortcut.From].updateOutIncidentEdge(shortcut.To, newCost)
				graph.Vertices[shortcut.To].updateInIncidentEdge(shortcut.From, newCost)
			}
		}
	}

	return nil
}

// getEdgeCost Returns the cost of edge from source to target (internal IDs).
// Returns -1 if edge is not found.
func (graph *Graph) getEdgeCost(from, to int64) float64 {
	for _, edge := range graph.Vertices[from].outIncidentEdges {
		if edge.vertexID == to {
			return edge.weight
		}
	}
	return -1
}
