package ch

import (
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestQueryPoolShortestPath(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestQueryPoolShortestPath is starting...")

	pool := g.NewQueryPool()

	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		rndV := g.Vertices[rand.Intn(len(g.Vertices))].Label
		rndU := g.Vertices[rand.Intn(len(g.Vertices))].Label

		ansPool, pathPool := pool.ShortestPath(rndV, rndU)
		ansGraph, pathGraph := g.ShortestPath(rndV, rndU)

		if len(pathPool) != len(pathGraph) {
			t.Errorf("Num of vertices in path should be %d, but got %d", len(pathGraph), len(pathPool))
			return
		}
		if math.Abs(ansPool-ansGraph) > eps {
			t.Errorf("Cost of path should be %f, but got %f", ansGraph, ansPool)
			return
		}
	}
	t.Log("TestQueryPoolShortestPath is Ok!")
}

func TestQueryPoolConcurrentQueries(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestQueryPoolConcurrentQueries is starting...")

	pool := g.NewQueryPool()

	source := g.Vertices[0].Label
	target := g.Vertices[len(g.Vertices)-1].Label
	expectedCost, _ := g.ShortestPath(source, target)

	numGoroutines := 100
	numQueriesPerGoroutine := 100

	var wg sync.WaitGroup
	errors := make(chan string, numGoroutines*numQueriesPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numQueriesPerGoroutine; j++ {
				ans, _ := pool.ShortestPath(source, target)
				if math.Abs(ans-expectedCost) > eps {
					errors <- "wrong cost"
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	var errCount int
	for range errors {
		errCount++
	}
	if errCount > 0 {
		t.Errorf("Concurrent queries produced %d errors", errCount)
		return
	}
	t.Log("TestQueryPoolConcurrentQueries is Ok!")
}

func TestQueryPoolShortestPathOneToMany(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestQueryPoolShortestPathOneToMany is starting...")

	pool := g.NewQueryPool()

	source := g.Vertices[0].Label
	targets := make([]int64, 10)
	for i := 0; i < 10; i++ {
		targets[i] = g.Vertices[i*100].Label
	}

	costsPool, pathsPool := pool.ShortestPathOneToMany(source, targets)
	costsGraph, pathsGraph := g.ShortestPathOneToMany(source, targets)

	for i := range targets {
		if math.Abs(costsPool[i]-costsGraph[i]) > eps {
			t.Errorf("Cost mismatch for target %d: pool=%v, graph=%v", targets[i], costsPool[i], costsGraph[i])
			return
		}
		if len(pathsPool[i]) != len(pathsGraph[i]) {
			t.Errorf("Path length mismatch for target %d: pool=%v, graph=%v", targets[i], len(pathsPool[i]), len(pathsGraph[i]))
			return
		}
	}
	t.Log("TestQueryPoolShortestPathOneToMany is Ok!")
}

func TestQueryPoolConcurrentOneToMany(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestQueryPoolConcurrentOneToMany is starting...")

	pool := g.NewQueryPool()

	source := g.Vertices[0].Label
	targets := make([]int64, 10)
	for i := 0; i < 10; i++ {
		targets[i] = g.Vertices[i*100].Label
	}

	expectedCosts, _ := g.ShortestPathOneToMany(source, targets)

	numGoroutines := 50
	numQueriesPerGoroutine := 50

	var wg sync.WaitGroup
	errors := make(chan string, numGoroutines*numQueriesPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numQueriesPerGoroutine; j++ {
				costs, _ := pool.ShortestPathOneToMany(source, targets)
				for k, cost := range costs {
					if math.Abs(cost-expectedCosts[k]) > eps {
						errors <- "cost mismatch"
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	var errCount int
	for range errors {
		errCount++
	}
	if errCount > 0 {
		t.Errorf("Concurrent OneToMany queries produced %d errors", errCount)
		return
	}
	t.Log("TestQueryPoolConcurrentOneToMany is Ok!")
}

func BenchmarkQueryPoolShortestPath(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
		return
	}
	g.PrepareContractionHierarchies()
	pool := g.NewQueryPool()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.ShortestPath(10, 7281)
	}
}

func BenchmarkQueryPoolConcurrent(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
		return
	}
	g.PrepareContractionHierarchies()
	pool := g.NewQueryPool()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.ShortestPath(10, 7281)
		}
	})
}

func TestQueryPoolShortestPathManyToMany(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestQueryPoolShortestPathManyToMany is starting...")

	pool := g.NewQueryPool()

	sources := []int64{106600, 69618}
	targets := []int64{5924, 81611, 69618}

	// Get results from pool
	costsPool, pathsPool := pool.ShortestPathManyToMany(sources, targets)
	// Get results from graph directly
	costsGraph, pathsGraph := g.ShortestPathManyToMany(sources, targets)

	// Compare results
	for i := range sources {
		for j := range targets {
			if math.Abs(costsPool[i][j]-costsGraph[i][j]) > eps {
				t.Errorf("Cost mismatch for [%d][%d]: pool=%v, graph=%v", i, j, costsPool[i][j], costsGraph[i][j])
			}
			if len(pathsPool[i][j]) != len(pathsGraph[i][j]) {
				t.Errorf("Path length mismatch for [%d][%d]: pool=%v, graph=%v", i, j, len(pathsPool[i][j]), len(pathsGraph[i][j]))
			}
		}
	}
	t.Log("TestQueryPoolShortestPathManyToMany is Ok!")
}

func TestQueryPoolConcurrentManyToMany(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestQueryPoolConcurrentManyToMany is starting...")

	pool := g.NewQueryPool()

	// Run many concurrent ManyToMany queries
	numGoroutines := 10
	numQueries := 50

	var wg sync.WaitGroup
	errors := make(chan string, numGoroutines*numQueries)

	sources := []int64{106600, 69618}
	targets := []int64{5924, 81611, 69618}

	// Pre-compute expected results
	expectedCosts, expectedPaths := g.ShortestPathManyToMany(sources, targets)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numQueries; j++ {
				costs, paths := pool.ShortestPathManyToMany(sources, targets)
				for si := range sources {
					for ti := range targets {
						if math.Abs(costs[si][ti]-expectedCosts[si][ti]) > eps {
							errors <- "cost mismatch"
						}
						if len(paths[si][ti]) != len(expectedPaths[si][ti]) {
							errors <- "path mismatch"
						}
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	var errCount int
	for range errors {
		errCount++
	}
	if errCount > 0 {
		t.Errorf("Concurrent ManyToMany queries produced %d errors", errCount)
		return
	}
	t.Log("TestQueryPoolConcurrentManyToMany is Ok!")
}
