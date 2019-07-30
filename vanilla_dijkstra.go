package ch

import (
	"container/heap"
	"log"
	"math"
)

// VanillaShortestPath Computes and returns shortest path and it's cost (vanilla Dijkstra's algorithm)
//
// If there are some errors then function returns '-1.0' as cost and nil as shortest path
//
// source User's definied ID of source vertex
// target User's definied ID of target vertex
//
// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
func (graph *Graph) VanillaShortestPath(source, target int) (float64, []int) {

	if source == target {
		return 0, []int{source}
	}
	ok := true

	if source, ok = graph.mapping[source]; !ok {
		log.Println("No such source")
		return -1.0, nil
	}
	if target, ok = graph.mapping[target]; !ok {
		log.Println("No such target")
		return -1.0, nil
	}

	// create vertex set Q
	Q := &minheapSTD{}

	// dist[source] ← 0
	distance := make(map[int]float64, len(graph.Vertices))
	distance[source] = 0

	// prev[v] ← UNDEFINED
	prev := make(map[int]int, len(graph.Vertices))
	// st := time.Now()
	// for each vertex v in Graph:
	for i := range graph.Vertices {
		// if v ≠ source:
		if graph.Vertices[i].vertexNum != source {
			// dist[v] = INFINITY
			distance[graph.Vertices[i].vertexNum] = math.MaxFloat64
		}
		// prev[v] ← UNDEFINED
		// nothing here
	}
	Q.add_with_priority(graph.Vertices[source].vertexNum, distance[graph.Vertices[source].vertexNum])
	// lograph.Println("Simple: elapsed for init:", time.Since(st))
	heap.Init(Q)
	prevNodeID := -1
	// while Q is not empty:
	for Q.Len() != 0 {
		// fmt.Println("Iteration")
		// u ← Q.extract_min()
		u := heap.Pop(Q).(minheapNode)
		// fmt.Println("\t Current vertex", u.id)
		restrictions := make(map[int]int)
		ok := false
		viaRestrictionID := -1
		if restrictions, ok = graph.restrictions[prevNodeID]; ok {
			// fmt.Println("restrictions met", restrictions, u.id)
			if viaRestrictionID, ok = restrictions[u.id]; ok {
				// fmt.Println("restriction id", viaRestrictionID)
			}
		}

		// popcounter++
		// if u == target:
		if u.id == target {
			// break
			// log.Println("====================>>>break")
			break
		}

		vertexList := graph.Vertices[u.id].outEdges
		costList := graph.Vertices[u.id].outECost

		// fmt.Println("\t vertex list", vertexList)
		// fmt.Println("\t vertex costs", costList)

		// for each neighbor v of u:

		for v := range vertexList {
			neighbor := vertexList[v]
			if neighbor == viaRestrictionID {
				// distance[neighbor] = math.MaxFloat64
				distance[u.id] = math.MaxFloat64
				// log.Println("con", prevNodeID, neighbor, u.id, distance[u.id])
				continue
			}
			cost := costList[v]
			// alt ← dist[u] + length(u, v)
			alt := distance[u.id] + cost
			// fmt.Println("\t\tneightbor", v, neighbor)
			// fmt.Println("\t\t\tcost to v", cost)
			// fmt.Println("\t\t\tdistance to u", distance[u.id])
			// fmt.Println("\t\t\tdistance to neighbor", neighbor, distance[neighbor])
			// fmt.Println("\t\t\talt", alt)
			// if alt < dist[v]
			if distance[neighbor] > alt {
				// dist[v] ← alt
				distance[neighbor] = alt
				// prev[v] ← u
				prev[neighbor] = u.id
				// Q.decrease_priority(v, alt)
				// Q.decrease_priority(v, alt)
				// fmt.Println("\t\t\tAdded!", neighbor, alt)
				Q.add_with_priority(neighbor, alt)
			}
		}

		prevNodeID = u.id
		// heap.Init(Q)
	}
	// lograph.Println("Simple: elapsed for run:", time.Since(st))
	// lograph.Println("iterations to get path", iters, popcounter)

	// path = []
	var path []int
	// u = target
	u := target

	// while prev[u] is defined:
	for {
		if _, ok := prev[u]; !ok {
			break
		}
		// path.push_front(u)
		temp := make([]int, len(path)+1)
		temp[0] = u
		copy(temp[1:], path)
		path = temp

		// u = prev[u]
		u = prev[u]
	}

	temp := make([]int, len(path)+1)
	temp[0] = source
	copy(temp[1:], path)
	path = temp

	usersLabelsPath := make([]int, len(path))
	for e := 0; e < len(usersLabelsPath); e++ {
		usersLabelsPath[e] = graph.Vertices[path[e]].Label //append(path, graph.Vertices[e.Value.(int)].Label)
	}

	// return path, prev
	return distance[target], usersLabelsPath
}
