package ch

import (
	"math"
	"testing"
)

func TestImportedFileShortestPath(t *testing.T) {
	g, err := ImportFromFile("data/export_pgrouting.csv", "data/export_pgrouting_vertices.csv", "data/export_pgrouting_shortcuts.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("TestImportedFileShortestPath is starting...")
	u := int64(69618)
	v := int64(5924)

	correctNumShortcuts := int64(394840)
	correctNumVertices := 187853
	evaluatedShortcuts := g.GetShortcutsNum()
	if evaluatedShortcuts != correctNumShortcuts {
		t.Errorf("Number of contractions should be %d, but got %d", correctNumShortcuts, evaluatedShortcuts)
	}
	if len(g.Vertices) != correctNumVertices {
		t.Errorf("Number of vertices should be %d, but got %d", correctNumVertices, len(g.Vertices))
	}

	ans, path := g.ShortestPath(u, v)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
		return
	}
	correctCost := 19135.6581215226
	if math.Abs(ans-correctCost) > eps {
		t.Errorf("Cost of path should be %f, but got %f", correctCost, ans)
		return
	}
	t.Log("TestImportedFileShortestPath is Ok!")
}

func TestImportedFileRecustomization(t *testing.T) {
	g, err := ImportFromFile("data/export_pgrouting.csv", "data/export_pgrouting_vertices.csv", "data/export_pgrouting_shortcuts.csv")
	if err != nil {
		t.Error(err)
		return
	}

	// Verify CH is marked as prepared after import
	if !g.chPrepared {
		t.Error("chPrepared should be true after import")
		return
	}

	// Verify contractionOrder was built
	if len(g.contractionOrder) != len(g.Vertices) {
		t.Errorf("contractionOrder length (%d) should match vertices count (%d)", len(g.contractionOrder), len(g.Vertices))
		return
	}

	// Verify shortcutsByVia was populated
	totalShortcutsInIndex := 0
	for _, shortcuts := range g.shortcutsByVia {
		totalShortcutsInIndex += len(shortcuts)
	}
	if int64(totalShortcutsInIndex) != g.GetShortcutsNum() {
		t.Errorf("shortcutsByVia should contain all shortcuts: got %d, expected %d", totalShortcutsInIndex, g.GetShortcutsNum())
		return
	}

	source := int64(69618)
	target := int64(5924)

	// Get initial path
	initialCost, initialPath := g.ShortestPath(source, target)
	if len(initialPath) < 2 {
		t.Error("Initial path too short for testing")
		return
	}

	// Find an edge to update
	from := initialPath[0]
	to := initialPath[1]
	internalFrom := g.mapping[from]
	internalTo := g.mapping[to]
	var originalWeight float64
	for _, e := range g.Vertices[internalFrom].outIncidentEdges {
		if e.vertexID == internalTo {
			originalWeight = e.weight
			break
		}
	}

	if originalWeight <= 0 {
		t.Error("Could not find edge weight")
		return
	}

	// Update edge weight and recustomize
	err = g.UpdateEdgeWeight(from, to, originalWeight*2, true)
	if err != nil {
		t.Errorf("UpdateEdgeWeight failed: %v", err)
		return
	}

	// Verify cost changed
	newCost, _ := g.ShortestPath(source, target)
	if newCost <= initialCost {
		t.Errorf("Cost should increase after doubling edge weight: initial=%f, new=%f", initialCost, newCost)
		return
	}

	// Restore original weight
	err = g.UpdateEdgeWeight(from, to, originalWeight, true)
	if err != nil {
		t.Errorf("UpdateEdgeWeight (restore) failed: %v", err)
		return
	}

	// Verify cost restored
	restoredCost, _ := g.ShortestPath(source, target)
	if math.Abs(restoredCost-initialCost) > eps {
		t.Errorf("Cost should be restored: initial=%f, restored=%f", initialCost, restoredCost)
		return
	}

	t.Logf("Recustomization on imported graph works: initial=%f, after update=%f, restored=%f", initialCost, newCost, restoredCost)
}
