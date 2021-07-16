package physics

import "testing"

//benchmark the access to ground collision

func BenchmarkDirectAccess(b *testing.B) {
	newBody := Body{}
	var counter int
	for i := 0; i < b.N; i++ {
		if newBody.Collisions.Below {
			counter += 1
		}
	}
}

func BenchmarkGetterMethod(b *testing.B) {
	newBody := Body{}
	var counter int
	for i := 0; i < b.N; i++ {
		if newBody.IsOnGround() {
			counter += 1
		}
	}
}
