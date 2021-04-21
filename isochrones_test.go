package ch

import (
	"testing"
)

func TestIsochrones(t *testing.T) {
	correctIsochrones := map[int64]float64{
		5: 0.0,
		3: 1.0,
		4: 1.0,
		6: 1.0,
		7: 3.0,
		1: 3.0,
		8: 4.0, // <---- Because of breadth-first search
		9: 2.0,
	}
	graph := Graph{}

	vertices := []V{
		V{from: 5, to: 3, weight: 1.0},
		V{from: 5, to: 4, weight: 1.0},
		V{from: 5, to: 6, weight: 1.0},
		V{from: 5, to: 7, weight: 2.0},
		V{from: 3, to: 7, weight: 2.0},
		V{from: 6, to: 9, weight: 1.0},
		V{from: 7, to: 8, weight: 4.0},
		V{from: 7, to: 3, weight: 2.0},
		V{from: 9, to: 8, weight: 2.0},
		V{from: 8, to: 10, weight: 3.0},
		V{from: 3, to: 1, weight: 2.0},
		V{from: 1, to: 2, weight: 3.0},
		V{from: 4, to: 11, weight: 7.0},
		V{from: 11, to: 2, weight: 2.0},
		V{from: 2, to: 11, weight: 2.0},
	}

	for i := range vertices {
		err := graph.CreateVertex(vertices[i].from)
		if err != nil {
			t.Error(err)
			return
		}
		err = graph.CreateVertex(vertices[i].to)
		if err != nil {
			t.Error(err)
			return
		}
		err = graph.AddEdge(vertices[i].from, vertices[i].to, vertices[i].weight)
		if err != nil {
			t.Error(err)
			return
		}
	}

	graph.PrepareContractionHierarchies() // This is excess in current example, but just for proof that contraction map isn't used.

	sourceVertex := int64(5)
	maxCost := 5.0
	isochrones, err := graph.Isochrones(sourceVertex, maxCost)
	if err != nil {
		t.Error(err)
		return
	}

	if len(isochrones) != len(correctIsochrones) {
		t.Errorf("Number of isochrones should be %d, but got %d", len(correctIsochrones), len(isochrones))
		return
	}

	for k, val := range isochrones {
		correctValue, ok := correctIsochrones[k]
		if !ok {
			t.Errorf("Isochrones should contain vertex %d, but it does not", k)
			return
		}
		if val != correctValue {
			t.Errorf("Travel cost to vertex %d should be %f, but got %f", k, correctValue, val)
			return
		}
	}
}
