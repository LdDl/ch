package ch

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestShortestPath is starting...")
	u := 144031
	v := 452090
	ans, path := g.ShortestPath(u, v)
	if len(path) != 199 {
		t.Errorf("Num of vertices in path should be 199, but got %d", len(path))
	}
	if ans != 19815.1930182987 {
		t.Errorf("Length of path should be 19815.1930182987, but got %f", ans)
	}
}

func BenchmarkShortestPath(b *testing.B) {
	g := Graph{}
	graphFromCSV(&g, "pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("BenchmarkShortestPath is starting...")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := 144031
		v := 452090
		ans, path := g.ShortestPath(u, v)
		_, _ = ans, path
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
		source, err := strconv.Atoi(record[0])
		if err != nil {
			log.Panicln(err)
		}
		target, err := strconv.Atoi(record[1])
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

		graph.AddEdge(source, target, weight)
		if oneway == "B" {
			graph.AddEdge(target, source, weight)
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
