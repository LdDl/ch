package ch

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	t.Log("TestShortestPath is starting...")
	u := int64(69618)
	v := int64(5924)

	ans, path := g.ShortestPath(u, v)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
	}
	t.Log("TestShortestPath is Ok!")
}

func TestBothVanillaAndCH(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	t.Log("TestAndSHVanPath is starting...")

	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		rndV := g.Vertices[rand.Intn(len(g.Vertices))].Label
		rndU := g.Vertices[rand.Intn(len(g.Vertices))].Label
		ansCH, pathCH := g.ShortestPath(rndV, rndU)
		ansVanilla, pathVanilla := g.VanillaShortestPath(rndV, rndU)
		if len(pathCH) != len(pathVanilla) {
			t.Errorf("Num of vertices in path should be %d, but got %d", len(pathVanilla), len(pathCH))
		}
		if Round(ansCH, 0.00005) != Round(ansVanilla, 0.00005) {
			t.Errorf("Length of path should be %f, but got %f", ansVanilla, ansCH)
		}
	}
	t.Log("TestAndSHVanPath is Ok!")
}

func BenchmarkShortestPath(b *testing.B) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	b.Log("BenchmarkShortestPath is starting...")
	b.ResetTimer()

	for k := 0.; k <= 12; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%s/%d/vertices-%d", "CH shortest path", n, len(g.Vertices)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := int64(69618)
				v := int64(5924)
				ans, path := g.ShortestPath(u, v)
				_, _ = ans, path
			}
		})
	}
}

func BenchmarkPrepareContracts(b *testing.B) {
	g := Graph{}
	graphFromCSV(&g, "pgrouting_osm.csv")
	b.ResetTimer()
	g.PrepareContracts()
}

func graphFromCSV(graph *Graph, fname string) {
	file, err := os.Open(fname)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))

	reader.Comma = ';'
	// reader.LazyQuotes = true

	// Read header
	_, err = reader.Read()
	if err != nil {
		log.Panicln(err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		source, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Panicln(err)
		}
		target, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Panicln(err)
		}

		// if intinslice(source, []int{5606, 5607, 255077, 238618}) == false && intinslice(target, []int{5606, 5607, 255077, 238618}) == false {
		// 	continue
		// }

		oneway := record[2]
		weight, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Panicln(err)
		}

		graph.CreateVertex(source)
		graph.CreateVertex(target)

		err = graph.AddEdge(source, target, weight)
		if err != nil {
			log.Panicln(err)
			continue
		}
		if oneway == "B" {
			graph.AddEdge(target, source, weight)
			if err != nil {
				log.Panicln(err)
				continue
			}
		}
	}
}

func intinslice(elem int, sl []int) bool {
	for i := range sl {
		if sl[i] == elem {
			return true
		}
	}
	return false
}
