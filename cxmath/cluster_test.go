package cxmath

import "testing"

func TestClusters(t *testing.T) {
	points := []Vec2i {
		Vec2i { 3,3 }, Vec2i {1001, 1001}, Vec2i {4,4 }, Vec2i {1000,1000},
	}
	clusters := FindClusters(points,2)
	t.Logf("%v", clusters)
}
