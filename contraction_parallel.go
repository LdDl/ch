package ch

import (
	"container/heap"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// PreprocessParallel Same as Preprocess() but with parallelism
func (graph *Graph) PreprocessParallel() []int64 {
	numThreads := runtime.GOMAXPROCS(0)
	nodeOrdering := make([]int64, len(graph.Vertices))
	var extractNum int
	var iter int

	// distChan := make(chan *distanceParallel, numThreads)
	// for i := 0; i < numThreads; i++ {
	// 	buf := createDistanceParallel(graph)
	// 	distChan <- &buf
	// }

	for graph.pqImportance.Len() != 0 {
		iter++
		vertex := heap.Pop(graph.pqImportance).(*Vertex)
		vertex.computeImportance()
		if graph.pqImportance.Len() != 0 && vertex.importance > graph.pqImportance.Peek().importance {
			graph.pqImportance.Push(vertex)
			continue
		}

		nodeOrdering[extractNum] = vertex.vertexNum
		vertex.orderPos = extractNum
		extractNum = extractNum + 1
		// graph.contractNodeParallel(vertex, int64(extractNum-1), distChan, numThreads)
		graph.contractNodeParallel_v2(vertex, int64(extractNum-1), numThreads)
		if DEBUG_PREPROCESSING {
			if iter > 0 && graph.pqImportance.Len()%1000 == 0 {
				fmt.Printf("Contraction Order: %d / %d, Remain vertices in heap: %d. Currect shortcuts num: %d Time: %v\n", extractNum, len(graph.Vertices), graph.pqImportance.Len(), graph.shortcutsNum(), time.Now())
			}
		}
	}
	return nodeOrdering
}

// contractNodeParallel_v2 Same as contractNode() but with but with parallelism
func (graph *Graph) contractNodeParallel_v2(vertex *Vertex, contractID int64, threads int) {
	inEdges := vertex.inIncidentEdges
	outEdges := vertex.outIncidentEdges
	vertex.contracted = true
	inMax := 0.0
	outMax := 0.0
	graph.markNeighbors(inEdges, outEdges)
	for i := 0; i < len(inEdges); i++ {
		if graph.Vertices[inEdges[i].vertexID].contracted {
			continue
		}
		if inMax < inEdges[i].cost {
			inMax = inEdges[i].cost
		}
	}
	for i := 0; i < len(outEdges); i++ {
		if graph.Vertices[outEdges[i].vertexID].contracted {
			continue
		}
		if outMax < outEdges[i].cost {
			outMax = outEdges[i].cost
		}
	}
	max := inMax + outMax

	res := make(chan bool, threads)
	limit := len(inEdges)
	lastIdx := 0
	done := false
	graph.pqComparators = make([]*distanceHeap, threads)
	if threads < limit {
		fmt.Printf("Here threading for %d pathes\n", limit)
		for i := 0; i < threads; i++ {
			// fmt.Println("Start thread #", i)
			go func(i int, r chan<- bool) {
				start := (limit / threads) * i
				end := start + (limit / threads)
				edgesSet := inEdges[start:end]

				lastIdx = end
				// fmt.Println(i, start, end, limit, threads, edgesSet)
				graph.workWithIncidentEdges(edgesSet, outEdges, max, contractID, vertex.vertexNum, i)
				// fmt.Println("Done thread #", i)
				r <- true
			}(i, res)
		}
		for i := 0; i < threads; i++ {
			done = <-res
		}
		_ = done
		if lastIdx < len(inEdges) {
			graph.workWithIncidentEdges(inEdges[lastIdx:], outEdges, max, contractID, vertex.vertexNum, 0)
		}
	} else {
		graph.workWithIncidentEdgesSingle(inEdges, outEdges, max, contractID, vertex.vertexNum)
	}

	// fmt.Println("`231")
	// panic("done")
}

func (graph *Graph) workWithIncidentEdges(inEdges []incidentEdge, outEdges []incidentEdge, max float64, contractID, vertexID int64, threadID int) {
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inEdges[i].cost
		graph.dijkstra_v2(inVertex, max, contractID, int64(i), threadID) // Finds the shortest distances from the inVertex to all outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j].vertexID
			outcost := outEdges[j].cost
			if graph.Vertices[outVertex].contracted {
				continue
			}
			summaryCost := incost + outcost
			if graph.Vertices[outVertex].distance.contractID != contractID || graph.Vertices[outVertex].distance.sourceID != int64(i) || graph.Vertices[outVertex].distance.distance > summaryCost {
				graph.createOrUpdateShortcut(inVertex, outVertex, vertexID, summaryCost)
			}
		}
	}
}

// contractNodeParallel Same as contractNode() but with but with parallelism
func (graph *Graph) contractNodeParallel(vertex *Vertex, contractID int64, distChan chan *distanceParallel, threads int) {

	inEdges := vertex.inIncidentEdges
	outEdges := vertex.outIncidentEdges

	vertex.contracted = true

	inMax := 0.0
	outMax := 0.0

	outChan := make(chan *distanceParallel, threads)
	graph.markNeighbors(inEdges, outEdges)

	for i := 0; i < len(inEdges); i++ {
		if graph.Vertices[inEdges[i].vertexID].contracted {
			continue
		}
		if inMax < inEdges[i].cost {
			inMax = inEdges[i].cost
		}
	}

	for i := 0; i < len(outEdges); i++ {
		if graph.Vertices[outEdges[i].vertexID].contracted {
			continue
		}
		if outMax < outEdges[i].cost {
			outMax = outEdges[i].cost
		}
	}

	max := inMax + outMax

	wg := sync.WaitGroup{}
	counter := 0
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		incost := inEdges[i].cost
		if !graph.Vertices[inVertex].contracted {
			wg.Add(1)
			counter++
			go graph.dijkstraParallel(inVertex, i, int(contractID)+counter, max, distChan, outChan, &wg)
		}
		if ((counter+1)%threads == 0) || i == len(inEdges)-1 {
			wg.Wait()
		outloop:
			for {
				var dist *distanceParallel
				select {
				case dist = <-outChan:
					{
						inVertex = dist.sourceID
						incost = inEdges[dist.edgeID].cost
						for j := 0; j < len(outEdges); j++ {
							outVertex := outEdges[j].vertexID
							outcost := outEdges[j].cost
							if graph.Vertices[outVertex].contracted {
								continue
							}
							distance := dist.distance[outVertex]
							summaryCost := incost + outcost
							if dist.contract[outVertex] != dist.contractID || distance > summaryCost {
								graph.createOrUpdateShortcut(inVertex, outVertex, vertex.vertexNum, summaryCost)
							}
						}
						distChan <- dist
					}
				default:
					break outloop
				}
			}
		}
	}
}

func (graph *Graph) workWithIncidentEdgesSingle(inEdges []incidentEdge, outEdges []incidentEdge, max float64, contractID, vertexID int64) {
	for i := 0; i < len(inEdges); i++ {
		inVertex := inEdges[i].vertexID
		if graph.Vertices[inVertex].contracted {
			continue
		}
		incost := inEdges[i].cost
		graph.dijkstra(inVertex, max, contractID, int64(i)) // Finds the shortest distances from the inVertex to all outVertices.
		for j := 0; j < len(outEdges); j++ {
			outVertex := outEdges[j].vertexID
			outcost := outEdges[j].cost
			if graph.Vertices[outVertex].contracted {
				continue
			}
			summaryCost := incost + outcost
			if graph.Vertices[outVertex].distance.contractID != contractID || graph.Vertices[outVertex].distance.sourceID != int64(i) || graph.Vertices[outVertex].distance.distance > summaryCost {
				graph.createOrUpdateShortcut(inVertex, outVertex, vertexID, summaryCost)
			}
		}
	}
}
