package ch

import (
	"log"
	"testing"
)

func TestExport(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestExport is starting...")
	log.Println(len(g.contracts)) // 268420
	log.Println(len(g.Vertices))  // 588804
	err := g.ExportToFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
	}
}

func TestImportedFileShortestPath(t *testing.T) {
	g, err := ImportFromFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
	}
	log.Println("TestImportedFileShortestPath is starting...")
	u := 69618
	v := 5924
	ans, path := g.ShortestPath(u, v)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
	}
}
