package ch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecustomizeNotPrepared(t *testing.T) {
	g := NewGraph()
	g.CreateVertex(0)
	g.CreateVertex(1)
	g.AddEdge(0, 1, 1.0)

	err := g.Recustomize()
	assert.Equal(t, ErrCHNotPrepared, err)
}

func TestUpdateEdgeWeightVertexNotFound(t *testing.T) {
	g := NewGraph()
	g.CreateVertex(0)
	g.CreateVertex(1)
	g.AddEdge(0, 1, 1.0)
	g.PrepareContractionHierarchies()

	err := g.UpdateEdgeWeight(999, 1, 2.0, false)
	assert.Equal(t, ErrVertexNotFound, err)
	err = g.UpdateEdgeWeight(0, 999, 2.0, false)
	assert.Equal(t, ErrVertexNotFound, err)
}

func TestUpdateEdgeWeightEdgeNotFound(t *testing.T) {
	g := NewGraph()
	g.CreateVertex(0)
	g.CreateVertex(1)
	g.CreateVertex(2)
	g.AddEdge(0, 1, 1.0)
	g.PrepareContractionHierarchies()

	// Edge 0->2 doesn't exist
	err := g.UpdateEdgeWeight(0, 2, 2.0, false)
	assert.Equal(t, ErrEdgeNotFound, err)
}

func TestRecustomizeSimpleGraph(t *testing.T) {
	// Create a simple graph: 0 -> 1 -> 2
	// Initial weights: 0->1: 1.0, 1->2: 1.0
	// Shortest path 0->2 should be 2.0
	g := NewGraph()
	g.CreateVertex(0)
	g.CreateVertex(1)
	g.CreateVertex(2)
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.PrepareContractionHierarchies()

	// Verify initial shortest path
	cost, path := g.ShortestPath(0, 2)
	assert.Equal(t, 2.0, cost)
	assert.Equal(t, []int64{0, 1, 2}, path)

	// Update edge weight: 0->1 becomes 5.0
	err := g.UpdateEdgeWeight(0, 1, 5.0, true)
	assert.NoError(t, err)

	// After recustomization, path 0->2 should cost 6.0
	cost, path = g.ShortestPath(0, 2)
	assert.Equal(t, 6.0, cost)
	assert.Equal(t, []int64{0, 1, 2}, path)
}

func TestRecustomizeDiamondGraph(t *testing.T) {
	// Create a diamond graph:
	//     1
	//    / \
	//   0   3
	//    \ /
	//     2
	// Initial: 0->1: 1.0, 0->2: 2.0, 1->3: 1.0, 2->3: 1.0
	// Shortest 0->3: via 1, cost 2.0

	g := NewGraph()
	for i := int64(0); i <= 3; i++ {
		g.CreateVertex(i)
	}
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(0, 2, 2.0)
	g.AddEdge(1, 3, 1.0)
	g.AddEdge(2, 3, 1.0)
	g.PrepareContractionHierarchies()

	// Initial shortest path 0->3
	cost, path := g.ShortestPath(0, 3)
	assert.Equal(t, 2.0, cost)
	assert.Equal(t, []int64{0, 1, 3}, path)

	// Make path via 1 expensive: 0->1 becomes 10.0
	err := g.UpdateEdgeWeight(0, 1, 10.0, true)
	assert.NoError(t, err)

	// Now shortest path should be via 2: cost 3.0
	cost, path = g.ShortestPath(0, 3)
	assert.Equal(t, 3.0, cost)
	assert.Equal(t, []int64{0, 2, 3}, path)
}

func TestRecustomize_PathChangesAfterBlockingEdge(t *testing.T) {
	// Create a longer chain graph to force shortcut creation:
	//
	//   0 --1-- 1 --1-- 2 --1-- 3 --1-- 4    (top path: cost 4)
	//   |                               |
	//   5 ----------------------------- +    (bottom path: cost 5)
	//
	// The chain 0->1->2->3->4 should create shortcuts when middle vertices are contracted
	// Initially: shortest path 0->4 via top path, cost 4
	// After blocking edge 1->2: shortest path becomes 0->4 via bottom, cost 5

	g := NewGraph()
	for i := int64(0); i <= 4; i++ {
		g.CreateVertex(i)
	}

	// Top path: chain 0->1->2->3->4 (cost 4)
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(2, 3, 1.0)
	g.AddEdge(3, 4, 1.0)

	// Bottom path: direct 0->4 (cost 5)
	g.AddEdge(0, 4, 5.0)

	g.PrepareContractionHierarchies()

	// Print shortcuts and contraction order
	t.Logf("Contraction order: %v", g.contractionOrder)
	t.Logf("Number of shortcuts: %d", g.shortcutsNum)
	for from, toMap := range g.shortcuts {
		for to, shortcut := range toMap {
			t.Logf("Shortcut: %d -> %d via %d, cost %f", from, to, shortcut.Via, shortcut.Cost)
		}
	}

	// Verify initial shortest path
	cost, path := g.ShortestPath(0, 4)
	t.Logf("Initial path: %v, cost: %f", path, cost)
	assert.Equal(t, 4.0, cost)
	assert.Equal(t, []int64{0, 1, 2, 3, 4}, path)

	// "Block" edge 1->2 by making it very expensive
	err := g.UpdateEdgeWeight(1, 2, 9999999.0, true)
	assert.NoError(t, err)

	// Print shortcuts after recustomization
	t.Log("After recustomization:")
	for from, toMap := range g.shortcuts {
		for to, shortcut := range toMap {
			t.Logf("Shortcut: %d -> %d via %d, cost %f", from, to, shortcut.Via, shortcut.Cost)
		}
	}

	// Now shortest path should use direct edge 0->4, cost 5
	cost, path = g.ShortestPath(0, 4)
	t.Logf("Path after blocking 1->2: %v, cost: %f", path, cost)
	assert.Equal(t, 5.0, cost)
	assert.Equal(t, []int64{0, 4}, path)

	// Restore edge 1->2
	err = g.UpdateEdgeWeight(1, 2, 1.0, true)
	assert.NoError(t, err)

	// Path should return to chain, cost 4
	cost, path = g.ShortestPath(0, 4)
	t.Logf("Path after restoring 1->2: %v, cost: %f", path, cost)
	assert.Equal(t, 4.0, cost)
	assert.Equal(t, []int64{0, 1, 2, 3, 4}, path)
}

func TestRecustomizeBatchUpdates(t *testing.T) {
	// Test multiple updates with needRecustom=false, then manual Recustomize()
	g := NewGraph()
	for i := int64(0); i <= 3; i++ {
		g.CreateVertex(i)
	}
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(2, 3, 1.0)
	g.PrepareContractionHierarchies()

	// Initial cost 0->3 should be 3.0
	cost, path := g.ShortestPath(0, 3)
	assert.Equal(t, 3.0, cost)
	assert.Equal(t, []int64{0, 1, 2, 3}, path)

	// Batch update multiple edges without recustomizing
	g.UpdateEdgeWeight(0, 1, 2.0, false)
	g.UpdateEdgeWeight(1, 2, 2.0, false)
	g.UpdateEdgeWeight(2, 3, 2.0, false)

	// Manual recustomize
	err := g.Recustomize()
	assert.NoError(t, err)

	// Now cost should be 6.0
	cost, path = g.ShortestPath(0, 3)
	assert.Equal(t, 6.0, cost)
	assert.Equal(t, []int64{0, 1, 2, 3}, path)
}

func TestRecustomizeLargeGraph(t *testing.T) {
	// Use the test graph from CSV file
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	g.PrepareContractionHierarchies()

	source := int64(69618)
	target := int64(5924)

	// Get initial path
	initialCost, initialPath := g.ShortestPath(source, target)
	assert.InDelta(t, 19135.6581215226, initialCost, 10e-6)
	assert.Equal(t, 160, len(initialPath))
	t.Logf("Initial cost: %f, path length: %d", initialCost, len(initialPath))

	// Update a random edge and recustomize
	// Find an edge to update - just use some vertex from the path
	if len(initialPath) > 2 {
		from := initialPath[0]
		to := initialPath[1]
		// Double the weight
		internalFrom := g.mapping[from]
		internalTo := g.mapping[to]
		var originalWeight float64
		for _, e := range g.Vertices[internalFrom].outIncidentEdges {
			if e.vertexID == internalTo {
				originalWeight = e.weight
				break
			}
		}
		t.Logf("Updating edge %d -> %d, original weight: %f", from, to, originalWeight)
		if originalWeight > 0 {
			// Increase weight slightly
			err := g.UpdateEdgeWeight(from, to, originalWeight*2, true)
			assert.NoError(t, err)
			// Path should still exist (maybe different cost)
			newCost, newPath := g.ShortestPath(source, target)
			assert.NotEmpty(t, newPath)
			// New cost should be >= initial cost (we made an edge more expensive)
			assert.GreaterOrEqual(t, newCost, initialCost)
			assert.InDelta(t, newCost, initialCost+(originalWeight), 10e-6)
			assert.Equal(t, len(initialPath), len(newPath))
			t.Logf("New cost after update: %f, path length: %d", newCost, len(newPath))

			// Increase weight back to original
			err = g.UpdateEdgeWeight(from, to, originalWeight, true)
			assert.NoError(t, err)
			// Path should revert to original cost
			finalCost, finalPath := g.ShortestPath(source, target)
			assert.InDelta(t, initialCost, finalCost, 10e-6)
			assert.Equal(t, len(initialPath), len(finalPath))
			t.Logf("Final cost after restoring: %f, path length: %d", finalCost, len(finalPath))

			// Increase weight significantly
			err = g.UpdateEdgeWeight(from, to, originalWeight*200, true)
			assert.NoError(t, err)
			// Path should still exist (maybe different cost)
			highCost, highPath := g.ShortestPath(source, target)
			assert.NotEmpty(t, highPath)
			// New cost should be >= initial cost (we made an edge more expensive)
			assert.GreaterOrEqual(t, highCost, initialCost)
			assert.InDelta(t, highCost, 22474.770773, 10e-5)
			assert.Equal(t, 178, len(highPath))
			t.Logf("High cost after large update: %f, path length: %d", highCost, len(highPath))
		}
	}
}

func TestRecustomizePreservesContractedOrder(t *testing.T) {
	// Verify that orderPos values are preserved after recustomization
	g := NewGraph()
	for i := int64(0); i <= 5; i++ {
		g.CreateVertex(i)
	}
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(2, 3, 1.0)
	g.AddEdge(3, 4, 1.0)
	g.AddEdge(4, 5, 1.0)
	g.PrepareContractionHierarchies()

	// Capture orderPos values
	orderPosBefore := make([]int64, len(g.Vertices))
	for i, v := range g.Vertices {
		orderPosBefore[i] = v.orderPos
	}

	// Update and recustomize
	g.UpdateEdgeWeight(0, 1, 10.0, true)

	// Verify orderPos unchanged
	for i, v := range g.Vertices {
		assert.Equalf(t, orderPosBefore[i], v.orderPos, "Vertex %d orderPos changed after recustomization", i)
	}
}
