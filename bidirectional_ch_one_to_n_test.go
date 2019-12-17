package ch

import (
	"fmt"
	"math"
	"testing"
)

func TestOneToManyShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	t.Log("TestShortestPath is starting...")
	u := int64(106600)
	v := []int64{5924, 81611, 69618, 68427, 68490}
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

func BenchmarkShortestPathOneToMany(b *testing.B) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	b.Log("BenchmarkShortestPathOneToMany is starting...")
	b.ResetTimer()

	for k := 0.; k <= 12; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%s/%d/vertices-%d", "CH shortest path (one to many)", n, len(g.Vertices)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := int64(106600)
				v := []int64{5924, 81611, 69618, 68427, 68490}
				ans, path := g.ShortestPathOneToMany(u, v)
				_, _ = ans, path
			}
		})
	}
}

func BenchmarkOldWayShortestPathOneToMany(b *testing.B) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	b.Log("BenchmarkOldWayShortestPathOneToMany is starting...")
	b.ResetTimer()

	for k := 0.; k <= 12; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%s/%d/vertices-%d", "CH shortest path (one to many)", n, len(g.Vertices)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := int64(106600)
				v := []int64{5924, 81611, 69618, 68427, 68490}
				for vv := range v {
					ans, path := g.ShortestPath(u, v[vv])
					_, _ = ans, path
				}
			}
		})
	}
}
