package benchmarks

import "testing"

var array_length = 10000

func BenchmarkAppendWithZeroCapacity(b *testing.B) {
	newSlice := make([]int, 0)
	for i := 0; i < b.N; i++ {
		for j := 0; j < array_length; j++ {
			newSlice = append(newSlice, j)
		}
	}
}

//second fastest
func BenchmarkAppendWithFullCapacity(b *testing.B) {
	newSlice := make([]int, 0, array_length)
	for i := 0; i < b.N; i++ {
		for j := 0; j < array_length; j++ {
			newSlice = append(newSlice, j)
		}
	}
}

//fastest,but requires to know the length of the slice beforehand
func BenchmarkSetValues(b *testing.B) {
	newSlice := make([]int, array_length)

	for i := 0; i < b.N; i++ {
		for j := 0; j < array_length; j++ {
			newSlice[j] = j
		}
	}
}
