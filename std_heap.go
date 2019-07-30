package ch

import "container/heap"

type minheapNode struct {
	id       int
	distance float64
}

type minheapSTD []minheapNode

func (h minheapSTD) Len() int           { return len(h) }
func (h minheapSTD) Less(i, j int) bool { return h[i].distance < h[j].distance } // Min-Heap
func (h minheapSTD) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minheapSTD) Push(x interface{}) {
	*h = append(*h, x.(minheapNode))
}

func (h *minheapSTD) Pop() interface{} {
	heapSize := len(*h)
	lastNode := (*h)[heapSize-1]
	*h = (*h)[0 : heapSize-1]
	return lastNode
}

func (h *minheapSTD) decrease_priority(id int, val float64) {
	for i := 0; i < len(*h); i++ {
		if (*h)[i].id == id {
			(*h)[i].distance = val
			break
		}
	}
}

func (h *minheapSTD) add_with_priority(id int, val float64) {
	nds := minheapNode{
		id:       id,
		distance: val,
	}
	heap.Push(h, nds)
}
