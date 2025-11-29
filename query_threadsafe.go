package ch

import (
	"container/heap"
	"sync"
)

// QueryState holds all the buffers needed for a single shortest path query.
// This is used by the thread-safe query methods to avoid sharing state between goroutines.
type QueryState struct {
	// Current query epoch (incremented each query within this state)
	epoch int64
	// Distance arrays for bidirectional search
	dist [directionsCount][]float64
	// Epoch markers (if != epoch, distance is Infinity)
	epochs [directionsCount][]int64
	// Previous vertex maps for path reconstruction
	prev [directionsCount]map[int64]int64
	// Priority queues for bidirectional search
	queues [directionsCount]*vertexDistHeap
}

// QueryPool provides thread-safe access to pooled QueryState objects.
// Use this when you need to call shortest path queries from multiple goroutines.
type QueryPool struct {
	pool  sync.Pool
	graph *Graph
}

// NewQueryPool creates a new QueryPool for concurrent query execution.
// The pool lazily initializes QueryState objects as needed.
func (graph *Graph) NewQueryPool() *QueryPool {
	return &QueryPool{
		graph: graph,
		pool: sync.Pool{
			New: func() interface{} {
				return &QueryState{}
			},
		},
	}
}

// acquireState gets a QueryState from the pool and initializes it if needed
func (qp *QueryPool) acquireState() *QueryState {
	state := qp.pool.Get().(*QueryState)
	n := len(qp.graph.Vertices)

	// Lazy initialization of buffers (only on first use of this state)
	if state.dist[forward] == nil || len(state.dist[forward]) != n {
		for d := forward; d < directionsCount; d++ {
			state.dist[d] = make([]float64, n)
			state.epochs[d] = make([]int64, n)
			state.prev[d] = make(map[int64]int64)
			state.queues[d] = &vertexDistHeap{}
		}
	}

	// Increment epoch to invalidate previous distances
	state.epoch++

	// Clear prev maps
	for d := forward; d < directionsCount; d++ {
		state.prev[d] = make(map[int64]int64)
		state.queues[d] = &vertexDistHeap{}
		heap.Init(state.queues[d])
	}

	return state
}

// releaseState returns a QueryState to the pool
func (qp *QueryPool) releaseState(state *QueryState) {
	qp.pool.Put(state)
}

// ShortestPath computes shortest path using a pooled QueryState (thread-safe).
// This method can be safely called from multiple goroutines concurrently.
//
// source - user's defined ID of source vertex
// target - user's defined ID of target vertex
func (qp *QueryPool) ShortestPath(source, target int64) (float64, []int64) {
	if source == target {
		return 0, []int64{source}
	}

	endpoints := [directionsCount]int64{source, target}
	for d, endpoint := range endpoints {
		var ok bool
		if endpoints[d], ok = qp.graph.mapping[endpoint]; !ok {
			return -1.0, nil
		}
	}

	state := qp.acquireState()
	defer qp.releaseState(state)

	return qp.shortestPath(state, endpoints)
}

func (qp *QueryPool) shortestPath(state *QueryState, endpoints [directionsCount]int64) (float64, []int64) {
	for d := forward; d < directionsCount; d++ {
		state.epochs[d][endpoints[d]] = state.epoch
		state.dist[d][endpoints[d]] = 0
		heapEndpoint := &vertexDist{
			id:   endpoints[d],
			dist: 0,
		}
		heap.Push(state.queues[d], heapEndpoint)
	}
	return qp.shortestPathCore(state)
}

func (qp *QueryPool) shortestPathCore(state *QueryState) (float64, []int64) {
	estimate := Infinity
	middleID := int64(-1)

	for {
		queuesProcessed := false
		for d := forward; d < directionsCount; d++ {
			if state.queues[d].Len() == 0 {
				continue
			}
			queuesProcessed = true
			reverseDirection := (d + 1) % directionsCount
			qp.directionalSearch(state, d, reverseDirection, &estimate, &middleID)
		}
		if !queuesProcessed {
			break
		}
	}

	if estimate == Infinity {
		return -1.0, nil
	}
	return estimate, qp.graph.ComputePath(middleID, state.prev[forward], state.prev[backward])
}

func (qp *QueryPool) directionalSearch(state *QueryState, d direction, reverseDirection direction, estimate *float64, middleID *int64) {
	vertex := heap.Pop(state.queues[d]).(*vertexDist)
	if vertex.dist <= *estimate {
		state.epochs[d][vertex.id] = state.epoch
		// Edge relaxation
		var vertexList []incidentEdge
		if d == forward {
			vertexList = qp.graph.Vertices[vertex.id].outIncidentEdges
		} else {
			vertexList = qp.graph.Vertices[vertex.id].inIncidentEdges
		}
		for i := range vertexList {
			temp := vertexList[i].vertexID
			cost := vertexList[i].weight
			if qp.graph.Vertices[vertex.id].orderPos < qp.graph.Vertices[temp].orderPos {
				alt := state.dist[d][vertex.id] + cost
				if state.epochs[d][temp] != state.epoch || state.dist[d][temp] > alt {
					state.dist[d][temp] = alt
					state.epochs[d][temp] = state.epoch
					state.prev[d][temp] = vertex.id
					node := &vertexDist{
						id:   temp,
						dist: alt,
					}
					heap.Push(state.queues[d], node)
				}
			}
		}
	}
	// Check if reverse direction has processed this vertex
	if state.epochs[reverseDirection][vertex.id] == state.epoch {
		if vertex.dist+state.dist[reverseDirection][vertex.id] < *estimate {
			*middleID = vertex.id
			*estimate = vertex.dist + state.dist[reverseDirection][vertex.id]
		}
	}
}

// ShortestPathWithAlternatives computes shortest path with multiple source/target alternatives (thread-safe).
// This method can be safely called from multiple goroutines concurrently.
//
// sources - user's defined source vertices with additional penalties
// targets - user's defined target vertices with additional penalties
func (qp *QueryPool) ShortestPathWithAlternatives(sources, targets []VertexAlternative) (float64, []int64) {
	endpoints := [directionsCount][]VertexAlternative{sources, targets}
	var endpointsInternal [directionsCount][]vertexAlternativeInternal
	for d, alternatives := range endpoints {
		endpointsInternal[d] = qp.graph.vertexAlternativesToInternal(alternatives)
	}

	state := qp.acquireState()
	defer qp.releaseState(state)

	return qp.shortestPathWithAlternatives(state, endpointsInternal)
}

func (qp *QueryPool) shortestPathWithAlternatives(state *QueryState, endpoints [directionsCount][]vertexAlternativeInternal) (float64, []int64) {
	for d := forward; d < directionsCount; d++ {
		for _, endpoint := range endpoints[d] {
			if endpoint.vertexNum == vertexNotFound {
				continue
			}
			state.epochs[d][endpoint.vertexNum] = state.epoch
			state.dist[d][endpoint.vertexNum] = endpoint.additionalDistance
			heapEndpoint := &vertexDist{
				id:   endpoint.vertexNum,
				dist: endpoint.additionalDistance,
			}
			heap.Push(state.queues[d], heapEndpoint)
		}
	}
	return qp.shortestPathCore(state)
}

// ShortestPathOneToMany computes shortest paths from single source to multiple targets (thread-safe).
// This method can be safely called from multiple goroutines concurrently.
//
// source - user's defined ID of source vertex
// targets - set of user's defined IDs of target vertices
func (qp *QueryPool) ShortestPathOneToMany(source int64, targets []int64) ([]float64, [][]int64) {
	estimateAll := make([]float64, 0, len(targets))
	pathAll := make([][]int64, 0, len(targets))

	var ok bool
	if source, ok = qp.graph.mapping[source]; !ok {
		estimateAll = append(estimateAll, -1.0)
		pathAll = append(pathAll, nil)
		return estimateAll, pathAll
	}

	state := qp.acquireState()
	defer qp.releaseState(state)

	for _, target := range targets {
		// Increment epoch for each target query
		state.epoch++
		epoch := state.epoch

		if source == target {
			estimateAll = append(estimateAll, 0)
			pathAll = append(pathAll, []int64{source})
			continue
		}
		var ok bool
		if target, ok = qp.graph.mapping[target]; !ok {
			estimateAll = append(estimateAll, -1.0)
			pathAll = append(pathAll, nil)
			continue
		}

		state.epochs[forward][source] = epoch
		state.epochs[backward][target] = epoch
		state.dist[forward][source] = 0
		state.dist[backward][target] = 0

		// Reset prev maps for this query
		state.prev[forward] = make(map[int64]int64)
		state.prev[backward] = make(map[int64]int64)

		forwQ := &vertexDistHeap{}
		backwQ := &vertexDistHeap{}
		heap.Init(forwQ)
		heap.Init(backwQ)

		heapSource := &vertexDist{id: source, dist: 0}
		heapTarget := &vertexDist{id: target, dist: 0}
		heap.Push(forwQ, heapSource)
		heap.Push(backwQ, heapTarget)

		estimate, path := qp.shortestPathOneToManyCore(state, epoch, forwQ, backwQ)
		estimateAll = append(estimateAll, estimate)
		pathAll = append(pathAll, path)
	}

	return estimateAll, pathAll
}

func (qp *QueryPool) shortestPathOneToManyCore(state *QueryState, epoch int64, forwQ *vertexDistHeap, backwQ *vertexDistHeap) (float64, []int64) {
	estimate := Infinity
	var middleID int64

	for forwQ.Len() != 0 || backwQ.Len() != 0 {
		if forwQ.Len() != 0 {
			vertex1 := heap.Pop(forwQ).(*vertexDist)
			if vertex1.dist <= estimate {
				state.epochs[forward][vertex1.id] = epoch
				qp.relaxEdgesBiForward(state, vertex1, forwQ, epoch)
			}
			if state.epochs[backward][vertex1.id] == epoch {
				if vertex1.dist+state.dist[backward][vertex1.id] < estimate {
					middleID = vertex1.id
					estimate = vertex1.dist + state.dist[backward][vertex1.id]
				}
			}
		}

		if backwQ.Len() != 0 {
			vertex2 := heap.Pop(backwQ).(*vertexDist)
			if vertex2.dist <= estimate {
				state.epochs[backward][vertex2.id] = epoch
				qp.relaxEdgesBiBackward(state, vertex2, backwQ, epoch)
			}
			if state.epochs[forward][vertex2.id] == epoch {
				if vertex2.dist+state.dist[forward][vertex2.id] < estimate {
					middleID = vertex2.id
					estimate = vertex2.dist + state.dist[forward][vertex2.id]
				}
			}
		}
	}

	if estimate == Infinity {
		return -1, nil
	}
	return estimate, qp.graph.ComputePath(middleID, state.prev[forward], state.prev[backward])
}

func (qp *QueryPool) relaxEdgesBiForward(state *QueryState, vertex *vertexDist, forwQ *vertexDistHeap, epoch int64) {
	vertexList := qp.graph.Vertices[vertex.id].outIncidentEdges
	for i := range vertexList {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if qp.graph.Vertices[vertex.id].orderPos < qp.graph.Vertices[temp].orderPos {
			alt := state.dist[forward][vertex.id] + cost
			if state.epochs[forward][temp] != epoch || state.dist[forward][temp] > alt {
				state.dist[forward][temp] = alt
				state.prev[forward][temp] = vertex.id
				state.epochs[forward][temp] = epoch
				node := &vertexDist{id: temp, dist: alt}
				heap.Push(forwQ, node)
			}
		}
	}
}

func (qp *QueryPool) relaxEdgesBiBackward(state *QueryState, vertex *vertexDist, backwQ *vertexDistHeap, epoch int64) {
	vertexList := qp.graph.Vertices[vertex.id].inIncidentEdges
	for i := range vertexList {
		temp := vertexList[i].vertexID
		cost := vertexList[i].weight
		if qp.graph.Vertices[vertex.id].orderPos < qp.graph.Vertices[temp].orderPos {
			alt := state.dist[backward][vertex.id] + cost
			if state.epochs[backward][temp] != epoch || state.dist[backward][temp] > alt {
				state.dist[backward][temp] = alt
				state.prev[backward][temp] = vertex.id
				state.epochs[backward][temp] = epoch
				node := &vertexDist{id: temp, dist: alt}
				heap.Push(backwQ, node)
			}
		}
	}
}

// ShortestPathOneToManyWithAlternatives computes shortest paths with alternatives (thread-safe).
// This method can be safely called from multiple goroutines concurrently.
func (qp *QueryPool) ShortestPathOneToManyWithAlternatives(sourceAlternatives []VertexAlternative, targetsAlternatives [][]VertexAlternative) ([]float64, [][]int64) {
	estimateAll := make([]float64, 0, len(targetsAlternatives))
	pathAll := make([][]int64, 0, len(targetsAlternatives))

	sourceAlternativesInternal := qp.graph.vertexAlternativesToInternal(sourceAlternatives)

	state := qp.acquireState()
	defer qp.releaseState(state)

	for _, targetAlternatives := range targetsAlternatives {
		state.epoch++
		epoch := state.epoch

		targetAlternativesInternal := qp.graph.vertexAlternativesToInternal(targetAlternatives)

		// Reset prev maps for this query
		state.prev[forward] = make(map[int64]int64)
		state.prev[backward] = make(map[int64]int64)

		forwQ := &vertexDistHeap{}
		backwQ := &vertexDistHeap{}
		heap.Init(forwQ)
		heap.Init(backwQ)

		for _, sourceAlternative := range sourceAlternativesInternal {
			if sourceAlternative.vertexNum == vertexNotFound {
				continue
			}
			state.epochs[forward][sourceAlternative.vertexNum] = epoch
			state.dist[forward][sourceAlternative.vertexNum] = sourceAlternative.additionalDistance
			heapSource := &vertexDist{
				id:   sourceAlternative.vertexNum,
				dist: sourceAlternative.additionalDistance,
			}
			heap.Push(forwQ, heapSource)
		}

		for _, targetAlternative := range targetAlternativesInternal {
			if targetAlternative.vertexNum == vertexNotFound {
				continue
			}
			state.epochs[backward][targetAlternative.vertexNum] = epoch
			state.dist[backward][targetAlternative.vertexNum] = targetAlternative.additionalDistance
			heapTarget := &vertexDist{
				id:   targetAlternative.vertexNum,
				dist: targetAlternative.additionalDistance,
			}
			heap.Push(backwQ, heapTarget)
		}

		estimate, path := qp.shortestPathOneToManyCore(state, epoch, forwQ, backwQ)
		estimateAll = append(estimateAll, estimate)
		pathAll = append(pathAll, path)
	}

	return estimateAll, pathAll
}
