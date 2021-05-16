package ch

type importanceHeap []*Vertex

func (h importanceHeap) Len() int            { return len(h) }
func (h importanceHeap) Less(i, j int) bool  { return h[i].importance < h[j].importance }
func (h importanceHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *importanceHeap) Push(x interface{}) { *h = append(*h, x.(*Vertex)) }
func (h *importanceHeap) Pop() interface{} {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}
func (h importanceHeap) Peek() *Vertex {
	lastNode := h[0]
	return lastNode
}

type distanceHeap struct {
	importanceHeap
}

func (h distanceHeap) Less(i, j int) bool {
	return h.importanceHeap[i].distance.distance < h.importanceHeap[j].distance.distance
}

type distanceHeapParallel []*Vertex

func (h distanceHeapParallel) Len() int { return len(h) }
func (h distanceHeapParallel) Less(i, j, threadID int) bool {
	return h[i].distance_v2[threadID].distance < h[j].distance_v2[threadID].distance
}
func (h distanceHeapParallel) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *distanceHeapParallel) Pop() (v *Vertex) {
	old := *h
	v = old[0]
	*h = old[1:]
	return
}
func (h *distanceHeapParallel) Push(v *Vertex, threadID int) {
	var changed bool
	old := *h
	updated := *h
	for i, w := range old {
		if v.vertexNum == w.vertexNum {
			if changed {
				if i+1 < len(updated) {
					updated = append(updated[:i], updated[i+1:]...)
				} else {
					updated = updated[:i]
				}
			} else if v.distance_v2[threadID].distance < w.distance_v2[threadID].distance {
				updated[i] = v
			}
			changed = true
		} else if v.distance_v2[threadID].distance < w.distance_v2[threadID].distance {
			changed = true
			updated = append(old[:i], v)
			updated = append(updated, w)
			updated = append(updated, old[i+1:]...)
		}
	}
	if !changed {
		updated = append(old, v)
	}
	*h = updated
}

type distanceHeapExplicit []*Vertex

func (h distanceHeapExplicit) Len() int      { return len(h) }
func (h distanceHeapExplicit) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h distanceHeapExplicit) Less(i, j int) bool {
	return h[i].distance.distance < h[j].distance.distance
}
func (pq *distanceHeapExplicit) Pop() *Vertex {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (h *distanceHeapExplicit) Push(v *Vertex) {
	var changed bool
	old := *h
	updated := *h
	for i, w := range old {
		if v.vertexNum == w.vertexNum {
			if changed {
				if i+1 < len(updated) {
					updated = append(updated[:i], updated[i+1:]...)
				} else {
					updated = updated[:i]
				}
			} else if v.distance.distance < w.distance.distance {
				updated[i] = v
			}
			changed = true
		} else if v.distance.distance < w.distance.distance {
			changed = true
			updated = append(old[:i], v)
			updated = append(updated, w)
			updated = append(updated, old[i+1:]...)
		}
	}
	if !changed {
		updated = append(old, v)
	}
	*h = updated
}
