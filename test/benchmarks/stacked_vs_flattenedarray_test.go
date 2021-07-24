package benchmarks

import "testing"

func BenchmarkFlattenedArray(b *testing.B) {
	var array = []float32{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}
	for i := 0; i < b.N; i++ {
		for j := range array {
			array[j] += 1
		}
	}
}

//not much difference
func BenchmarkStackedArray(b *testing.B) {
	var array = [][]float32{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10},
	}
	for i := 0; i < b.N; i++ {
		for j := range array {
			for k := range array[j] {
				array[j][k] += 1
			}
		}
	}
}
