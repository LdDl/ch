package ch

import (
	"testing"
)

func TestImportedFileShortestPath(t *testing.T) {
	g, err := ImportFromFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
	}
	t.Log("TestImportedFileShortestPath is starting...")
	u := int64(69618)
	v := int64(5924)

	ans, path := g.ShortestPath(u, v)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
	}
	t.Log("TestImportedFileShortestPath is Ok!")
}
