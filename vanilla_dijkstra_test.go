package ch

import (
	"math"
	"testing"
)

const (
	eps = 0.0001
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
	correctCost := 19135.6581215226
	if math.Abs(ans-correctCost) > eps {
		t.Errorf("Cost of path should be %f, but got %f", correctCost, ans)
		return
	}
	t.Log("TestVanillaShortestPath is Ok!")
}
