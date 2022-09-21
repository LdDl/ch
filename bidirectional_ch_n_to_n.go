package ch

import (
	"container/heap"
)

// ShortestPath Computes and returns shortest path and it's cost (extended Dijkstra's algorithm)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source User's definied ID of source vertex
// target User's definied ID of target vertex
func (graph *Graph) ShortestPathManyToMany(sources, targets []int64) ([][]float64, [][][]int64) {
	endpoints := [directionsCount][]int64{sources, targets}
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

func (graph *Graph) initShortestPathManyToMany(endpoints [directionsCount][]int64) (
	queryDist [directionsCount][][]float64,
	processed [directionsCount][][]bool,
	queues [directionsCount][]*vertexDistHeap,
) {
	for d := forward; d < directionsCount; d++ {
		queryDist[d] = make([][]float64, len(endpoints[d]))
		processed[d] = make([][]bool, len(endpoints[d]))
		queues[d] = make([]*vertexDistHeap, len(endpoints[d]))
		for endpointIdx := range endpoints[d] {
			queryDist[d][endpointIdx] = make([]float64, len(graph.Vertices))

			for i := range queryDist[d][endpointIdx] {
				queryDist[d][endpointIdx][i] = Infinity
			}

			processed[d][endpointIdx] = make([]bool, len(graph.Vertices))

			queues[d][endpointIdx] = &vertexDistHeap{}

			heap.Init(queues[d][endpointIdx])
		}
	}

	return
}

func (graph *Graph) shortestPathManyToMany(endpoints [directionsCount][]int64) ([][]float64, [][][]int64) {
	queryDist, processed, queues := graph.initShortestPathManyToMany(endpoints)

	for d := forward; d < directionsCount; d++ {

		for endpointIdx, endpoint := range endpoints[d] {
			processed[d][endpointIdx][endpoint] = true

			queryDist[d][endpointIdx][endpoint] = 0

			heapEndpoint := &vertexDist{
				id:   endpoint,
				dist: 0,
			}
			heap.Push(queues[d][endpointIdx], heapEndpoint)
		}
	}

	return graph.shortestPathManyToManyCore(queryDist, processed, queues)
}

func (graph *Graph) shortestPathManyToManyCore(
	queryDist [directionsCount][][]float64,
	processed [directionsCount][][]bool,
	queues [directionsCount][]*vertexDistHeap,
) ([][]float64, [][][]int64) {
	var prev [directionsCount][]map[int64]int64
	for d := forward; d < directionsCount; d++ {
		prev[d] = make([]map[int64]int64, len(queues[d]))
		for endpointIdx := range queues[d] {
			prev[d][endpointIdx] = make(map[int64]int64)
		}
	}
	estimates := make([][]float64, len(queues[forward]))
	middleIDs := make([][]int64, len(queues[forward]))
	for sourceEndpointIdx := range queues[forward] {
		sourceEstimates := make([]float64, len(queues[backward]))
		sourceMiddleIDs := make([]int64, len(queues[backward]))
		estimates[sourceEndpointIdx] = sourceEstimates
		middleIDs[sourceEndpointIdx] = sourceMiddleIDs
		for targetEndpointIdx := range queues[backward] {
			sourceEstimates[targetEndpointIdx] = Infinity
			sourceMiddleIDs[targetEndpointIdx] = int64(-1)
		}
	}

	for {
		queuesProcessed := false
		for d := forward; d < directionsCount; d++ {
			reverseDirection := (d + 1) % directionsCount
			for endpointIdx := range queues[d] {
				if queues[d][endpointIdx].Len() == 0 {
					continue
				}
				queuesProcessed = true
				graph.directionalSearchManyToMany(d, endpointIdx, queues[d][endpointIdx], processed[d][endpointIdx], processed[reverseDirection], queryDist[d][endpointIdx], queryDist[reverseDirection], prev[d][endpointIdx], estimates, middleIDs)
			}
		}
		if !queuesProcessed {
			break
		}
	}
	paths := make([][][]int64, len(estimates))
	for sourceEndpointIdx, targetEstimates := range estimates {
		targetPaths := make([][]int64, len(targetEstimates))
		paths[sourceEndpointIdx] = targetPaths
		for targetEndpointIdx, estimate := range targetEstimates {
			if estimate == Infinity {
				targetEstimates[targetEndpointIdx] = -1
				continue
			}
			targetPaths[targetEndpointIdx] = graph.ComputePath(middleIDs[sourceEndpointIdx][targetEndpointIdx], prev[forward][sourceEndpointIdx], prev[backward][targetEndpointIdx])
		}
	}
	return estimates, paths
}

func (graph *Graph) directionalSearchManyToMany(
	d direction, endpointIndex int, q *vertexDistHeap,
	localProcessed []bool, reverseProcessed [][]bool,
	localQueryDist []float64, reverseQueryDist [][]float64,
	prev map[int64]int64, estimates [][]float64, middleIDs [][]int64) {

	vertex := heap.Pop(q).(*vertexDist)
	// if vertex.dist <= *estimate { // TODO: move to another place
	localProcessed[vertex.id] = true
	// Edge relaxation in a forward propagation
	var vertexList []*incidentEdge
	if d == forward {
		vertexList = graph.Vertices[vertex.id].outIncidentEdges
	} else {
		vertexList = graph.Vertices[vertex.id].inIncidentEdges
	}
	for i := range vertexList {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if graph.Vertices[vertex.id].orderPos < graph.Vertices[temp].orderPos {
			alt := localQueryDist[vertex.id] + cost
			if localQueryDist[temp] > alt {
				localQueryDist[temp] = alt
				prev[temp] = vertex.id
				node := &vertexDist{
					id:   temp,
					dist: alt,
				}
				heap.Push(q, node)
			}
		}
	}
	// }
	for revEndpointIdx, revEndpointProcessed := range reverseProcessed {
		if revEndpointProcessed[vertex.id] {
			var sourceEndpoint, targetEndpoint int
			if d == forward {
				sourceEndpoint, targetEndpoint = endpointIndex, revEndpointIdx
			} else {
				targetEndpoint, sourceEndpoint = endpointIndex, revEndpointIdx
			}
			if vertex.dist+reverseQueryDist[revEndpointIdx][vertex.id] < estimates[sourceEndpoint][targetEndpoint] {
				middleIDs[sourceEndpoint][targetEndpoint] = vertex.id
				estimates[sourceEndpoint][targetEndpoint] = vertex.dist + reverseQueryDist[revEndpointIdx][vertex.id]
			}
		}
	}
}
