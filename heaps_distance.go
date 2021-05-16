package ch

/* Copied from container/heap - https://golang.org/pkg/container/heap/ */
// Why make copy? Just want to avoid type conversion

type distanceHeap []*Vertex

func (h distanceHeap) Len() int           { return len(h) }
func (h distanceHeap) Less(i, j int) bool { return h[i].distance.distance < h[j].distance.distance }
func (h distanceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

/* Actual code of interface шьздуьутефешщт */
// type Interface interface {
// 	Len() int <-- taken from sort.Interface
// 	Less(i, j int) bool <-- taken from sort.Interface
// 	Swap(i, j int) <-- taken from sort.Interface
// 	Push(x *Vertex) // add x as element Len()
// 	Pop() *Vertex   // remove and return element Len() - 1.
// }

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *distanceHeap) Push(x *Vertex) {
	*h = append(*h, x)
	h.up(h.Len() - 1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *distanceHeap) Pop() *Vertex {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *distanceHeap) Remove(i int) *Vertex {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.Pop()
}

func (h *distanceHeap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func (h *distanceHeap) down(i0, n int) bool {
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
