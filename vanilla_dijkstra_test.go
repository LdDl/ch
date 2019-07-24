package ch

import (
	"log"
	"testing"
	"time"
)

func TestVanillaShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	// log.Println("Please wait until contraction hierarchy is prepared")
	// g.PrepareContracts()
	log.Println("TestShortestPath is starting...")
	u := 69618
	v := 5924
	st := time.Now()

	ans, path := g.VanillaShortestPath(u, v)
	// for i := 0; i < 300; i++ {
	// 	ans, path = g.VanillaShortestPath(u, v)
	// }
	log.Println(time.Since(st))
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
	}
}

func Round(x, unit float64) float64 {
	if x > 0 {
		return float64(int64(x/unit+0.5)) * unit
	}
	return float64(int64(x/unit-0.5)) * unit
}
