package benchmarks

import (
	"testing"
)

// check which way is faster - passing by parameter or operating on global value

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
