package ch

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestShortestPath is starting...")
	// u := 69618
	// v := 5924
	u := 126826
	v := 5485

	ans, path := g.ShortestPath(u, v)
	fmt.Println(path)
	// 126826 127332 126739 126597 65282 125469 125734 127709 125604 125727 124763 18483 30052 18750 10843 28676 12050 26004 115934 112349 112996 3509 5485
	// 126826 127332 126739 126597 65282 125469 125734 127709 125604 125727 124763 18483 30052 18750 10843 28676 12050 26004 115934 112349 112996 3509 892 19008 17938 10029 26965 27744 20729 1237 18594 7000 2582 26857 11607 28943 11767 27697 12097 15992 12097 27697 11767 28943 11607 26857 2582 7000 18594 1237 20729 27744 26965 10029 17938 19008 892 3509 5485
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
	}
	log.Println("TestShortestPath is Ok!")
}

func BenchmarkShortestPath(b *testing.B) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("BenchmarkShortestPath is starting...")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := 69618
		v := 5924
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
