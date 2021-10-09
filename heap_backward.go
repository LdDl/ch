package ch

// We can just wrap and override Less(...) function, but I'd prefer copy code for convenience
// type backwardHeap struct {
// 	forwardHeap
// }
// func (h backwardHeap) Less(i, j int) bool {
// 	return h.forwardHeap[i].revDistance < h.forwardHeap[j].revDistance
// }

type backwardHeap []*simpleNode

func (h backwardHeap) Len() int { return len(h) }
func (h backwardHeap) Less(i, j int) bool {
	return h[i].revDistance < h[j].revDistance
}
func (h backwardHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *backwardHeap) Push(x interface{}) { *h = append(*h, x.(*simpleNode)) }
func (h *backwardHeap) Pop() interface{} {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}
