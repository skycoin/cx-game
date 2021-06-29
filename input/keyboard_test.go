package input

import "testing"

func BenchmarkKeyBoardDeleteMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		myMap := make(map[int]bool)
		for k := 0; k < 1000; k++ {
			myMap[k] = true
		}
		for k := range myMap {
			delete(myMap, k)
		}
	}
}

func BenchmarkKeyBoardEmptyMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		myMap := make(map[int]bool)
		for k := 0; k < 1000; k++ {
			myMap[k] = true
		}
		for k := range myMap {
			myMap[k] = false
		}
	}
}
