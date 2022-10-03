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

func TestOneToManyAlternatives(t *testing.T) {
	//  S-(1)-0-(1)-1-(1)-2
	//  |     |     |     |
	// (2)   (1)   (2)   (2)
	//  |     |     |     |
	//  3-(1)-4-(1)-5-(1)-T

	g := Graph{}
	g.CreateVertex(0)
	g.CreateVertex(1)
	g.CreateVertex(2)
	g.CreateVertex(3)
	g.CreateVertex(4)
	g.CreateVertex(5)
	g.CreateVertex(6)
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(0, 4, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(1, 5, 2.0)
	g.AddEdge(3, 4, 1.0)
	g.AddEdge(4, 5, 1.0)

	expectedPath := []int64{0, 4, 5}

	g.PrepareContractionHierarchies()
	t.Log("TestOneToManyAlternatives is starting...")
	source := []VertexAlternative{
		{Label: 0, AdditionalDistance: 1.0},
		{Label: 3, AdditionalDistance: 2.0},
	}
	targets := []VertexAlternative{
		{Label: 2, AdditionalDistance: 2.0},
		{Label: 5, AdditionalDistance: 1.0},
	}
	ans, paths := g.ShortestPathOneToManyWithAlternatives(source, [][]VertexAlternative{targets})
	t.Log("TestOneToManyAlternatives returned", ans, paths)
	path := paths[0]
	if len(path) != len(expectedPath) {
		t.Errorf("Num of vertices in path should be %d, but got %d", len(expectedPath), len(path))
	}
	for i := range expectedPath {
		if path[i] != expectedPath[i] {
			t.Errorf("Path item %d should be %d, but got %d", i, expectedPath[i], path[i])
		}
	}
	if Round(ans[0], 0.00005) != Round(4.0, 0.00005) {
		t.Errorf("Length of path should be 4.0, but got %f", ans)
	}

	t.Log("TestOneToManyAlternatives is Ok!")
}
