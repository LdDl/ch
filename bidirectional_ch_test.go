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
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
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
		if Round(ansCH, 0.00005) != Round(ansVanilla, 0.00005) {
			t.Errorf("Length of path should be %f, but got %f", ansVanilla, ansCH)
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
		b.Run(fmt.Sprintf("%s/%d/vertices-%d-shortcuts-%d", "CH shortest path", n, len(g.Vertices), g.shortcutsNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := int64(rand.Intn(len(g.Vertices)))
				v := int64(rand.Intn(len(g.Vertices)))
				ans, path := g.ShortestPath(u, v)
				_, _ = ans, path
			}
		})
	}
}

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
	if Round(ans, 0.00005) != Round(2.348720, 0.00005) {
		t.Errorf("Length of path should be 2.348720, but got %f", ans)
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
	if Round(ans, 0.00005) != Round(20.5, 0.00005) {
		t.Errorf("Length of path should be 20.0, but got %f", ans)
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
	if Round(ans, 0.00005) != Round(4.0, 0.00005) {
		t.Errorf("Length of path should be 4.0, but got %f", ans)
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
	if Round(ans, 0.00005) != Round(2.0, 0.00005) {
		t.Errorf("Length of path should be 2.0, but got %f", ans)
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
			weight := rand.Float64()
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
