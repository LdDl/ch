package ch

import (
	"fmt"
	"testing"
)

func TestLoadOsmGraph(t *testing.T) {
	cfg := OsmConfiguration{
		TagName: "highway",
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
	g, err := ImportFromOSMFile("data/moscow_center_reduced.osm.pbf", &cfg)
	if err != nil {
		t.Error(err)
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	t.Log("TestLoadOsmGraph is starting...")

	u := int64(272650046)
	v := int64(7012442362)

	correctPath := 96
	correctAns := 2952.003039
	ans, path := g.ShortestPath(u, v)
	if len(path) != correctPath {
		t.Errorf("Num of vertices in path should be %d, but got %d", correctPath, len(path))
	}
	if Round(ans, 0.00005) != Round(correctAns, 0.00005) {
		t.Errorf("Length of path should be %f, but got %f", correctAns, ans)
	}
	t.Log("TestLoadOsmGraph is Ok!")
	t.Error(0)
}

func TestMiddlePoint(t *testing.T) {
	p1 := Coord{
		Lon: 37.6417350769043,
		Lat: 55.751849391735284,
	}
	p2 := Coord{
		Lon: 37.668514251708984,
		Lat: 55.73261980350401,
	}
	res := Coord{
		Lon: 37.65512796336629,
		Lat: 55.742235325526806,
	}
	fmt.Println(MiddlePoint(p1, p2) == res)
	t.Error(0)
}
