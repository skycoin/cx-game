package benchmarks

import (
	"testing"
)

func BenchmarkByValueLarge(b *testing.B) {
	largeStruct := LargeStruct{Age: 15, Name: "John"}
	for i := 0; i < b.N; i++ {
		PassByValueLarge(largeStruct)
	}
}

func BenchmarkByPointerLarge(b *testing.B) {
	largeStruct := LargeStruct{Age: 15, Name: "John"}
	for i := 0; i < b.N; i++ {
		PassByPointerLarge(&largeStruct)
	}
}

func BenchmarkByInterfaceLarge(b *testing.B) {
	largeStruct := LargeStruct{Age: 15, Name: "John"}
	for i := 0; i < b.N; i++ {
		PassByInterfaceLarge(largeStruct)
	}
}
