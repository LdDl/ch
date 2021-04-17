package ch

import (
	"log"
	"testing"
)

func TestExport(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	t.Log("TestExport is starting...")
	correctNumShortcuts := 91757
	correctNumVertices := 187853
	if len(g.shortcuts) != correctNumShortcuts {
		t.Errorf("Number of contractions should be %d, but got %d", correctNumShortcuts, len(g.shortcuts))
	}
	if len(g.Vertices) != correctNumVertices {
		t.Errorf("Number of vertices should be %d, but got %d", correctNumVertices, len(g.Vertices))
	}
	err := g.ExportToFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("TestExport is Ok!")
}
