package ch

import (
	"container/heap"
)

// ShortestPathManyToMany computes and returns shortest paths and theirs's costs (extended Dijkstra's algorithm) between multiple sources and targets
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// sources - set of user's definied IDs of source vertices
// targets - set of user's definied IDs of target vertices
func (graph *Graph) ShortestPathManyToMany(sources, targets []int64) ([][]float64, [][][]int64) {
	// Copy input slices to avoid modifying caller's data
	sourcesCopy := make([]int64, len(sources))
	targetsCopy := make([]int64, len(targets))
	copy(sourcesCopy, sources)
	copy(targetsCopy, targets)

	endpoints := [directionsCount][]int64{sourcesCopy, targetsCopy}
	for d, directionEndpoints := range endpoints {
		for i, endpoint := range directionEndpoints {
			var ok bool
			if endpoints[d][i], ok = graph.mapping[endpoint]; !ok {
				endpoints[d][i] = -1
			}
		}
	}
	return graph.shortestPathManyToMany(endpoints)
}

// initManyToManyBuffers ensures buffers are allocated and properly sized for the query
func (graph *Graph) initManyToManyBuffers(numSources, numTargets int) {
	n := len(graph.Vertices)
	endpointCounts := [directionsCount]int{numSources, numTargets}

	for d := forward; d < directionsCount; d++ {
		// Grow outer slices if needed
		if len(graph.manyToManyDist[d]) < endpointCounts[d] {
			newDist := make([][]float64, endpointCounts[d])
			newEpochs := make([][]int64, endpointCounts[d])
			newPrev := make([]map[int64]int64, endpointCounts[d])
			copy(newDist, graph.manyToManyDist[d])
			copy(newEpochs, graph.manyToManyEpochs[d])
			copy(newPrev, graph.manyToManyPrev[d])
			graph.manyToManyDist[d] = newDist
			graph.manyToManyEpochs[d] = newEpochs
			graph.manyToManyPrev[d] = newPrev
		}

		// Ensure each endpoint has properly sized inner slices
		for i := 0; i < endpointCounts[d]; i++ {
			if len(graph.manyToManyDist[d][i]) < n {
				graph.manyToManyDist[d][i] = make([]float64, n)
				graph.manyToManyEpochs[d][i] = make([]int64, n)
			}
			if graph.manyToManyPrev[d][i] == nil {
				graph.manyToManyPrev[d][i] = make(map[int64]int64)
			}
		}
	}
}

// getManyToManyDist returns distance for endpoint at given vertex, using epoch for lazy clearing
func (graph *Graph) getManyToManyDist(d direction, endpointIdx int, vertexID int64) float64 {
	if graph.manyToManyEpochs[d][endpointIdx][vertexID] != graph.manyToManyEpoch {
		return Infinity
	}
	return graph.manyToManyDist[d][endpointIdx][vertexID]
}

// setManyToManyDist sets distance for endpoint at given vertex
func (graph *Graph) setManyToManyDist(d direction, endpointIdx int, vertexID int64, dist float64) {
	graph.manyToManyDist[d][endpointIdx][vertexID] = dist
	graph.manyToManyEpochs[d][endpointIdx][vertexID] = graph.manyToManyEpoch
}

func (graph *Graph) shortestPathManyToMany(endpoints [directionsCount][]int64) ([][]float64, [][][]int64) {
	numSources := len(endpoints[forward])
	numTargets := len(endpoints[backward])

	// Increment epoch for lazy buffer clearing
	graph.manyToManyEpoch++
	graph.initManyToManyBuffers(numSources, numTargets)

	// Clear prev maps (these still need explicit clearing)
	for d := forward; d < directionsCount; d++ {
		endpointCount := numSources
		if d == backward {
			endpointCount = numTargets
		}
		for i := 0; i < endpointCount; i++ {
			for k := range graph.manyToManyPrev[d][i] {
				delete(graph.manyToManyPrev[d][i], k)
			}
		}
	}

	// Initialize queues
	queues := [directionsCount][]*vertexDistHeap{}
	for d := forward; d < directionsCount; d++ {
		endpointCount := numSources
		if d == backward {
			endpointCount = numTargets
		}
		queues[d] = make([]*vertexDistHeap, endpointCount)
		for i := 0; i < endpointCount; i++ {
			queues[d][i] = &vertexDistHeap{}
			heap.Init(queues[d][i])
		}
	}

	// Initialize sources and targets
	for d := forward; d < directionsCount; d++ {
		for endpointIdx, endpoint := range endpoints[d] {
			if endpoint == -1 {
				continue
			}
			graph.setManyToManyDist(d, endpointIdx, endpoint, 0)
			heap.Push(queues[d][endpointIdx], &vertexDist{id: endpoint, dist: 0})
		}
	}

	// Initialize estimates matrix
	estimates := make([][]float64, numSources)
	middleIDs := make([][]int64, numSources)
	for i := 0; i < numSources; i++ {
		estimates[i] = make([]float64, numTargets)
		middleIDs[i] = make([]int64, numTargets)
		for j := 0; j < numTargets; j++ {
			estimates[i][j] = Infinity
			middleIDs[i][j] = -1
		}
	}

	// Main search loop
	for {
		queuesProcessed := false
		for d := forward; d < directionsCount; d++ {
			endpointCount := numSources
			if d == backward {
				endpointCount = numTargets
			}
			for endpointIdx := 0; endpointIdx < endpointCount; endpointIdx++ {
				if queues[d][endpointIdx].Len() == 0 {
					continue
				}
				queuesProcessed = true
				graph.directionalSearchManyToMany(d, endpointIdx, queues, estimates, middleIDs, numSources, numTargets)
			}
		}
		if !queuesProcessed {
			break
		}
	}

	// Build paths
	paths := make([][][]int64, numSources)
	for sourceIdx := 0; sourceIdx < numSources; sourceIdx++ {
		paths[sourceIdx] = make([][]int64, numTargets)
		for targetIdx := 0; targetIdx < numTargets; targetIdx++ {
			if estimates[sourceIdx][targetIdx] == Infinity {
				estimates[sourceIdx][targetIdx] = -1
				continue
			}
			paths[sourceIdx][targetIdx] = graph.ComputePath(
				middleIDs[sourceIdx][targetIdx],
				graph.manyToManyPrev[forward][sourceIdx],
				graph.manyToManyPrev[backward][targetIdx],
			)
		}
	}

	return estimates, paths
}

func (graph *Graph) directionalSearchManyToMany(d direction, endpointIdx int, queues [directionsCount][]*vertexDistHeap, estimates [][]float64, middleIDs [][]int64, numSources, numTargets int) {
	q := queues[d][endpointIdx]
	vertex := heap.Pop(q).(*vertexDist)

	// Skip if we've already found a better path
	currentDist := graph.getManyToManyDist(d, endpointIdx, vertex.id)
	if vertex.dist > currentDist {
		return
	}

	// Edge relaxation
	var vertexList []incidentEdge
	if d == forward {
		vertexList = graph.Vertices[vertex.id].outIncidentEdges
	} else {
		vertexList = graph.Vertices[vertex.id].inIncidentEdges
	}

	for i := range vertexList {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight

		// Only explore upward in CH
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := vertex.dist + cost
			tempDist := graph.getManyToManyDist(d, endpointIdx, temp)
			if alt < tempDist {
				graph.setManyToManyDist(d, endpointIdx, temp, alt)
				graph.manyToManyPrev[d][endpointIdx][temp] = vertex.id
				heap.Push(q, &vertexDist{id: temp, dist: alt})
			}
		}
	}

	// Check for meeting points with reverse direction
	reverseEndpointCount := numTargets
	if d == backward {
		reverseEndpointCount = numSources
	}

	for revIdx := 0; revIdx < reverseEndpointCount; revIdx++ {
		revDist := graph.getManyToManyDist(1-d, revIdx, vertex.id)
		if revDist == Infinity {
			continue
		}

		var sourceIdx, targetIdx int
		if d == forward {
			sourceIdx, targetIdx = endpointIdx, revIdx
		} else {
			sourceIdx, targetIdx = revIdx, endpointIdx
		}

		newEstimate := vertex.dist + revDist
		if newEstimate < estimates[sourceIdx][targetIdx] {
			estimates[sourceIdx][targetIdx] = newEstimate
			middleIDs[sourceIdx][targetIdx] = vertex.id
		}
	}
}

// ShortestPathManyToManyWithAlternatives Computes and returns shortest paths and their cost (extended Dijkstra's algorithm),
// with multiple alternatives for source and target vertices with additional distances to reach the vertices
// (useful if source and target are outside of the graph)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// sourcesAlternatives - set of user's definied IDs of source vertices with additional penalty
// targetsAlternatives - set of user's definied IDs of target vertices  with additional penalty
func (graph *Graph) ShortestPathManyToManyWithAlternatives(sourcesAlternatives, targetsAlternatives [][]VertexAlternative) ([][]float64, [][][]int64) {
	endpoints := [directionsCount][][]VertexAlternative{sourcesAlternatives, targetsAlternatives}
	var endpointsInternal [directionsCount][][]vertexAlternativeInternal
	for d, directionEndpoints := range endpoints {
		endpointsInternal[d] = make([][]vertexAlternativeInternal, 0, len(directionEndpoints))
		for _, alternatives := range directionEndpoints {
			endpointsInternal[d] = append(endpointsInternal[d], graph.vertexAlternativesToInternal(alternatives))
		}
	}
	return graph.shortestPathManyToManyWithAlternatives(endpointsInternal)
}

func (graph *Graph) shortestPathManyToManyWithAlternatives(endpoints [directionsCount][][]vertexAlternativeInternal) ([][]float64, [][][]int64) {
	numSources := len(endpoints[forward])
	numTargets := len(endpoints[backward])

	// Increment epoch for lazy buffer clearing
	graph.manyToManyEpoch++
	graph.initManyToManyBuffers(numSources, numTargets)

	// Clear prev maps
	for d := forward; d < directionsCount; d++ {
		endpointCount := numSources
		if d == backward {
			endpointCount = numTargets
		}
		for i := 0; i < endpointCount; i++ {
			for k := range graph.manyToManyPrev[d][i] {
				delete(graph.manyToManyPrev[d][i], k)
			}
		}
	}

	// Initialize queues
	queues := [directionsCount][]*vertexDistHeap{}
	for d := forward; d < directionsCount; d++ {
		endpointCount := numSources
		if d == backward {
			endpointCount = numTargets
		}
		queues[d] = make([]*vertexDistHeap, endpointCount)
		for i := 0; i < endpointCount; i++ {
			queues[d][i] = &vertexDistHeap{}
			heap.Init(queues[d][i])
		}
	}

	// Initialize with alternatives
	for d := forward; d < directionsCount; d++ {
		for endpointIdx, alternatives := range endpoints[d] {
			for _, alt := range alternatives {
				if alt.vertexNum == vertexNotFound {
					continue
				}
				currentDist := graph.getManyToManyDist(d, endpointIdx, alt.vertexNum)
				if alt.additionalDistance < currentDist {
					graph.setManyToManyDist(d, endpointIdx, alt.vertexNum, alt.additionalDistance)
					heap.Push(queues[d][endpointIdx], &vertexDist{id: alt.vertexNum, dist: alt.additionalDistance})
				}
			}
		}
	}

	// Initialize estimates matrix
	estimates := make([][]float64, numSources)
	middleIDs := make([][]int64, numSources)
	for i := 0; i < numSources; i++ {
		estimates[i] = make([]float64, numTargets)
		middleIDs[i] = make([]int64, numTargets)
		for j := 0; j < numTargets; j++ {
			estimates[i][j] = Infinity
			middleIDs[i][j] = -1
		}
	}

	// Main search loop
	for {
		queuesProcessed := false
		for d := forward; d < directionsCount; d++ {
			endpointCount := numSources
			if d == backward {
				endpointCount = numTargets
			}
			for endpointIdx := 0; endpointIdx < endpointCount; endpointIdx++ {
				if queues[d][endpointIdx].Len() == 0 {
					continue
				}
				queuesProcessed = true
				graph.directionalSearchManyToMany(d, endpointIdx, queues, estimates, middleIDs, numSources, numTargets)
			}
		}
		if !queuesProcessed {
			break
		}
	}

	// Build paths
	paths := make([][][]int64, numSources)
	for sourceIdx := 0; sourceIdx < numSources; sourceIdx++ {
		paths[sourceIdx] = make([][]int64, numTargets)
		for targetIdx := 0; targetIdx < numTargets; targetIdx++ {
			if estimates[sourceIdx][targetIdx] == Infinity {
				estimates[sourceIdx][targetIdx] = -1
				continue
			}
			paths[sourceIdx][targetIdx] = graph.ComputePath(
				middleIDs[sourceIdx][targetIdx],
				graph.manyToManyPrev[forward][sourceIdx],
				graph.manyToManyPrev[backward][targetIdx],
			)
		}
	}

	return estimates, paths
}
