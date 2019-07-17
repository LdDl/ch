package ch

import (
	"log"
	"testing"
)

func TestExport(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestExport is starting...")
	log.Println(len(g.contracts)) // 268420
	log.Println(len(g.Vertices))  // 588804
	err := g.ExportToFile("export_pgrouting.csv")
	if err != nil {
		t.Error(err)
	}
}

func TestImportedFileShortestPath(t *testing.T) {
	g, err := ImportFromFile("export_pgrouting.csv")
	if err != nil {
		t.Error(err)
	}
	log.Println("TestImportedFileShortestPath is starting...")
	u := 144031
	v := 452090
	ans, path := g.ShortestPath(u, v)
	if len(path) != 199 {
		t.Errorf("Num of vertices in path should be 1966, but got %d", len(path))
	}
	if ans != 19815.1930182987 {
		t.Errorf("Length of path should be 19815.1930182987, but got %f", ans)
	}
}
