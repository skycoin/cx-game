package tiling

const (
	fullTilingWidth int = 11
	fullTilingHeight int = 5
)

type FullTiling struct {}

func (t FullTiling) Count() int { return fullTilingWidth * fullTilingHeight }

func (t FullTiling) Index(detailed DetailedNeighbours) int {
	n := detailed.Simplify()
	innerCorners := n.countInnerCorners()
	// default to the solid square
	x := 1
	y := 1
	if innerCorners == 0 {
		// block from (0,0) to (3,3)
		x = 1 + boolToInt(n.Left) - boolToInt(n.Right)
		y = 1 + boolToInt(n.Up) - boolToInt(n.Down)
		if !n.Left && !n.Right { x = 3 }
		if !n.Up && !n.Down { y = 3 }
	}
	if innerCorners == 1 {
		// block from (4,0) to (7,3)
		if !n.Left { x = 4 }
		if !n.Right { x = 7 }
		if n.Left && n.Right {
			if n.hasRightInnerCorner() { x = 5 }
			if n.hasLeftInnerCorner() { x = 6 }
		}

		if !n.Up { y = 0 }
		if !n.Down { y = 3 }
		if n.Up && n.Down {
			if n.hasDownInnerCorner() { y = 1 }
			if n.hasUpInnerCorner() { y = 2 }
		}
	}
	if innerCorners == 2 {
		// horizontal strip from (4,4) to (7,4)
		if n.upRightInnerCorner() && n.downRightInnerCorner() {
			y = 4
			if !n.Left { x = 4 } else { x = 5 }
		}
		if n.upLeftInnerCorner() && n.downLeftInnerCorner() {
			y = 4
			if !n.Right { x = 7 } else { x = 6 }
		}

		// vertical strip from (8,0) to (8,3)
		if n.downLeftInnerCorner() && n.downRightInnerCorner() {
			x = 8
			if !n.Up { y = 0 } else { y = 1 }
		}
		if n.upLeftInnerCorner() && n.upRightInnerCorner() {
			x = 8
			if !n.Down { y = 3 } else { y = 2 }
		}

		// vertical strip from (9,0) to (9,1)
		if n.upRightInnerCorner() && n.downLeftInnerCorner() {
			x = 9; y = 0;
		}
		if n.upLeftInnerCorner() && n.downRightInnerCorner() {
			x = 9; y = 1;
		}
	}
	if innerCorners == 3 {
		// block from (9,2) to (10,3)
		if !n.downRightInnerCorner() { x = 9; y = 2 }
		if !n.downLeftInnerCorner() { x = 10; y = 2 }
		if !n.upRightInnerCorner() { x = 9; y = 3 }
		if !n.upLeftInnerCorner() { x= 10; y = 3 }
	}
	if innerCorners == 4 {
		// tile at (8,4)
		x = 8
		y = 4
	}
	return fullTilingWidth*y + x
}
