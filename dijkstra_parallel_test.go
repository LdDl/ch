package ch

import (
	"fmt"
	"math"
	"runtime"
	"testing"
	"time"
)

func TestParallelDijkstraPath(t *testing.T) {
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
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(0, 2, 1.0)
	g.AddEdge(0, 3, 1.0)
	g.AddEdge(0, 4, 1.0)
	g.AddEdge(0, 5, 1.0)
	g.AddEdge(0, 6, 1.0)
	g.AddEdge(0, 7, 1.0)
	g.AddEdge(8, 0, 1.0)
	g.AddEdge(9, 0, 1.0)
	g.AddEdge(1, 0, 1.0)
	g.AddEdge(2, 0, 1.0)
	g.AddEdge(3, 0, 1.0)
	g.AddEdge(4, 0, 1.0)
	g.AddEdge(5, 0, 1.0)
	g.AddEdge(6, 0, 1.0)
	g.AddEdge(7, 0, 1.0)
	g.AddEdge(7, 9, 1.0)
	g.AddEdge(6, 4, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(2, 3, 1.0)
	g.AddEdge(3, 4, 1.0)
	g.AddEdge(1, 3, 1.0)

	vertex := g.Vertices[0]
	inEdges := g.Vertices[0].inIncidentEdges
	outEdges := g.Vertices[0].outIncidentEdges
	contractID := int64(0)
	vertex.contracted = true
	inMax := 0.0
	outMax := 0.0
	for i := 0; i < len(inEdges); i++ {
		if g.Vertices[inEdges[i].vertexID].contracted {
			continue
		}
		if inMax < inEdges[i].cost {
			inMax = inEdges[i].cost
		}
	}
	for i := 0; i < len(outEdges); i++ {
		if g.Vertices[outEdges[i].vertexID].contracted {
			continue
		}
		if outMax < outEdges[i].cost {
			outMax = outEdges[i].cost
		}
	}
	max := inMax + outMax

	st := time.Now()
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if g.Vertices[inVertex].contracted {
			continue
		}
		incost := inEdges[i].cost
		g.dijkstra(inVertex, max, contractID, int64(i)) // Finds the shortest distances from the inVertex to all outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j].vertexID
			outcost := outEdges[j].cost
			if g.Vertices[outVertex].contracted {
				continue
			}
			summaryCost := incost + outcost
			if g.Vertices[outVertex].distance.contractID != contractID || g.Vertices[outVertex].distance.sourceID != int64(i) || g.Vertices[outVertex].distance.distance > summaryCost {
				g.createOrUpdateShortcut(inVertex, outVertex, vertex.vertexNum, summaryCost)
			}
		}
	}
	fmt.Println("Since single core", time.Since(st).Nanoseconds())

	st = time.Now()
	threads := runtime.GOMAXPROCS(0)
	g.pqComparators = make([]*distanceHeap, threads)
	res := make(chan bool, threads)
	limit := len(inEdges)
	lastIdx := 0
	for bn := 0; bn < threads; bn++ {
		go func(bn int, r chan<- bool) {
			start := (limit / threads) * bn
			end := start + (limit / threads)
			lastIdx = end
			edgesSet := inEdges[start:end]
			for i := 0; i < len(edgesSet); i++ {
				inVertex := edgesSet[i].vertexID
				if g.Vertices[inVertex].contracted {
					continue
				}
				incost := edgesSet[i].cost
				g.dijkstra_v2(inVertex, max, contractID, int64(i), bn) // Finds the shortest distances from the inVertex to all outVertices.
				for j := 0; j < len(outEdges); j++ {
					outVertex := outEdges[j].vertexID
					outcost := outEdges[j].cost
					if g.Vertices[outVertex].contracted {
						continue
					}
					summaryCost := incost + outcost
					if g.Vertices[outVertex].distance.contractID != contractID || g.Vertices[outVertex].distance.sourceID != int64(i) || g.Vertices[outVertex].distance.distance > summaryCost {
						g.createOrUpdateShortcut(inVertex, outVertex, vertex.vertexNum, summaryCost)
					}
				}
			}
			r <- true
		}(bn, res)
	}
	if lastIdx < len(inEdges) {
		edgesSet := inEdges[lastIdx:]
		for i := 0; i < len(edgesSet); i++ {
			inVertex := edgesSet[i].vertexID
			if g.Vertices[inVertex].contracted {
				continue
			}
			incost := edgesSet[i].cost
			g.dijkstra_v2(inVertex, max, contractID, int64(i), i) // Finds the shortest distances from the inVertex to all outVertices.
			for j := 0; j < len(outEdges); j++ {
				outVertex := outEdges[j].vertexID
				outcost := outEdges[j].cost
				if g.Vertices[outVertex].contracted {
					continue
				}
				summaryCost := incost + outcost
				if g.Vertices[outVertex].distance.contractID != contractID || g.Vertices[outVertex].distance.sourceID != int64(i) || g.Vertices[outVertex].distance.distance > summaryCost {
					g.createOrUpdateShortcut(inVertex, outVertex, vertex.vertexNum, summaryCost)
				}
			}
		}
	}
	fmt.Println("Since multi-core", time.Since(st).Nanoseconds())
	fmt.Println("DONE")
}

func BenchmarkSingleCore(b *testing.B) {
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
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(0, 2, 1.0)
	g.AddEdge(0, 3, 1.0)
	g.AddEdge(0, 4, 1.0)
	g.AddEdge(0, 5, 1.0)
	g.AddEdge(0, 6, 1.0)
	g.AddEdge(0, 7, 1.0)
	g.AddEdge(8, 0, 1.0)
	g.AddEdge(9, 0, 1.0)
	g.AddEdge(1, 0, 1.0)
	g.AddEdge(2, 0, 1.0)
	g.AddEdge(3, 0, 1.0)
	g.AddEdge(4, 0, 1.0)
	g.AddEdge(5, 0, 1.0)
	g.AddEdge(6, 0, 1.0)
	g.AddEdge(7, 0, 1.0)
	g.AddEdge(7, 9, 1.0)
	g.AddEdge(6, 4, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(2, 3, 1.0)
	g.AddEdge(3, 4, 1.0)
	g.AddEdge(1, 3, 1.0)

	vertex := g.Vertices[0]
	inEdges := g.Vertices[0].inIncidentEdges
	outEdges := g.Vertices[0].outIncidentEdges
	contractID := int64(0)
	vertex.contracted = true
	inMax := 0.0
	outMax := 0.0
	for i := 0; i < len(inEdges); i++ {
		if g.Vertices[inEdges[i].vertexID].contracted {
			continue
		}
		if inMax < inEdges[i].cost {
			inMax = inEdges[i].cost
		}
	}
	for i := 0; i < len(outEdges); i++ {
		if g.Vertices[outEdges[i].vertexID].contracted {
			continue
		}
		if outMax < outEdges[i].cost {
			outMax = outEdges[i].cost
		}
	}
	max := inMax + outMax

	b.Log("BenchmarkSingleCore is starting...")
	b.ResetTimer()

	for k := 0.; k <= 12; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%s/%d/vertices-%d", "CH shortest path", n, len(g.Vertices)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i := 0; i < len(inEdges); i++ {
					inVertex := inEdges[i].vertexID
					if g.Vertices[inVertex].contracted {
						continue
					}
					incost := inEdges[i].cost
					g.dijkstra(inVertex, max, contractID, int64(i)) // Finds the shortest distances from the inVertex to all outVertices.
					for j := 0; j < len(outEdges); j++ {
						outVertex := outEdges[j].vertexID
						outcost := outEdges[j].cost
						if g.Vertices[outVertex].contracted {
							continue
						}
						summaryCost := incost + outcost
						if g.Vertices[outVertex].distance.contractID != contractID || g.Vertices[outVertex].distance.sourceID != int64(i) || g.Vertices[outVertex].distance.distance > summaryCost {
							g.createOrUpdateShortcut(inVertex, outVertex, vertex.vertexNum, summaryCost)
						}
					}
				}
			}
		})
	}
}

func BenchmarkMultiCore(b *testing.B) {
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
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(0, 2, 1.0)
	g.AddEdge(0, 3, 1.0)
	g.AddEdge(0, 4, 1.0)
	g.AddEdge(0, 5, 1.0)
	g.AddEdge(0, 6, 1.0)
	g.AddEdge(0, 7, 1.0)
	g.AddEdge(8, 0, 1.0)
	g.AddEdge(9, 0, 1.0)
	g.AddEdge(1, 0, 1.0)
	g.AddEdge(2, 0, 1.0)
	g.AddEdge(3, 0, 1.0)
	g.AddEdge(4, 0, 1.0)
	g.AddEdge(5, 0, 1.0)
	g.AddEdge(6, 0, 1.0)
	g.AddEdge(7, 0, 1.0)
	g.AddEdge(7, 9, 1.0)
	g.AddEdge(6, 4, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(2, 3, 1.0)
	g.AddEdge(3, 4, 1.0)
	g.AddEdge(1, 3, 1.0)

	vertex := g.Vertices[0]
	inEdges := g.Vertices[0].inIncidentEdges
	outEdges := g.Vertices[0].outIncidentEdges
	contractID := int64(0)
	vertex.contracted = true
	inMax := 0.0
	outMax := 0.0
	for i := 0; i < len(inEdges); i++ {
		if g.Vertices[inEdges[i].vertexID].contracted {
			continue
		}
		if inMax < inEdges[i].cost {
			inMax = inEdges[i].cost
		}
	}
	for i := 0; i < len(outEdges); i++ {
		if g.Vertices[outEdges[i].vertexID].contracted {
			continue
		}
		if outMax < outEdges[i].cost {
			outMax = outEdges[i].cost
		}
	}
	max := inMax + outMax

	b.Log("BenchmarkSingleCore is starting...")
	b.ResetTimer()
	threads := runtime.GOMAXPROCS(0)
	g.pqComparators = make([]*distanceHeap, threads)
	for k := 0.; k <= 12; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%s/%d/vertices-%d", "CH shortest path", n, len(g.Vertices)), func(b *testing.B) {
			for bn := 0; bn < b.N; bn++ {
				res := make(chan bool, threads)
				limit := len(inEdges)
				lastIdx := 0
				for kf := 0; kf < threads; kf++ {
					go func(kf int, r chan<- bool) {
						start := (limit / threads) * kf
						end := start + (limit / threads)
						lastIdx = end
						edgesSet := inEdges[start:end]
						for i := 0; i < len(edgesSet); i++ {
							inVertex := edgesSet[i].vertexID
							if g.Vertices[inVertex].contracted {
								continue
							}
							incost := edgesSet[i].cost
							g.dijkstra_v2(inVertex, max, contractID, int64(i), kf) // Finds the shortest distances from the inVertex to all outVertices.
							for j := 0; j < len(outEdges); j++ {
								outVertex := outEdges[j].vertexID
								outcost := outEdges[j].cost
								if g.Vertices[outVertex].contracted {
									continue
								}
								summaryCost := incost + outcost
								if g.Vertices[outVertex].distance.contractID != contractID || g.Vertices[outVertex].distance.sourceID != int64(i) || g.Vertices[outVertex].distance.distance > summaryCost {
									g.createOrUpdateShortcut(inVertex, outVertex, vertex.vertexNum, summaryCost)
								}
							}
						}
						r <- true
					}(kf, res)
				}
				_ = lastIdx
				// if lastIdx < len(inEdges) {
				// 	fmt.Println("Done")
				// 	edgesSet := inEdges[lastIdx:]
				// 	_ = edgesSet
				// 	for i := 0; i < len(edgesSet); i++ {
				// 		// 	inVertex := edgesSet[i].vertexID
				// 		// 	if g.Vertices[inVertex].contracted {
				// 		// 		continue
				// 		// 	}
				// 		// 	incost := edgesSet[i].cost
				// 		// 	g.dijkstra_v2(inVertex, max, contractID, int64(i)) // Finds the shortest distances from the inVertex to all outVertices.
				// 		// 	for j := 0; j < len(outEdges); j++ {
				// 		// 		outVertex := outEdges[j].vertexID
				// 		// 		outcost := outEdges[j].cost
				// 		// 		if g.Vertices[outVertex].contracted {
				// 		// 			continue
				// 		// 		}
				// 		// 		summaryCost := incost + outcost
				// 		// 		if g.Vertices[outVertex].distance.contractID != contractID || g.Vertices[outVertex].distance.sourceID != int64(i) || g.Vertices[outVertex].distance.distance > summaryCost {
				// 		// 			g.createOrUpdateShortcut(inVertex, outVertex, vertex.vertexNum, summaryCost)
				// 		// 		}
				// 		// 	}
				// 	}
				// }
			}
		})
	}
}
