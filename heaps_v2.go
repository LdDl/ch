package ch

type potatoHeap []*Vertex

func (h potatoHeap) Len() int           { return len(h) }
func (h potatoHeap) Less(i, j int) bool { return h[i].distance.distance < h[j].distance.distance }
func (h potatoHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *potatoHeap) Push(x *Vertex) {
	*h = append(*h, x)
	h.up(h.Len() - 1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *potatoHeap) Pop() *Vertex {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}

func (h *potatoHeap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func (h *potatoHeap) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}
