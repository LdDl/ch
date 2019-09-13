package ch

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestShortestPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestShortestPath is starting...")
	u := 69618
	v := 5924

	ans, path := g.ShortestPath(u, v)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
	}
	log.Println("TestShortestPath is Ok!")
}

func TestAndSHVanPath(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestAndSHVanPath is starting...")

	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000; i++ {
		rndV := g.Vertices[rand.Intn(len(g.Vertices))].Label
		rndU := g.Vertices[rand.Intn(len(g.Vertices))].Label
		st := time.Now()
		ansCH, pathCH := g.ShortestPath(rndV, rndU)
		fmt.Println("ch:", time.Since(st))
		st = time.Now()
		ansVanilla, pathVanilla := g.VanillaShortestPath(rndV, rndU)
		fmt.Println("vanilla:", time.Since(st))
		if len(pathCH) != len(pathVanilla) {
			log.Println(rndV, rndU, pathVanilla)
			t.Errorf("Num of vertices in path should be %d, but got %d", len(pathVanilla), len(pathCH))
		}
		if Round(ansCH, 0.00005) != Round(ansVanilla, 0.00005) {
			t.Errorf("Length of path should be %f, but got %f", ansVanilla, ansCH)
		}
	}
	log.Println("TestAndSHVanPath is Ok!")
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

func TestShortestPath_v2(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestShortestPath is starting...")
	// u := 21831
	// v := 77701
	u := 69618
	v := 5924

	fmt.Println("source:", u, ", target:", v)

	fmt.Println("time", "ans", "path")
	ts := time.Now()
	ans, path := g.VanillaShortestPath(u, v) // Верный результат, [21831 51225 65164 23323 82109 29766 77701]
	dts := ts.Sub(time.Now())
	fmt.Println("vanila", dts, ans, path)
	ts = time.Now()
	ans, path = g.ShortestPath(u, v) // И близко неверный, 21831 51225 65164 23323 82109 29766 77701 (wtf???) 4651 34542 79740 89676 78394 14683 7899 34962 84862 40509 37933 35372 40605 34460 26637 29766 77701
	dts = ts.Sub(time.Now())
	fmt.Println("bi ch ", dts, ans, path)
	ts = time.Now()
	ans, path = g.ShortestPathProc(u, v) // И близко неверный, 21831 51225 65164 23323 82109 29766 77701 (wtf???) 4651 34542 79740 89676 78394 14683 7899 34962 84862 40509 37933 35372 40605 34460 26637 29766 77701
	dts = ts.Sub(time.Now())
	fmt.Println("bi go ", dts, ans, path)
	if len(path) != 7 {
		t.Errorf("Num of vertices in path should be 7, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(600.074025, 0.00005) {
		t.Errorf("Length of path should be 600.074025, but got %f", ans)
	}
	log.Println("TestShortestPath is Ok!")
	t.Errorf("s")
}
