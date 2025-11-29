package ch

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestShortestPath(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestShortestPath is starting...")
	u := int64(69618)
	v := int64(5924)
	//
	ans, path := g.ShortestPath(u, v)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
		return
	}
	correctCost := 19135.6581215226
	if math.Abs(ans-correctCost) > eps {
		t.Errorf("Cost of path should be %f, but got %f", correctCost, ans)
		return
	}
	t.Log("TestShortestPath is Ok!")
}

func TestBothVanillaAndCH(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestAndSHVanPath is starting...")

	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		rndV := g.Vertices[rand.Intn(len(g.Vertices))].Label
		rndU := g.Vertices[rand.Intn(len(g.Vertices))].Label
		ansCH, pathCH := g.ShortestPath(rndV, rndU)
		ansVanilla, pathVanilla := g.VanillaShortestPath(rndV, rndU)
		if len(pathCH) != len(pathVanilla) {
			t.Errorf("Num of vertices in path should be %d, but got %d", len(pathVanilla), len(pathCH))
			return
		}
		if math.Abs(ansCH-ansVanilla) > eps {
			t.Errorf("Cost of path should be %f, but got %f", ansVanilla, ansCH)
			return
		}
	}
	t.Log("TestAndSHVanPath is Ok!")
}

func BenchmarkShortestPath(b *testing.B) {
	b.Log("BenchmarkShortestPath is starting...")
	rand.Seed(1337)
	for k := 2.; k <= 8; k++ {
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
				v := int64(rand.Intn(len(g.Vertices)))
				ans, path := g.ShortestPath(u, v)
				_, _ = ans, path
			}
		})
	}
}

func BenchmarkStaticCaseShortestPath(b *testing.B) {
	b.Log("BenchmarkStaticCaseShortestPath is starting...")
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	g.PrepareContractionHierarchies()
	b.Run(fmt.Sprintf("%s/vertices-%d-edges-%d-shortcuts-%d", "CH shortest path", len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum()), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			u := int64(69618)
			v := int64(5924)
			ans, path := g.ShortestPath(u, v)
			_, _ = ans, path
		}
	})
}

// BenchmarkOneCaseShortestPath/CH_shortest_path/vertices-187853-edges-366113-shortcuts-394840-12         	     891	   1412347 ns/op	 3460158 B/op	    1027 allocs/op

func BenchmarkPrepareContracts(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	g.PrepareContractionHierarchies()
}

func TestBadSpatialShortestPath(t *testing.T) {
	rand.Seed(1337)
	g := Graph{}
	numVertices := 5
	lastVertex := int64(numVertices + 1)
	for i := 0; i < numVertices; i++ {
		idx := int64(i)
		g.CreateVertex(idx + 1)
		g.CreateVertex(idx + 2)
		g.AddEdge(idx+1, idx+2, rand.Float64())
	}
	g.AddEdge(lastVertex, 1, rand.Float64())
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestBadSpatialShortestPath is starting...")
	u := int64(1)
	v := int64(5)

	ans, path := g.ShortestPath(u, v)
	if len(path) != 5 {
		t.Errorf("Num of vertices in path should be 5, but got %d", len(path))
		return
	}
	correctCost := 2.348720
	if math.Abs(ans-correctCost) > eps {
		t.Errorf("Cost of path should be %f, but got %f", correctCost, ans)
		return
	}
	t.Log("TestBadSpatialShortestPath is Ok!")
}

func TestLittleShortestPath(t *testing.T) {
	g := Graph{}
	g.CreateVertex(0)
	g.CreateVertex(1)
	g.CreateVertex(2)
	g.CreateVertex(3)
	g.CreateVertex(4)
	g.CreateVertex(5)
	g.CreateVertex(6)
	g.CreateVertex(7)
	g.CreateVertex(8)
	g.CreateVertex(9)
	g.AddEdge(0, 1, 6.0)
	g.AddEdge(1, 0, 5.0)
	g.AddEdge(1, 9, 3.0)
	g.AddEdge(1, 2, 4.0)
	g.AddEdge(2, 3, 2.0)
	g.AddEdge(3, 2, 2.0)
	g.AddEdge(3, 4, 2.0)
	g.AddEdge(4, 3, 1.0)
	g.AddEdge(0, 4, 0.5)
	g.AddEdge(0, 4, 3.0)
	g.AddEdge(9, 8, 2.0)
	g.AddEdge(4, 8, 13.0)
	g.AddEdge(8, 5, 6.5)
	g.AddEdge(5, 4, 3.5)
	g.AddEdge(7, 8, 1.0)
	g.AddEdge(6, 7, 1.0)
	g.AddEdge(5, 6, 2.0)
	g.AddEdge(5, 6, 4.0)

	g.PrepareContractionHierarchies()
	t.Log("TestLittleShortestPath is starting...")
	u := int64(0)
	v := int64(7)
	//
	ans, path := g.ShortestPath(u, v)
	if len(path) != 7 {
		t.Errorf("Num of vertices in path should be 7, but got %d", len(path))
	}

	correctCost := 20.5
	if math.Abs(ans-correctCost) > eps {
		t.Errorf("Cost of path should be %f, but got %f", correctCost, ans)
		return
	}
	t.Log("TestLittleShortestPath is Ok!")
}

func TestVertexAlternatives(t *testing.T) {
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
	t.Log("TestVertexAlternatives is starting...")
	sources := []VertexAlternative{
		{Label: 0, AdditionalDistance: 1.0},
		{Label: 3, AdditionalDistance: 2.0},
	}
	targets := []VertexAlternative{
		{Label: 2, AdditionalDistance: 2.0},
		{Label: 5, AdditionalDistance: 1.0},
	}
	ans, path := g.ShortestPathWithAlternatives(sources, targets)
	t.Log("ShortestPathWithAlternatives returned", ans, path)
	if len(path) != len(expectedPath) {
		t.Errorf("Num of vertices in path should be %d, but got %d", len(expectedPath), len(path))
	}
	for i := range expectedPath {
		if path[i] != expectedPath[i] {
			t.Errorf("Path item %d should be %d, but got %d", i, expectedPath[i], path[i])
		}
	}
	correctCost := 4.0
	if math.Abs(ans-correctCost) > eps {
		t.Errorf("Cost of path should be %f, but got %f", correctCost, ans)
		return
	}
	t.Log("TestVertexAlternatives is Ok!")
}

func TestVertexAlternativesConnected(t *testing.T) {
	//  S-(1)-0-----\
	//  |     |     |
	// (1)   (1)   (3)
	//  |     |     |
	//  \-----1-(1)-T

	g := Graph{}
	g.CreateVertex(0)
	g.CreateVertex(1)
	g.AddEdge(0, 1, 1.0)

	expectedPath := []int64{1}

	g.PrepareContractionHierarchies()
	t.Log("TestVertexAlternativesConnected is starting...")
	sources := []VertexAlternative{
		{Label: 0, AdditionalDistance: 1.0},
		{Label: 1, AdditionalDistance: 1.0},
	}
	targets := []VertexAlternative{
		{Label: 0, AdditionalDistance: 3.0},
		{Label: 1, AdditionalDistance: 1.0},
	}
	ans, path := g.ShortestPathWithAlternatives(sources, targets)
	t.Log("ShortestPathWithAlternatives returned", ans, path)
	if len(path) != len(expectedPath) {
		t.Errorf("Num of vertices in path should be %d, but got %d", len(expectedPath), len(path))
	}
	for i := range expectedPath {
		if path[i] != expectedPath[i] {
			t.Errorf("Path item %d should be %d, but got %d", i, expectedPath[i], path[i])
		}
	}
	correctCost := 2.0
	if math.Abs(ans-correctCost) > eps {
		t.Errorf("Cost of path should be %f, but got %f", correctCost, ans)
		return
	}
	t.Log("TestVertexAlternativesConnected is Ok!")
}

func graphFromCSV(graph *Graph, fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))

	reader.Comma = ';'
	// reader.LazyQuotes = true

	// Read header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		source, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return err
		}
		target, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return err
		}

		oneway := record[2]
		weight, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return err
		}

		err = graph.CreateVertex(source)
		if err != nil {
			return err
		}
		err = graph.CreateVertex(target)
		if err != nil {
			return err
		}

		err = graph.AddEdge(source, target, weight)
		if err != nil {
			return err
		}
		if oneway == "B" {
			err = graph.AddEdge(target, source, weight)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func generateSyntheticGraph(verticesNum int) (*Graph, error) {
	rand.Seed(1337)
	graph := Graph{}
	var i int
	for i = 1; i < verticesNum; i++ {
		source := int64(i - 1)
		err := graph.CreateVertex(source)
		if err != nil {
			return nil, err
		}
		for j := 1; j < verticesNum; j++ {
			if j == i {
				continue
			}
			target := int64(j)
			err := graph.CreateVertex(target)
			if err != nil {
				return nil, err
			}
			// weight := rand.Float64()
			weight := 0.01 + rand.Float64()*(10-0.01) // Make more dispersion
			err = graph.AddEdge(source, target, weight)
			if err != nil {
				return nil, err
			}
			addReverse := rand.Intn(2)
			if addReverse != 0 {
				// Add reverse edge imitating bidirectional=true
				err = graph.AddEdge(target, source, weight)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	graph.PrepareContractionHierarchies()
	return &graph, nil
}

// TestAllShortestPathMethods compares all shortest path methods to ensure they produce identical results.
// It picks 10 random sources and 10 random targets, then compares:
// - ShortestPath (NxM calls)
// - ShortestPathOneToMany (N calls)
// - ShortestPathManyToMany (1 call)
// - VanillaShortestPath (NxM calls)
func TestAllShortestPathMethods(t *testing.T) {
	rand.Seed(1337)
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestAllShortestPathMethods is starting...")

	const numSources = 10
	const numTargets = 10

	sources := make([]int64, numSources)
	targets := make([]int64, numTargets)
	for i := 0; i < numSources; i++ {
		sources[i] = g.Vertices[rand.Intn(len(g.Vertices))].Label
	}
	for i := 0; i < numTargets; i++ {
		targets[i] = g.Vertices[rand.Intn(len(g.Vertices))].Label
	}

	t.Logf("Sources (%d): %v", numSources, sources)
	t.Logf("Targets (%d): %v", numTargets, targets)

	// 1. ShortestPath (NxM calls)
	t.Log("Running ShortestPath (NxM calls)...")
	costsSP := make([][]float64, numSources)
	pathLensSP := make([][]int, numSources)
	startSP := time.Now()
	for i, src := range sources {
		costsSP[i] = make([]float64, numTargets)
		pathLensSP[i] = make([]int, numTargets)
		for j, tgt := range targets {
			cost, path := g.ShortestPath(src, tgt)
			costsSP[i][j] = cost
			pathLensSP[i][j] = len(path)
		}
	}
	durationSP := time.Since(startSP)
	t.Logf("ShortestPath (NxM = %d calls): %v", numSources*numTargets, durationSP)

	// 2. ShortestPathOneToMany (N calls)
	t.Log("Running ShortestPathOneToMany (N calls)...")
	costsO2M := make([][]float64, numSources)
	pathLensO2M := make([][]int, numSources)
	startO2M := time.Now()
	for i, src := range sources {
		costs, paths := g.ShortestPathOneToMany(src, targets)
		costsO2M[i] = costs
		pathLensO2M[i] = make([]int, numTargets)
		for j := range paths {
			pathLensO2M[i][j] = len(paths[j])
		}
	}
	durationO2M := time.Since(startO2M)
	t.Logf("ShortestPathOneToMany (N = %d calls): %v", numSources, durationO2M)

	// 3. ShortestPathManyToMany (1 call)
	t.Log("Running ShortestPathManyToMany (1 call)...")
	startM2M := time.Now()
	costsM2M, pathsM2M := g.ShortestPathManyToMany(sources, targets)
	durationM2M := time.Since(startM2M)
	pathLensM2M := make([][]int, numSources)
	for i := range pathsM2M {
		pathLensM2M[i] = make([]int, numTargets)
		for j := range pathsM2M[i] {
			pathLensM2M[i][j] = len(pathsM2M[i][j])
		}
	}
	t.Logf("ShortestPathManyToMany (1 call): %v", durationM2M)

	// 4. VanillaShortestPath (NxM calls)
	t.Log("Running VanillaShortestPath (NxM calls)...")
	costsVanilla := make([][]float64, numSources)
	pathLensVanilla := make([][]int, numSources)
	startVanilla := time.Now()
	for i, src := range sources {
		costsVanilla[i] = make([]float64, numTargets)
		pathLensVanilla[i] = make([]int, numTargets)
		for j, tgt := range targets {
			cost, path := g.VanillaShortestPath(src, tgt)
			costsVanilla[i][j] = cost
			pathLensVanilla[i][j] = len(path)
		}
	}
	durationVanilla := time.Since(startVanilla)
	t.Logf("VanillaShortestPath (NxM = %d calls): %v", numSources*numTargets, durationVanilla)

	// Compare results
	t.Log("Comparing results...")
	allMatch := true

	for i := 0; i < numSources; i++ {
		for j := 0; j < numTargets; j++ {
			src, tgt := sources[i], targets[j]

			// Get all costs
			costSP := costsSP[i][j]
			costO2M := costsO2M[i][j]
			costM2M := costsM2M[i][j]
			costVanilla := costsVanilla[i][j]

			// Compare costs (use Vanilla as ground truth)
			if math.Abs(costSP-costVanilla) > eps {
				t.Errorf("[%d][%d] src=%d tgt=%d: ShortestPath cost=%f != Vanilla cost=%f",
					i, j, src, tgt, costSP, costVanilla)
				allMatch = false
			}
			if math.Abs(costO2M-costVanilla) > eps {
				t.Errorf("[%d][%d] src=%d tgt=%d: OneToMany cost=%f != Vanilla cost=%f",
					i, j, src, tgt, costO2M, costVanilla)
				allMatch = false
			}
			if math.Abs(costM2M-costVanilla) > eps {
				t.Errorf("[%d][%d] src=%d tgt=%d: ManyToMany cost=%f != Vanilla cost=%f",
					i, j, src, tgt, costM2M, costVanilla)
				allMatch = false
			}

			// Compare path lengths
			lenSP := pathLensSP[i][j]
			lenO2M := pathLensO2M[i][j]
			lenM2M := pathLensM2M[i][j]
			lenVanilla := pathLensVanilla[i][j]

			if lenSP != lenVanilla {
				t.Errorf("[%d][%d] src=%d tgt=%d: ShortestPath pathLen=%d != Vanilla pathLen=%d",
					i, j, src, tgt, lenSP, lenVanilla)
				allMatch = false
			}
			if lenO2M != lenVanilla {
				t.Errorf("[%d][%d] src=%d tgt=%d: OneToMany pathLen=%d != Vanilla pathLen=%d",
					i, j, src, tgt, lenO2M, lenVanilla)
				allMatch = false
			}
			if lenM2M != lenVanilla {
				t.Errorf("[%d][%d] src=%d tgt=%d: ManyToMany pathLen=%d != Vanilla pathLen=%d",
					i, j, src, tgt, lenM2M, lenVanilla)
				allMatch = false
			}
		}
	}

	if allMatch {
		t.Log("All methods produce identical results!")
	}

	// Summary
	t.Log("=== TIMING SUMMARY ===")
	t.Logf("VanillaShortestPath (NxM = %d): %v (baseline)", numSources*numTargets, durationVanilla)
	t.Logf("ShortestPath        (NxM = %d): %v (%.1fx faster than Vanilla)",
		numSources*numTargets, durationSP, float64(durationVanilla)/float64(durationSP))
	t.Logf("ShortestPathOneToMany  (N = %d): %v (%.1fx faster than Vanilla)",
		numSources, durationO2M, float64(durationVanilla)/float64(durationO2M))
	t.Logf("ShortestPathManyToMany (1 call): %v (%.1fx faster than Vanilla)",
		durationM2M, float64(durationVanilla)/float64(durationM2M))

	t.Log("TestAllShortestPathMethods is Ok!")
}
