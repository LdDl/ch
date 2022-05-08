package ch

import (
	"log"
	"testing"
)

func TestExport(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestExport is starting...")
	correctNumShortcuts := int64(394840)
	correctNumVertices := 187853
	evaluatedShortcuts := g.GetShortcutsNum()
	if evaluatedShortcuts != correctNumShortcuts {
		t.Errorf("Number of contractions should be %d, but got %d", correctNumShortcuts, evaluatedShortcuts)
	}
	if len(g.Vertices) != correctNumVertices {
		t.Errorf("Number of vertices should be %d, but got %d", correctNumVertices, len(g.Vertices))
	}
	err = g.ExportToFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("TestExport is Ok!")
}
