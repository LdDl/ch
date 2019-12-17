package ch

import (
	"log"
	"testing"
)

func TestLoadOsmGraph(t *testing.T) {
	cfg := Config{
		Name: "highway",
		Tags: []string{
			"motorway",
			"primary",
			"primary_link",
			"road",
			"secondary",
			"secondary_link",
			"residential",
			"tertiary",
			"tertiary_link",
			"unclassified",
			"trunk",
			"trunk_link",
		},
	}
	g, err := LoadOsmGraph("data/moscow_tinao.pbf", cfg)
	if err != nil {
		t.Error(err)
	}
	t.Error("Please wait until contraction hierarchy is prepared")
	// g.PrepareContracts()
	t.Log("TestExport is starting...")
	t.Log(len(g.contracts)) // 268420
	t.Log(len(g.Vertices))  // 588804

	log.Println("TestExport is Ok!")
}
