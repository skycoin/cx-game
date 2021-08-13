package cxmath

import (
	"testing"
)

func TestPackRectangles(T *testing.T) {
	width := 5
	packed := PackRectangles(width,[]Vec2i {
		Vec2i{1,1}, Vec2i{2,1},          Vec2i{2,2},
		Vec2i{3,1},
		Vec2i{5,1},
	})

	for idx,rect := range packed {
		if rect.Right() > int32(width) {
			T.Errorf("rect #%d [%+v] is too far right",idx,rect)
		}
		if rect.Bottom() > 3 {
			T.Errorf("rect #%d[%+v] is not packed tightly enough",idx,rect)
		}
		for otherIdx:=idx+1; otherIdx<len(packed); otherIdx++ {
			otherRect := packed[otherIdx]
			if rect.Intersects(otherRect) {
				T.Errorf(
					"rect #%d [%+v] intersects with rect #%d [%+v]",
					idx,rect, otherIdx,otherRect,
				)
			}
		}
	}
}

func containsVec2i( list []Vec2i, target Vec2i ) bool {
	for _,candidate := range list {
		if candidate.Eq(target) { return true }
	}
	return false
}

func TestRectNeighbours(T *testing.T) {
	rect := Rect { Vec2i {2,2}, Vec2i {2,1} }
	got := rect.Neighbours()
	want := []Vec2i {
		Vec2i { 1,1 }, Vec2i { 2,1 }, Vec2i { 3,1 }, Vec2i { 4,1 },
		Vec2i { 1,3 }, Vec2i { 2,3 }, Vec2i { 3,3 }, Vec2i { 4,3 },
		Vec2i { 1,2 },                               Vec2i { 4,2 },
	}
	if len(got) != len(want) {
		T.Errorf("wrong number of neighbours")
	}
	for _,pos := range got {
		if !containsVec2i( want, pos ) {
			T.Errorf(
				"Cannot find neighbour %v in neighbours for rect with\n"+
				"origin=%v and size=%v",
				pos,rect.Origin,rect.Size,
			)
		}
	}
}
