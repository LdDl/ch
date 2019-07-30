package ch

import (
	"fmt"
	"log"
	"testing"
)

type V struct {
	from   int
	to     int
	weight float64
}

func TestRestriction(t *testing.T) {
	vertices := []V{
		V{from: 1, to: 2, weight: 1.0},
		V{from: 2, to: 3, weight: 3.0},
		V{from: 3, to: 4, weight: 1.0},
		V{from: 4, to: 5, weight: 1.0},
		V{from: 5, to: 6, weight: 1.0},
		V{from: 5, to: 7, weight: 1.0},
		V{from: 2, to: 5, weight: 1.0},
		V{from: 8, to: 2, weight: 1.0},
	}

	graph := Graph{}

	for i := range vertices {
		graph.CreateVertex(vertices[i].from)
		graph.CreateVertex(vertices[i].to)
		graph.AddEdge(vertices[i].from, vertices[i].to, vertices[i].weight)
	}

	// restrictions := make(map[int]map[int]int)
	// restrictions[0] = make(map[int]int)
	// restrictions[0][4] = 1
	// restrictions[1] = make(map[int]int)
	// restrictions[1][6] = 4

	restrictions := make(map[int]map[int]int)
	restrictions[0] = make(map[int]int)
	restrictions[0][1] = 4
	restrictions[1] = make(map[int]int)
	restrictions[1][4] = 6

	// fmt.Println(restrictions)
	graph.restrictions = restrictions

	// hard coded
	// if _, ok := graph.contracts[0]; !ok {
	// 	graph.contracts[0] = make(map[int]int)
	// 	graph.contracts[0][4] = 1
	// }
	// graph.contracts[0][4] = 1

	// graph.Vertices[0].outEdges = append(graph.Vertices[0].outEdges, 4)
	// graph.Vertices[0].outECost = append(graph.Vertices[0].outECost, math.MaxFloat64)
	// graph.Vertices[4].inEdges = append(graph.Vertices[4].inEdges, 0)
	// graph.Vertices[4].inECost = append(graph.Vertices[4].inECost, math.MaxFloat64)

	// graph.computeImportance()
	// var extractNum int
	// for graph.pqImportance.Len() != 0 {
	// 	vertex := heap.Pop(graph.pqImportance).(*Vertex)
	// 	vertex.computeImportance()
	// 	if graph.pqImportance.Len() != 0 && vertex.importance > graph.pqImportance.Peek().(*Vertex).importance {
	// 		graph.pqImportance.Push(vertex)
	// 		continue
	// 	}
	// 	graph.Vertices[vertex.vertexNum].orderPos = extractNum
	// 	extractNum = extractNum + 1
	// }
	// graph.PrepareContracts()

	cost, path := graph.VanillaShortestPath(1, 5)
	fmt.Println(cost, path)

	cost, path = graph.VanillaShortestPath(2, 7)
	fmt.Println(cost, path)

	// cost1, path1 := graph.ShortestPath(1, 5)
	// fmt.Println(cost1, path1)

	t.Error("done")
}

func TestExport(t *testing.T) {
	g := Graph{}
	graphFromCSV(&g, "data/pgrouting_osm.csv")
	log.Println("Please wait until contraction hierarchy is prepared")
	g.PrepareContracts()
	log.Println("TestExport is starting...")
	log.Println(len(g.contracts)) // 268420
	log.Println(len(g.Vertices))  // 588804
	err := g.ExportToFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
	}
}

func TestImportedFileShortestPath(t *testing.T) {
	g, err := ImportFromFile("data/export_pgrouting.csv")
	if err != nil {
		t.Error(err)
	}
	log.Println("TestImportedFileShortestPath is starting...")
	u := 69618
	v := 5924
	ans, path := g.ShortestPath(u, v)
	if len(path) != 160 {
		t.Errorf("Num of vertices in path should be 160, but got %d", len(path))
	}
	if Round(ans, 0.00005) != Round(19135.6581215226, 0.00005) {
		t.Errorf("Length of path should be 19135.6581215226, but got %f", ans)
	}
}
