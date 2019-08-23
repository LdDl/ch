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
	// log.Println(len(g.contracts)) // 268420
	// log.Println(len(g.Vertices))  // 588804
	err := g.ExportToFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
	}
	log.Println("TestExport is Ok!")
}
