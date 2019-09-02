package ch

import (
	"fmt"
	"log"
	"testing"
)

func TestVanillaShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	log.Println("TestShortestPath is starting...")
	// u := 69618
	// v := 5924
	u := 126826
	v := 5485
	ans, path := g.VanillaShortestPath(u, v)
	// 126826 127332 126739 126597 65282 125469 125734 127709 125604 125727 124763 18483 30052 18750 10843 28676 12050 26004 115934 112349 112996 3509 5485
	fmt.Println(path)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
	}
	log.Println("TestVanillaShortestPath is Ok!")
}

func Round(x, unit float64) float64 {
	if x > 0 {
		return float64(int64(x/unit+0.5)) * unit
	}
	return float64(int64(x/unit-0.5)) * unit
}
