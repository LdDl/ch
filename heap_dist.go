package ch

type vertexDist struct {
	id   int64
	dist float64
}

type vertexDistHeap []*vertexDist

func (h vertexDistHeap) Len() int            { return len(h) }
func (h vertexDistHeap) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h vertexDistHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *vertexDistHeap) Push(x interface{}) { *h = append(*h, x.(*vertexDist)) }
func (h *vertexDistHeap) Pop() interface{} {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}
