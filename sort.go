package main

import (
	"container/heap"
	"runtime"
	"sort"
	"sync"
)

// MyMergeSort provides a multi-way merge sort algorithm
func MyMergeSort(arr []data) {
	shard := runtime.NumCPU()
	if len(arr) < shard {
		shard = len(arr)
	}

	resultChan := make(chan sliceHeader, shard)

	var wg sync.WaitGroup
	wg.Add(shard)

	size := len(arr) / shard

	for i := 0; i < shard; i++ {
		go func(idx int) {
			defer wg.Done()
			offset := idx * size

			dst := make([]data, size)
			copy(dst, arr[offset:offset+size])

			if idx == shard-1 {
				for j := offset + size; j < len(arr); j++ {
					dst = append(dst, arr[j])
				}
			}

			sort.Slice(dst, func(x, y int) bool {
				return dst[x].commit < dst[y].commit
			})
			resultChan <- sliceHeader{data: dst, idx: 0}
		}(i)
	}

	wg.Wait()
	close(resultChan)

	minHeap := &minHeap{}
	heap.Init(minHeap)

	for arr := range resultChan {
		heap.Push(minHeap, arr)
	}

	resultArr := arr[:0]

	for minHeap.Len() > 0 {
		min := minHeap.Peek()
		resultArr = append(resultArr, min.data[min.idx])

		min.idx++
		minHeap.increaseIdx()

		if min.idx <= len(min.data)-1 {
			heap.Fix(minHeap, 0)
		} else {
			heap.Pop(minHeap)
		}
	}
}

type sliceHeader struct {
	data []data
	idx  int
}

type minHeap []sliceHeader

func (h minHeap) Less(i, j int) bool {
	return h[i].data[h[i].idx].commit < h[j].data[h[j].idx].commit
}

func (h minHeap) Len() int {
	return len(h)
}

func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h minHeap) Peek() sliceHeader {
	return h[0]
}

func (h minHeap) increaseIdx() {
	h[0].idx++
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	re := old[n-1]
	*h = old[0 : n-1]
	return re
}

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(sliceHeader))
}
