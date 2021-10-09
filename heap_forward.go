package ch

type simpleNode struct {
	id          int64
	queryDist   float64
	revDistance float64
}

type forwardPropagationHeap []*simpleNode

func (h forwardPropagationHeap) Len() int            { return len(h) }
func (h forwardPropagationHeap) Less(i, j int) bool  { return h[i].queryDist < h[j].queryDist }
func (h forwardPropagationHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *forwardPropagationHeap) Push(x interface{}) { *h = append(*h, x.(*simpleNode)) }
func (h *forwardPropagationHeap) Pop() interface{} {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}
