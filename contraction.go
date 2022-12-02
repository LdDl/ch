package ch

import (
	"container/heap"
	"fmt"
	"time"
)

var (
	tmLayout = "2006-01-2T15:04:05.999999999"
)

// Preprocess Computes contraction hierarchies and returns node ordering
func (graph *Graph) Preprocess(pqImportance *importanceHeap) {
	extractionOrder := int64(0)
	for pqImportance.Len() != 0 {
		// Lazy update heuristic:
		// update Importance of vertex "on demand" as follows:
		// Before contracting vertex with currently smallest Importance, recompute its Importance and see if it is still the smallest
		// If not pick next smallest one, recompute its Importance and see if that is the smallest now; If not, continue in same way ...
		vertex := heap.Pop(pqImportance).(*Vertex)
		vertex.computeImportance()
		if pqImportance.Len() != 0 && vertex.importance > pqImportance.Peek().importance {
			pqImportance.Push(vertex)
			continue
		}
		vertex.orderPos = extractionOrder
		graph.contractNode(vertex)
		if graph.verbose {
			if extractionOrder > 0 && pqImportance.Len()%1000 == 0 {
				fmt.Printf("Contraction Order: %d / %d, Remain vertices in heap: %d. Currect shortcuts num: %d Initial edges num: %d Time: %v\n", extractionOrder, len(graph.Vertices), pqImportance.Len(), graph.shortcutsNum, graph.edgesNum, time.Now().Format(tmLayout))
			}
		}
		extractionOrder++
	}
}

// markNeighbors
//
// inEdges Incoming edges from vertex
// outEdges Outcoming edges from vertex
func (graph *Graph) markNeighbors(inEdges, outEdges []incidentEdge) {
	for i := range inEdges {
		temp := inEdges[i]
		graph.Vertices[temp.vertexID].delNeighbors++
	}
	for i := range outEdges {
		temp := outEdges[i]
		graph.Vertices[temp.vertexID].delNeighbors++
	}
}

// contractNode
//
// vertex Vertex to be contracted
func (graph *Graph) contractNode(vertex *Vertex) {
	// Consider all vertices with edges incoming TO current vertex as U
	incomingEdges := vertex.inIncidentEdges

	// Consider all vertices with edges incoming FROM current vertex as W
	outcomingEdges := vertex.outIncidentEdges

	// Exclude vertex for local shortest paths searches
	vertex.contracted = true
	// Tell neighbor vertices that current vertex has been contracted
	graph.markNeighbors(incomingEdges, outcomingEdges)

	// For every vertex 'w' in W, compute Pw as the cost from 'u' to 'w' through current vertex, which is the sum of the edge weights w(u, vertex) + w(vertex, w).
	inMax := 0.0
	outMax := 0.0
	for i := range incomingEdges {
		if graph.Vertices[incomingEdges[i].vertexID].contracted {
			continue
		}
		if inMax < incomingEdges[i].weight {
			inMax = incomingEdges[i].weight
		}
	}
	for i := range outcomingEdges {
		if graph.Vertices[outcomingEdges[i].vertexID].contracted {
			continue
		}
		if outMax < outcomingEdges[i].weight {
			outMax = outcomingEdges[i].weight
		}
	}
	// Then Pmax is the maximum pMax over all 'w' in W.
	pmax := inMax + outMax

	// Perform a standard Dijkstra’s shortest path search from 'u' on the subgraph excluding current vertex.
	graph.processIncidentEdges(vertex, pmax)
}

// processIncidentEdges Returns evaluated shorcuts
//
// vertex - Vertex for making possible shortcuts around
// pmax - path cost restriction
func (graph *Graph) processIncidentEdges(vertex *Vertex, pmax float64) {
	incomingEdges := vertex.inIncidentEdges
	outcomingEdges := vertex.outIncidentEdges
	if len(outcomingEdges) == 0 {
		return
	}

	batchShortcuts := make([]ShortcutPath, 0)

	previousOrderPos := int64(vertex.orderPos - 1)
	for _, u := range incomingEdges {
		inVertex := u.vertexID
		// Do not consider any vertex has been excluded earlier
		if graph.Vertices[inVertex].contracted {
			continue
		}
		inCost := u.weight
		graph.shortestPathsWithMaxCost(inVertex, pmax, previousOrderPos) // Finds the shortest distances from the inVertex to all outVertices.
		for _, w := range outcomingEdges {
			outVertex := w.vertexID
			outVertexPtr := graph.Vertices[outVertex]
			// Do not consider any vertex has been excluded earlier
			if outVertexPtr.contracted {
				continue
			}
			outCost := w.weight
			neighborsWeights := inCost + outCost
			// For each w, if dist(u, w) > Pw we add a shortcut edge uw with weight Pw.
			// If this condition doesn’t hold, no shortcut is added.
			if outVertexPtr.distance.distance > neighborsWeights ||
				outVertexPtr.distance.previousOrderPos != previousOrderPos || // Optional condition: if previous shortestPathsWithMaxCost(...) call has changed shortest path tree
				outVertexPtr.distance.previousSourceID != inVertex { // Optional condition: if previous shortestPathsWithMaxCost(...) call has changed shortest path tree

				// Collect needed shortcuts
				batchShortcuts = append(batchShortcuts, ShortcutPath{From: inVertex, To: outVertex, Via: vertex.vertexNum, Cost: neighborsWeights})
			}
		}
	}

	graph.insertShortcuts(batchShortcuts)
}

// insertShortcuts Creates (or updates: it depends on conditions) multiple shortcuts in graph structure
func (graph *Graph) insertShortcuts(batchShortcuts []ShortcutPath) {
	for i := range batchShortcuts {
		d := batchShortcuts[i]
		graph.createOrUpdateShortcut(d.From, d.To, d.Via, d.Cost)
	}
}

// createOrUpdateShortcut Creates (or updates: it depends on conditions) shortcut
//
// fromVertex - Library defined ID of source vertex where shortcut starts from
// fromVertex - Library defined ID of target vertex where shortcut leads to
// viaVertex - Library defined ID of vertex through which the shortcut exists
// summaryCost - Travel path of a shortcut
func (graph *Graph) createOrUpdateShortcut(fromVertex, toVertex, viaVertex int64, summaryCost float64) {
	if _, ok := graph.shortcuts[fromVertex]; !ok {
		// If there is no such shortcut then add one.
		graph.shortcuts[fromVertex] = make(map[int64]*ShortcutPath)
	}
	if existing, ok := graph.shortcuts[fromVertex][toVertex]; !ok {
		// Prepare shorcut pointer if there is no From-To-Via combo
		graph.shortcuts[fromVertex][toVertex] = &ShortcutPath{
			From: fromVertex,
			To:   toVertex,
			Via:  viaVertex,
			Cost: summaryCost,
		}
		graph.Vertices[fromVertex].addOutIncidentEdge(toVertex, summaryCost)
		graph.Vertices[toVertex].addInIncidentEdge(fromVertex, summaryCost)
		graph.shortcutsNum++
	} else {
		// If shortcut already exists
		if summaryCost < existing.Cost {
			// If middle vertex is not optimal for shortcut then change cost
			existing.Cost = summaryCost
			updatedOutSuccess := graph.Vertices[fromVertex].updateOutIncidentEdge(toVertex, summaryCost)
			if !updatedOutSuccess {
				panic(fmt.Sprintf("Should not happen [1]. Can't update outcoming incident edge. %d has no common edge with %d", fromVertex, toVertex))
			}
			updatedInSuccess := graph.Vertices[toVertex].updateInIncidentEdge(fromVertex, summaryCost)
			if !updatedInSuccess {
				panic(fmt.Sprintf("Should not happen [2]. Can't update incoming incident edge. %d has no common edge with %d", toVertex, fromVertex))
			}
			// We should check if the middle vertex is still the same
			// We could just do existing.ViaVertex = viaVertex, but it could be helpful for debugging purposes.
			if existing.Via != viaVertex {
				existing.Via = viaVertex
			}
		}
	}
}
