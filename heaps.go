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
