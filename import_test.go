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
