package ch

import "container/heap"

type minHeapVertex struct {
	id       int64
	distance float64
}

type minHeap []*minHeapVertex

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].distance < h[j].distance } // Min-Heap
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(*minHeapVertex))
}

func (h *minHeap) Pop() interface{} {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}

func (h *minHeap) add_with_priority(id int64, val float64) {
	nds := &minHeapVertex{
		id:       id,
		distance: val,
	}
	heap.Push(h, nds)
}
