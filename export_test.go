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
	// t.Log(len(g.contracts)) // 268420
	// t.Log(len(g.Vertices))  // 588804
	err := g.ExportToFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("TestExport is Ok!")
}
