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
