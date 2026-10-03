package ch

import (
	"container/heap"
	"math/rand"
	"reflect"
	"testing"
)

func TestQueryHeapMatchesReference(t *testing.T) {
	rng := rand.New(rand.NewSource(20261001))
	var optimized queryHeap
	var reference vertexDistHeap
	for trial := 0; trial < 500; trial++ {
		for i := 0; i < 100; i++ {
			if len(optimized) == 0 || rng.Intn(3) != 0 {
				value := vertexDist{id: int64(i), dist: float64(rng.Intn(8))}
				optimized.push(value)
				heap.Push(&reference, &value)
			} else if a, b := optimized.pop(), *heap.Pop(&reference).(*vertexDist); a != b {
				t.Fatalf("extraction changed: %v != %v", a, b)
			}
		}
		for len(optimized) != 0 {
			if a, b := optimized.pop(), *heap.Pop(&reference).(*vertexDist); a != b {
				t.Fatalf("drain changed: %v != %v", a, b)
			}
		}
	}
}

func TestQueryBuffersPreserveResultsAndOwnership(t *testing.T) {
	g := NewGraph()
	for i := int64(0); i < 8; i++ {
		if err := g.CreateVertex(i); err != nil {
			t.Fatal(err)
		}
	}
	for i := int64(0); i < 7; i++ {
		if err := g.AddEdge(i, i+1, float64(i%3)); err != nil {
			t.Fatal(err)
		}
		if err := g.AddEdge(i+1, i, float64((i+1)%3)); err != nil {
			t.Fatal(err)
		}
	}
	g.PrepareContractionHierarchies()
	pool := g.NewQueryPool()
	for source := int64(0); source < 8; source++ {
		for target := int64(0); target < 8; target++ {
			before, original := g.ShortestPath(source, target)
			cost, path := pool.ShortestPath(source, target)
			if cost != before || !reflect.DeepEqual(path, original) {
				t.Fatalf("result changed: %v %v != %v %v", cost, path, before, original)
			}
			snapshot := append([]int64(nil), path...)
			pool.ShortestPath(target, source)
			if !reflect.DeepEqual(snapshot, path) {
				t.Fatal("returned path aliases scratch")
			}
			starts := []VertexAlternative{{Label: source, AdditionalDistance: -2}}
			ends := []VertexAlternative{{Label: target, AdditionalDistance: 1}}
			before, original = g.ShortestPathWithAlternatives(starts, ends)
			cost, path = pool.ShortestPathWithAlternatives(starts, ends)
			if cost != before || !reflect.DeepEqual(path, original) {
				t.Fatalf("negative alternative changed: %v %v != %v %v", cost, path, before, original)
			}
		}
	}
}
