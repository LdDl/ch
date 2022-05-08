package ch

import (
	"testing"
)

func TestVanillaShortestPath(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
	}
	t.Log("TestShortestPath is starting...")
	u := int64(69618)
	v := int64(5924)
	ans, path := g.VanillaShortestPath(u, v)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
		return
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
		return
	}
	t.Log("TestVanillaShortestPath is Ok!")
}

func Round(x, unit float64) float64 {
	if x > 0 {
		return float64(int64(x/unit+0.5)) * unit
	}
	return float64(int64(x/unit-0.5)) * unit
}
