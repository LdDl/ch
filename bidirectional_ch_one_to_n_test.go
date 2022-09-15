package ch

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestOneToManyShortestPath(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestShortestPath is starting...")
	u := int64(106600)
	v := []int64{5924, 81611, 69618, 68427, 68490}
	correctAns, correctPath := []float64{61089.42195558673, 94961.78959757874, 78692.8292369651, 61212.00481622628, 71101.1080090782}, []int{418, 866, 591, 314, 353}
	ans, path := g.ShortestPathOneToMany(u, v)
	for i := range path {
		if len(path[i]) != correctPath[i] {
			t.Errorf("Num of vertices in path should be %d, but got %d", correctPath[i], len(path[i]))
			return
		}
	}
	for i := range ans {
		if Round(ans[i], 0.00005) != Round(correctAns[i], 0.00005) {
			t.Errorf("Length of path should be %f, but got %f", correctAns[i], ans[i])
			return
		}
	}

	t.Log("TestShortestPath is Ok!")
}

func BenchmarkShortestPathOneToMany(b *testing.B) {
	b.Log("BenchmarkShortestPathOneToMany is starting...")
	rand.Seed(1337)
	for k := 2.0; k <= 8; k++ {
		n := int(math.Pow(2, k))
		g, err := generateSyntheticGraph(n)
		if err != nil {
			b.Error(err)
			return
		}
		b.ResetTimer()
		b.Run(fmt.Sprintf("%s/%d/vertices-%d-edges-%d-shortcuts-%d", "CH shortest path", n, len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum()), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := int64(rand.Intn(len(g.Vertices)))
				v := []int64{
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
				}
				ans, path := g.ShortestPathOneToMany(u, v)
				_, _ = ans, path
			}
		})
	}
}

func BenchmarkOldWayShortestPathOneToMany(b *testing.B) {
	b.Log("BenchmarkOldWayShortestPathOneToMany is starting...")
	rand.Seed(1337)
	for k := 2.0; k <= 8; k++ {
		n := int(math.Pow(2, k))
		g, err := generateSyntheticGraph(n)
		if err != nil {
			b.Error(err)
			return
		}
		b.ResetTimer()
		b.Run(fmt.Sprintf("%s/%d/vertices-%d-edges-%d-shortcuts-%d", "CH shortest path", n, len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum()), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := int64(rand.Intn(len(g.Vertices)))
				v := []int64{
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
				}
				for vv := range v {
					ans, path := g.ShortestPath(u, v[vv])
					_, _ = ans, path
				}
			}
		})
	}
}

func BenchmarkTargetNodesShortestPathOneToMany(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
	}
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	b.Log("BenchmarkTargetNodesShortestPathOneToMany is starting...")
	b.ResetTimer()

	rand.Seed(1337)
	for k := 1.0; k <= 7; k++ {
		u := int64(106600)
		n := int(math.Pow(2, k))
		targets := 1 + rand.Intn(n)
		v := make([]int64, targets)
		for t := 0; t < targets; t++ {
			v[t] = int64(rand.Intn(len(g.Vertices)))
		}
		b.Run(fmt.Sprintf("%s/%d/vertices-%d-edges-%d-shortcuts-%d-targets-%d", "CH shortest path", n, len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum(), targets), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ans, path := g.ShortestPathOneToMany(u, v)
				_, _ = ans, path
			}
		})
	}
}

func BenchmarkTargetNodesOldWayShortestPathOneToMany(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
	}
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	b.Log("BenchmarkTargetNodesShortestPathOneToMany is starting...")
	b.ResetTimer()

	rand.Seed(1337)
	for k := 1.0; k <= 7; k++ {
		u := int64(106600)
		n := int(math.Pow(2, k))
		targets := 1 + rand.Intn(n)
		v := make([]int64, targets)
		for t := 0; t < targets; t++ {
			v[t] = int64(rand.Intn(len(g.Vertices)))
		}
		b.Run(fmt.Sprintf("%s/%d/vertices-%d-edges-%d-shortcuts-%d-targets-%d", "CH shortest path", n, len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum(), targets), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for vv := range v {
					ans, path := g.ShortestPath(u, v[vv])
					_, _ = ans, path
				}
			}
		})
	}
}

func BenchmarkStaticCaseShortestPathOneToMany(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
	}
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	b.Log("BenchmarkStaticCaseShortestPathOneToMany is starting...")
	b.ResetTimer()

	b.Run(fmt.Sprintf("%s/vertices-%d-edges-%d-shortcuts-%d", "CH shortest path (one to many)", len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum()), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			u := int64(106600)
			v := []int64{5924, 81611, 69618, 68427, 68490}
			ans, path := g.ShortestPathOneToMany(u, v)
			_, _ = ans, path
		}
	})
}

func BenchmarkStaticCaseOldWayShortestPathOneToMany(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
	}
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	b.Log("BenchmarkStaticCaseOldWayShortestPathOneToMany is starting...")
	b.ResetTimer()

	b.Run(fmt.Sprintf("%s/vertices-%d-edges-%d-shortcuts-%d", "CH shortest path (one to many)", len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum()), func(b *testing.B) {
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
