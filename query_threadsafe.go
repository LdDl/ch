package ch

import (
	"container/heap"
	"sync"
)

type queryHeap []vertexDist

func (h queryHeap) Len() int {
	return len(h)
}

// Use the same strict comparison and swaps as container/heap, including equal priorities.
func (h *queryHeap) push(value vertexDist) {
	*h = append(*h, value)
	j := len(*h) - 1
	for {
		i := (j - 1) / 2
		if i == j || !((*h)[j].dist < (*h)[i].dist) {
			break
		}
		(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
		j = i
	}
}

func (h *queryHeap) pop() vertexDist {
	n := len(*h) - 1
	(*h)[0], (*h)[n] = (*h)[n], (*h)[0]
	i := 0
	for {
		left := 2*i + 1
		if left >= n || left < 0 {
			break
		}
		j := left
		if right := left + 1; right < n && (*h)[right].dist < (*h)[left].dist {
			j = right
		}
		if !((*h)[j].dist < (*h)[i].dist) {
			break
		}
		(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
		i = j
	}
	value := (*h)[n]
	*h = (*h)[:n]
	return value
}

type queryArc struct {
	from int64
	to   int64
}

// Scratch belongs to the acquired query state; the returned path owns its storage.
func (qp *QueryPool) computeQueryPath(state *QueryState, middle int64) []int64 {
	chain := state.pathChain[:0]
	chain = append(chain, middle)
	for at := middle; ; {
		before := state.previous[forward][at]
		if before < 0 {
			break
		}
		chain = append(chain, before)
		at = before
	}
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	for at := middle; ; {
		after := state.previous[backward][at]
		if after < 0 {
			break
		}
		chain = append(chain, after)
		at = after
	}
	state.pathChain = chain
	path := state.pathOutput[:0]
	path = append(path, qp.graph.Vertices[chain[0]].Label)
	stack := state.pathStack[:0]
	for i := 1; i < len(chain); i++ {
		stack = append(stack, queryArc{chain[i-1], chain[i]})
		for len(stack) != 0 {
			arc := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if shortcut, ok := qp.graph.shortcuts[arc.from][arc.to]; ok {
				stack = append(stack, queryArc{shortcut.Via, arc.to}, queryArc{arc.from, shortcut.Via})
			} else {
				path = append(path, qp.graph.Vertices[arc.to].Label)
			}
		}
	}
	state.pathStack = stack
	state.pathOutput = path
	return append([]int64(nil), path...)
}

// QueryState holds all the buffers needed for a single shortest path query.
// This is used by the thread-safe query methods to avoid sharing state between goroutines.
type QueryState struct {
	// Current query epoch (incremented each query within this state)
	epoch int64
	// Distance arrays for bidirectional search
	dist [directionsCount][]float64
	// Epoch markers (if != epoch, distance is Infinity)
	epochs [directionsCount][]int64
	// Predecessor maps for one-to-many queries
	prev [directionsCount]map[int64]int64
	// Priority queues for bidirectional search
	queues         [directionsCount]queryHeap
	pathChain      []int64
	pathStack      []queryArc
	pathOutput     []int64
	stopAtEstimate bool
	// Dense predecessors are valid for the current distance epoch.
	previous [directionsCount][]int64

	// ManyToMany query state buffers
	// Outer slice indexed by endpoint, inner slice indexed by vertex
	manyToManyDist   [directionsCount][][]float64
	manyToManyEpochs [directionsCount][][]int64
	manyToManyPrev   [directionsCount][]map[int64]int64
	manyToManyEpoch  int64
}

// QueryPool provides thread-safe access to pooled QueryState objects.
// Use this when you need to call shortest path queries from multiple goroutines.
type QueryPool struct {
	pool        sync.Pool
	graph       *Graph
	nonnegative bool
}

// NewQueryPool creates a new QueryPool for concurrent query execution.
// The pool lazily initializes QueryState objects as needed.
// Keep the prepared graph unchanged while using the pool; recreate the pool after graph updates.
func (graph *Graph) NewQueryPool() *QueryPool {
	nonnegative := true
	for i := range graph.Vertices {
		for _, edge := range graph.Vertices[i].outIncidentEdges {
			if !(edge.weight >= 0) {
				nonnegative = false
			}
		}
	}
	return &QueryPool{
		graph:       graph,
		nonnegative: nonnegative,
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
			state.previous[d] = make([]int64, n)
			state.epochs[d] = make([]int64, n)
			state.prev[d] = make(map[int64]int64)
			state.queues[d] = nil
		}
	}

	// Increment epoch to invalidate previous distances
	state.epoch++
	state.stopAtEstimate = qp.nonnegative

	// Clear prev maps
	for d := forward; d < directionsCount; d++ {
		for vertex := range state.prev[d] {
			delete(state.prev[d], vertex)
		}
		state.queues[d] = state.queues[d][:0]
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
		state.previous[d][endpoints[d]] = -1
		heapEndpoint := vertexDist{
			id:   endpoints[d],
			dist: 0,
		}
		state.queues[d].push(heapEndpoint)
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
			// A strictly worse minimum cannot relax or improve a meeting on nonnegative paths.
			if state.stopAtEstimate && state.queues[d][0].dist > estimate {
				state.queues[d] = state.queues[d][:0]
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
	return estimate, qp.computeQueryPath(state, middleID)
}

func (qp *QueryPool) directionalSearch(state *QueryState, d direction, reverseDirection direction, estimate *float64, middleID *int64) {
	vertex := state.queues[d].pop()
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
					state.previous[d][temp] = vertex.id
					node := vertexDist{
						id:   temp,
						dist: alt,
					}
					state.queues[d].push(node)
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
			if !(endpoint.additionalDistance >= 0) {
				state.stopAtEstimate = false
			}
			if endpoint.vertexNum == vertexNotFound {
				continue
			}
			state.epochs[d][endpoint.vertexNum] = state.epoch
			state.dist[d][endpoint.vertexNum] = endpoint.additionalDistance
			state.previous[d][endpoint.vertexNum] = -1
			heapEndpoint := vertexDist{
				id:   endpoint.vertexNum,
				dist: endpoint.additionalDistance,
			}
			state.queues[d].push(heapEndpoint)
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

// ShortestPathManyToMany computes shortest paths between multiple sources and targets (thread-safe).
// This method can be safely called from multiple goroutines concurrently.
//
// sources - set of user's defined IDs of source vertices
// targets - set of user's defined IDs of target vertices
func (qp *QueryPool) ShortestPathManyToMany(sources, targets []int64) ([][]float64, [][][]int64) {
	// Copy input slices to avoid modifying caller's data
	sourcesCopy := make([]int64, len(sources))
	targetsCopy := make([]int64, len(targets))
	copy(sourcesCopy, sources)
	copy(targetsCopy, targets)

	endpoints := [directionsCount][]int64{sourcesCopy, targetsCopy}
	for d, directionEndpoints := range endpoints {
		for i, endpoint := range directionEndpoints {
			var ok bool
			if endpoints[d][i], ok = qp.graph.mapping[endpoint]; !ok {
				endpoints[d][i] = -1
			}
		}
	}

	state := qp.acquireState()
	defer qp.releaseState(state)

	return qp.shortestPathManyToMany(state, endpoints)
}

// initManyToManyBuffers ensures buffers are allocated and properly sized for the query
func (qp *QueryPool) initManyToManyBuffers(state *QueryState, numSources, numTargets int) {
	n := len(qp.graph.Vertices)
	endpointCounts := [directionsCount]int{numSources, numTargets}

	for d := forward; d < directionsCount; d++ {
		// Grow outer slices if needed
		if len(state.manyToManyDist[d]) < endpointCounts[d] {
			newDist := make([][]float64, endpointCounts[d])
			newEpochs := make([][]int64, endpointCounts[d])
			newPrev := make([]map[int64]int64, endpointCounts[d])
			copy(newDist, state.manyToManyDist[d])
			copy(newEpochs, state.manyToManyEpochs[d])
			copy(newPrev, state.manyToManyPrev[d])
			state.manyToManyDist[d] = newDist
			state.manyToManyEpochs[d] = newEpochs
			state.manyToManyPrev[d] = newPrev
		}

		// Ensure each endpoint has properly sized inner slices
		for i := 0; i < endpointCounts[d]; i++ {
			if len(state.manyToManyDist[d][i]) < n {
				state.manyToManyDist[d][i] = make([]float64, n)
				state.manyToManyEpochs[d][i] = make([]int64, n)
			}
			if state.manyToManyPrev[d][i] == nil {
				state.manyToManyPrev[d][i] = make(map[int64]int64)
			}
		}
	}
}

// getManyToManyDist returns distance for endpoint at given vertex, using epoch for lazy clearing
func (qp *QueryPool) getManyToManyDist(state *QueryState, d direction, endpointIdx int, vertexID int64) float64 {
	if state.manyToManyEpochs[d][endpointIdx][vertexID] != state.manyToManyEpoch {
		return Infinity
	}
	return state.manyToManyDist[d][endpointIdx][vertexID]
}

// setManyToManyDist sets distance for endpoint at given vertex
func (qp *QueryPool) setManyToManyDist(state *QueryState, d direction, endpointIdx int, vertexID int64, dist float64) {
	state.manyToManyDist[d][endpointIdx][vertexID] = dist
	state.manyToManyEpochs[d][endpointIdx][vertexID] = state.manyToManyEpoch
}

func (qp *QueryPool) shortestPathManyToMany(state *QueryState, endpoints [directionsCount][]int64) ([][]float64, [][][]int64) {
	numSources := len(endpoints[forward])
	numTargets := len(endpoints[backward])

	// Increment epoch for lazy buffer clearing
	state.manyToManyEpoch++
	qp.initManyToManyBuffers(state, numSources, numTargets)

	// Clear prev maps
	for d := forward; d < directionsCount; d++ {
		endpointCount := numSources
		if d == backward {
			endpointCount = numTargets
		}
		for i := 0; i < endpointCount; i++ {
			for k := range state.manyToManyPrev[d][i] {
				delete(state.manyToManyPrev[d][i], k)
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
			qp.setManyToManyDist(state, d, endpointIdx, endpoint, 0)
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
				qp.directionalSearchManyToMany(state, d, endpointIdx, queues, estimates, middleIDs, numSources, numTargets)
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
			paths[sourceIdx][targetIdx] = qp.graph.ComputePath(
				middleIDs[sourceIdx][targetIdx],
				state.manyToManyPrev[forward][sourceIdx],
				state.manyToManyPrev[backward][targetIdx],
			)
		}
	}

	return estimates, paths
}

func (qp *QueryPool) directionalSearchManyToMany(state *QueryState, d direction, endpointIdx int, queues [directionsCount][]*vertexDistHeap, estimates [][]float64, middleIDs [][]int64, numSources, numTargets int) {
	q := queues[d][endpointIdx]
	vertex := heap.Pop(q).(*vertexDist)

	// Skip if we've already found a better path
	currentDist := qp.getManyToManyDist(state, d, endpointIdx, vertex.id)
	if vertex.dist > currentDist {
		return
	}

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

		// Only explore upward in CH
		if qp.graph.Vertices[vertex.id].orderPos < qp.graph.Vertices[temp].orderPos {
			alt := vertex.dist + cost
			tempDist := qp.getManyToManyDist(state, d, endpointIdx, temp)
			if alt < tempDist {
				qp.setManyToManyDist(state, d, endpointIdx, temp, alt)
				state.manyToManyPrev[d][endpointIdx][temp] = vertex.id
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
		revDist := qp.getManyToManyDist(state, 1-d, revIdx, vertex.id)
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

// ShortestPathManyToManyWithAlternatives computes shortest paths with alternatives (thread-safe).
// This method can be safely called from multiple goroutines concurrently.
//
// sourcesAlternatives - set of user's defined IDs of source vertices with additional penalty
// targetsAlternatives - set of user's defined IDs of target vertices with additional penalty
func (qp *QueryPool) ShortestPathManyToManyWithAlternatives(sourcesAlternatives, targetsAlternatives [][]VertexAlternative) ([][]float64, [][][]int64) {
	endpoints := [directionsCount][][]VertexAlternative{sourcesAlternatives, targetsAlternatives}
	var endpointsInternal [directionsCount][][]vertexAlternativeInternal
	for d, directionEndpoints := range endpoints {
		endpointsInternal[d] = make([][]vertexAlternativeInternal, 0, len(directionEndpoints))
		for _, alternatives := range directionEndpoints {
			endpointsInternal[d] = append(endpointsInternal[d], qp.graph.vertexAlternativesToInternal(alternatives))
		}
	}

	state := qp.acquireState()
	defer qp.releaseState(state)

	return qp.shortestPathManyToManyWithAlternatives(state, endpointsInternal)
}

func (qp *QueryPool) shortestPathManyToManyWithAlternatives(state *QueryState, endpoints [directionsCount][][]vertexAlternativeInternal) ([][]float64, [][][]int64) {
	numSources := len(endpoints[forward])
	numTargets := len(endpoints[backward])

	// Increment epoch for lazy buffer clearing
	state.manyToManyEpoch++
	qp.initManyToManyBuffers(state, numSources, numTargets)

	// Clear prev maps
	for d := forward; d < directionsCount; d++ {
		endpointCount := numSources
		if d == backward {
			endpointCount = numTargets
		}
		for i := 0; i < endpointCount; i++ {
			for k := range state.manyToManyPrev[d][i] {
				delete(state.manyToManyPrev[d][i], k)
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
				currentDist := qp.getManyToManyDist(state, d, endpointIdx, alt.vertexNum)
				if alt.additionalDistance < currentDist {
					qp.setManyToManyDist(state, d, endpointIdx, alt.vertexNum, alt.additionalDistance)
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
				qp.directionalSearchManyToMany(state, d, endpointIdx, queues, estimates, middleIDs, numSources, numTargets)
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
			paths[sourceIdx][targetIdx] = qp.graph.ComputePath(
				middleIDs[sourceIdx][targetIdx],
				state.manyToManyPrev[forward][sourceIdx],
				state.manyToManyPrev[backward][targetIdx],
			)
		}
	}

	return estimates, paths
}
