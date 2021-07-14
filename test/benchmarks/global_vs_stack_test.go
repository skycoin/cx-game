package benchmarks

import (
	"testing"
)

func BenchmarkPassByParameter(b *testing.B) {
	myUser := User{
		age:  15,
		name: "John",
	}

	for i := 0; i < b.N; i++ {
		PassByParameter(myUser)
	}
}

func BenchmarkPassByGlobal(b *testing.B) {
	user = User{
		age:  15,
		name: "John",
	}
	for i := 0; i < b.N; i++ {
		PassByGlobal()
	}
}
