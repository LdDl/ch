package ch

import (
	"log"
	"testing"
)

type V struct {
	from   int64
	to     int64
	weight float64
}

func TestVanillaTurnRestrictedShortestPath(t *testing.T) {

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
		err := graph.AddEdge(vertices[i].from, vertices[i].to, vertices[i].weight)
		if err != nil {
			log.Panicln(err)
			continue
		}
	}

	restrictions := make(map[int64]map[int64]int64)
	restrictions[1] = make(map[int64]int64)
	restrictions[1][2] = 5
	restrictions[2] = make(map[int64]int64)
	restrictions[2][5] = 7

	for source, turn := range restrictions {
		for via, target := range turn {
			err := graph.AddTurnRestriction(source, via, target)
			if err != nil {
				log.Panicln(err)
				continue
			}
		}
	}

	ans, path := graph.VanillaTurnRestrictedShortestPath(1, 5)
	rightPath := []int64{1, 2, 3, 4, 5}
	if len(path) != 5 {
		t.Errorf("Run 1: num of vertices in path should be 5, but got %d", len(path))
	}
	for i := range path {
		if path[i] != rightPath[i] {
			t.Errorf("Run 1: vertex in path should be %d, but got %d", path[i], rightPath[i])
		}
	}
	if ans != 6 {
		t.Errorf("Run 1: length of path should be 6, but got %f", ans)
	}

	ans, path = graph.VanillaTurnRestrictedShortestPath(2, 7)
	rightPath = []int64{2, 3, 4, 5, 7}
	if len(path) != 5 {
		t.Errorf("Run 2: num of vertices in path should be 5, but got %d", len(path))
	}
	for i := range path {
		if path[i] != rightPath[i] {
			t.Errorf("Run 2: vertex in path should be %d, but got %d", path[i], rightPath[i])
		}
	}
	if ans != 6 {
		t.Errorf("Run 2: length of path should be 6, but got %f", ans)
	}

	ans, path = graph.VanillaTurnRestrictedShortestPath(1, 7)
	rightPath = []int64{1, 2, 3, 4, 5, 7}
	if len(path) != 6 {
		t.Errorf("Run 3: num of vertices in path should be 6, but got %d", len(path))
	}
	for i := range path {
		if path[i] != rightPath[i] {
			t.Errorf("Run 3: vertex in path should be %d, but got %d", path[i], rightPath[i])
		}
	}
	if ans != 7 {
		t.Errorf("Run 3: length of path should be 7, but got %f", ans)
	}

	t.Log("TestVanillaTurnRestrictedShortestPath is Ok!")
}
