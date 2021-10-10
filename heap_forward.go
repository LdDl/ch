package ch

type bidirectionalVertex struct {
	id               int64
	queryDist        float64
	revQueryDistance float64
}

type forwardHeap []*bidirectionalVertex

func (h forwardHeap) Len() int            { return len(h) }
func (h forwardHeap) Less(i, j int) bool  { return h[i].queryDist < h[j].queryDist }
func (h forwardHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *forwardHeap) Push(x interface{}) { *h = append(*h, x.(*bidirectionalVertex)) }
func (h *forwardHeap) Pop() interface{} {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}
