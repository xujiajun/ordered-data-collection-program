package main

import (
	"sort"
	"testing"
)

func BenchmarkMyMergeSort(b *testing.B) {
	numElements := 16 << 20
	src := make([]data, numElements)
	original := make([]data, numElements)
	prepare(original)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(src, original)
		b.StartTimer()
		MyMergeSort(src)
	}
}

func BenchmarkNormalSort(b *testing.B) {
	numElements := 16 << 20
	src := make([]data, numElements)
	original := make([]data, numElements)
	prepare(original)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(src, original)
		b.StartTimer()
		sort.Slice(src, func(i, j int) bool { return src[i].commit < src[j].commit })
	}
}
