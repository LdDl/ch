package ch

type potatoHeapParallel []*Vertex

func (h potatoHeapParallel) Len() int { return len(h) }
func (h potatoHeapParallel) Less(i, j, threadID int) bool {
	return h[i].distance_v2[threadID].distance < h[j].distance_v2[threadID].distance
}
func (h potatoHeapParallel) Swap(i, j int)   { h[i], h[j] = h[j], h[i] }
func (h *potatoHeapParallel) Push(x *Vertex) { *h = append(*h, x) }
func (h *potatoHeapParallel) Pop() *Vertex {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}

// The PotatoInterfaceParallel type describes the requirements
// for a type using the routines in this package.
// Any type that implements it may be used as a
// min-heap with the following invariants (established after
// Init has been called or if the data is empty or sorted):
//
//	!h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
//
// Note that Push and Pop in this interface are for package heap's
// implementation to call. To add and remove things from the heap,
// use heap.Push and heap.Pop.
type PotatoInterfaceParallel interface {
	Len() int
	Swap(i, j int)
	Less(i, j, threadID int) bool
	Push(x *Vertex) // add x as element Len()
	Pop() *Vertex   // remove and return element Len() - 1.
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func InitParallel(h PotatoInterfaceParallel, threads int) {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		for threadID := 0; threadID < threads; threadID++ {
			downParallel(h, i, n, threadID)
		}
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func PushParallel(h PotatoInterfaceParallel, x *Vertex, threadID int) {
	h.Push(x)
	upParallel(h, h.Len()-1, threadID)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func PopParallel(h PotatoInterfaceParallel, threadID int) *Vertex {
	n := h.Len() - 1
	h.Swap(0, n)
	downParallel(h, 0, n, threadID)
	return h.Pop()
}

func upParallel(h PotatoInterfaceParallel, j, threadID int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i, threadID) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func downParallel(h PotatoInterfaceParallel, i0, n, threadID int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1, threadID) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i, threadID) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}
