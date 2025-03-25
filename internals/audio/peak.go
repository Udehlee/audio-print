package audio

import (
	"container/heap"
)

// Peak represents a frequency peak.
type Peak struct {
	FreqIndx int
	Value    float64
}

type PeakHeap []Peak

func (h PeakHeap) Len() int            { return len(h) }
func (h PeakHeap) Less(i, j int) bool  { return h[i].Value < h[j].Value } // Min-Heap
func (h PeakHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PeakHeap) Push(x interface{}) { *h = append(*h, x.(Peak)) }
func (h *PeakHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// ExtractPeaks selects the top spectral frequency peaks from FFT magnitudes using a Min-Heap
func (a AudioService) ExtractPeaks(magnitudes []float64, numPeaks int) []int {
	h := &PeakHeap{}
	heap.Init(h)

	for i := 1; i < len(magnitudes)-1; i++ {
		if magnitudes[i] > magnitudes[i-1] && magnitudes[i] > magnitudes[i+1] {
			peak := Peak{FreqIndx: i, Value: magnitudes[i]}

			if h.Len() < numPeaks {
				heap.Push(h, peak)
			} else if peak.Value > (*h)[0].Value {
				heap.Pop(h)
				heap.Push(h, peak)
			}
		}
	}

	selectedPeaks := make([]int, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		selectedPeaks[i] = heap.Pop(h).(Peak).FreqIndx
	}

	return selectedPeaks
}

// HashPeaks combines two frequencies and time difference into a single 64-bit hash
func (a AudioService) HashPeaks(f1, f2, deltaTime int) uint64 {
	return uint64(f1)<<32 | uint64(f2)<<16 | uint64(deltaTime)
}
