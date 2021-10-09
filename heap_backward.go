package ch

// We can just wrap and override Less(...) function, but I'd prefer copy code for convenience
// type backwardPropagationHeap struct {
// 	forwardHeap
// }
// func (h backwardPropagationHeap) Less(i, j int) bool {
// 	return h.forwardHeap[i].revDistance < h.forwardHeap[j].revDistance
// }

type backwardPropagationHeap []*simpleNode

func (h backwardPropagationHeap) Len() int { return len(h) }
func (h backwardPropagationHeap) Less(i, j int) bool {
	return h[i].revDistance < h[j].revDistance
}
func (h backwardPropagationHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *backwardPropagationHeap) Push(x interface{}) { *h = append(*h, x.(*simpleNode)) }
func (h *backwardPropagationHeap) Pop() interface{} {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}
