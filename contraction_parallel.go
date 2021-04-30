package ch

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// PreprocessParallel Same as Preprocess() but with parallelism
func (graph *Graph) PreprocessParallel(numThreads int) []int64 {
	nodeOrdering := make([]int64, len(graph.Vertices))
	var extractNum int
	var iter int
	distChan := make(chan *distanceParallel, numThreads)
	for i := 0; i < numThreads; i++ {
		buf := createDistanceParallel(graph)
		distChan <- &buf
	}
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
		graph.contractNodeParallel(vertex, int64(extractNum-1), distChan, numThreads)
		if DEBUG_PREPROCESSING {
			if iter > 0 && graph.pqImportance.Len()%1000 == 0 {
				fmt.Printf("Contraction Order: %d / %d, Remain vertices in heap: %d. Currect shortcuts num: %d Time: %v\n", extractNum, len(graph.Vertices), graph.pqImportance.Len(), graph.shortcutsNum(), time.Now())
			}
		}
	}
	return nodeOrdering
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
							if dist.contract[outVertex] != dist.contractID || distance > incost+outcost {
								if _, ok := graph.shortcuts[inVertex]; !ok {
									// If there is no such shortcut then add one.
									graph.shortcuts[inVertex] = make(map[int64]*ContractionPath)
									graph.shortcuts[inVertex][outVertex] = &ContractionPath{
										ViaVertex: vertex.vertexNum,
										Cost:      summaryCost,
									}
									graph.Vertices[inVertex].outIncidentEdges = append(graph.Vertices[inVertex].outIncidentEdges, incidentEdge{outVertex, summaryCost})
									graph.Vertices[outVertex].inIncidentEdges = append(graph.Vertices[outVertex].inIncidentEdges, incidentEdge{inVertex, summaryCost})
								} else {
									if v, ok := graph.shortcuts[inVertex][outVertex]; ok {
										// If shortcut already exists
										// we should check if the middle vertex is still the same
										if v.ViaVertex == vertex.vertexNum {
											// If middle vertex is still the same then change cost of shortcut only [Additional conditional: previous estimated cost is less than current one]
											if summaryCost < graph.shortcuts[inVertex][outVertex].Cost {
												graph.shortcuts[inVertex][outVertex].Cost = summaryCost
												bk1 := graph.Vertices[inVertex].updateOutIncidentEdge(outVertex, summaryCost)
												if !bk1 {
													panic(fmt.Sprintf("Should not happen [1]. Can't update outcoming incident edge. %d has no common edge with %d", inVertex, outVertex))
												}
												bk2 := graph.Vertices[outVertex].updateInIncidentEdge(inVertex, summaryCost)
												if !bk2 {
													panic(fmt.Sprintf("Should not happen [2]. Can't update incoming incident edge. %d has no common edge with %d", outVertex, inVertex))
												}
												graph.Vertices[inVertex].outIncidentEdges = append(graph.Vertices[inVertex].outIncidentEdges, incidentEdge{outVertex, summaryCost})
												graph.Vertices[outVertex].inIncidentEdges = append(graph.Vertices[outVertex].inIncidentEdges, incidentEdge{inVertex, summaryCost})
											}
										} else {
											// If middle vertex is not optimal for shortcut then change both vertex ID and cost [Additional conditional: previous estimated cost is less than current one]
											if summaryCost < graph.shortcuts[inVertex][outVertex].Cost {
												graph.shortcuts[inVertex][outVertex].ViaVertex = vertex.vertexNum
												graph.shortcuts[inVertex][outVertex].Cost = summaryCost
												dk1 := graph.Vertices[inVertex].deleteOutIncidentEdge(outVertex)
												if !dk1 {
													panic(fmt.Sprintf("Should not happen [3]. Can't delete outcoming incident edge. %d has no common edge with %d", inVertex, outVertex))
												}
												dk2 := graph.Vertices[outVertex].deleteInIncidentEdge(inVertex)
												if !dk2 {
													panic(fmt.Sprintf("Should not happen [4]. Can't delete incoming incident edge. %d has no common edge with %d", outVertex, inVertex))
												}
												graph.Vertices[inVertex].outIncidentEdges = append(graph.Vertices[inVertex].outIncidentEdges, incidentEdge{outVertex, summaryCost})
												graph.Vertices[outVertex].inIncidentEdges = append(graph.Vertices[outVertex].inIncidentEdges, incidentEdge{inVertex, summaryCost})
											}
										}
									} else {
										graph.shortcuts[inVertex][outVertex] = &ContractionPath{
											ViaVertex: vertex.vertexNum,
											Cost:      summaryCost,
										}
										graph.Vertices[inVertex].outIncidentEdges = append(graph.Vertices[inVertex].outIncidentEdges, incidentEdge{outVertex, summaryCost})
										graph.Vertices[outVertex].inIncidentEdges = append(graph.Vertices[outVertex].inIncidentEdges, incidentEdge{inVertex, summaryCost})
									}
								}
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
