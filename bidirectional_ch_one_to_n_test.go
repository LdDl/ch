package ch

import (
	"testing"
)

func TestOneToManyShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	t.Log("TestShortestPath is starting...")
	u := 106600
	v := []int{5924, 81611, 69618, 68427, 68490}
	correctAns, correctPath := []float64{61089.42195558673, 94961.78959757874, 78692.8292369651, 61212.00481622628, 71101.1080090782}, []int{418, 866, 591, 314, 353}
	ans, path := g.ShortestPathOneToMany(u, v)
	for i := range path {
		if len(path[i]) != correctPath[i] {
			t.Errorf("Num of vertices in path should be %d, but got %d", correctPath[i], len(path[i]))
		}
	}
	for i := range ans {
		if Round(ans[i], 0.00005) != Round(correctAns[i], 0.00005) {
			t.Errorf("Length of path should be %f, but got %f", correctAns[i], ans[i])
		}
	}

	t.Log("TestShortestPath is Ok!")
}
