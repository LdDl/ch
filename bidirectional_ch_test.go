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

func TestExport(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "benchmark_graph.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestExport is starting...")
	log.Println(len(g.contracts)) // 268420
	log.Println(len(g.Vertices))  // 588804
	err := g.ExportToFile("export_bench.csv")
	if err != nil {
		t.Error(err)
	}
}

func TestImportedFileShortestPath(t *testing.T) {
	g, err := ImportFromFile("export_bench.csv")
	if err != nil {
		t.Error(err)
	}
	log.Println("TestImportedFileShortestPath is starting...")
	u := 144031
	v := 452090
	ans, path := g.ShortestPath(u, v)
	if len(path) != 1966 {
		t.Errorf("Num of vertices in path should be 1966, but got %d", len(path))
	}
	if ans != 329520.4412391192 {
		t.Errorf("Length of path should be 329520.4412391192, but got %f", ans)
	}
}

func TestShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "benchmark_graph.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestShortestPath is starting...")
	u := 144031
	v := 452090
	ans, path := g.ShortestPath(u, v)
	if len(path) != 1966 {
		t.Errorf("Num of vertices in path should be 1966, but got %d", len(path))
	}
	if ans != 329520.4412391192 {
		t.Errorf("Length of path should be 329520.4412391192, but got %f", ans)
	}
}

func BenchmarkShortestPath(b *testing.B) {
	g := Graph{}
	graphFromCSV(&g, "benchmark_graph.csv")
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
	graphFromCSV(&g, "benchmark_graph.csv")
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
