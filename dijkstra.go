package ch

import (
	"container/heap"
	"math"
)

// checkID Checks if both source's and target's contraction ID are not equal
func (graph *Graph) checkID(source, target int64) bool {
	return graph.Vertices[source].distance.contractID != graph.Vertices[target].distance.contractID || graph.Vertices[source].distance.sourceID != graph.Vertices[target].distance.sourceID
}

// relaxEdges Edge relaxation
func (graph *Graph) relaxEdges(vertex, contractID, sourceID int64) {
	vertexList := graph.Vertices[vertex].outIncidentEdges
	for i := 0; i < len(vertexList); i++ {
		temp := vertexList[i].vertexID
		cost := vertexList[i].cost
		// Skip shortcuts
		if graph.Vertices[temp].contracted {
			continue
		}
		if graph.checkID(vertex, temp) || graph.Vertices[temp].distance.distance > graph.Vertices[vertex].distance.distance+cost {
			graph.Vertices[temp].distance.distance = graph.Vertices[vertex].distance.distance + cost
			graph.Vertices[temp].distance.contractID = contractID
			graph.Vertices[temp].distance.sourceID = sourceID
			heap.Push(graph.pqComparator, graph.Vertices[temp])
		}
	}
}

// dijkstra Internal dijkstra algorithm to compute contraction hierarchies
func (graph *Graph) dijkstra(source int64, maxcost float64, contractID, sourceID int64) {
	// fmt.Println("DIJKSTRA START", source, maxcost)
	graph.pqComparator = &distanceHeap{}
	heap.Init(graph.pqComparator)
	heap.Push(graph.pqComparator, graph.Vertices[source])

	graph.Vertices[source].distance.distance = 0
	graph.Vertices[source].distance.contractID = contractID
	graph.Vertices[source].distance.sourceID = sourceID

	for graph.pqComparator.Len() != 0 {
		vertex := heap.Pop(graph.pqComparator).(*Vertex)
		if vertex.distance.distance > maxcost {
			// fmt.Println("Vertex done")
			return
		}
		// fmt.Println("Vertex relax", vertex.vertexNum, contractID, sourceID)
		graph.relaxEdges(vertex.vertexNum, contractID, sourceID)
	}
	// fmt.Println("DIJKSTRA END")
}

var (
	workset = make(map[int64]bool)
	dist_u  = make(map[int64]float64)
	hops_u  = make(map[int64]int)
	visited = make(map[int64]bool)
	pq      = &minheapSTD{}
)

func (graph *Graph) dijkstras_clear() {
	for i := range workset {
		dist_u[i] = math.MaxFloat64
		hops_u[i] = 0
		visited[i] = false
	}
}

func (graph *Graph) dijkstra_v2(source int64, max_tot float64) {
	graph.dijkstras_clear()
	dist_u[source] = 0
	hops_u[source] = 0
	visited[source] = true

	pq.add_with_priority(graph.Vertices[source].vertexNum, 0)
	if pq == nil {
		heap.Init(pq)
	}

	for pq.Len() != 0 {
		u := heap.Pop(pq).(minheapNode)

		u1 := u.id

		visited[u1] = true
		workset[u1] = true

		if dist_u[u1] > max_tot {
			break
		}
		if hops_u[u1] > 5 {
			continue
		}

		outgoing := graph.Vertices[u1].outIncidentEdges
		for i := range outgoing {
			v1 := outgoing[i].vertexID
			cost := outgoing[i].cost
			if graph.Vertices[v1].contracted || visited[v1] {
				continue
			}
			if dist_u[v1] > dist_u[u1]+cost || dist_u[v1] == math.MaxFloat64 {
				dist_u[v1] = dist_u[u1] + cost
				hops_u[v1] = 1 + hops_u[u1]
				pq.add_with_priority(v1, dist_u[v1])
				workset[v1] = true
			}
		}
	}
}
