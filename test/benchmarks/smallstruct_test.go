package benchmarks

import (
	"testing"
)

func BenchmarkByValueSmall(b *testing.B) {
	smallStruct := SmallStruct{Age: 15, Name: "John"}
	for i := 0; i < b.N; i++ {
		PassByValueSmall(smallStruct)
	}
}

func BenchmarkByPointerSmall(b *testing.B) {
	smallStruct := SmallStruct{Age: 15, Name: "John"}
	for i := 0; i < b.N; i++ {
		PassByPointerSmall(&smallStruct)
	}
}

func BenchmarkByInterfaceSmall(b *testing.B) {
	smallStruct := SmallStruct{Age: 15, Name: "John"}
	for i := 0; i < b.N; i++ {
		PassByInterfaceSmall(smallStruct)
	}
}
