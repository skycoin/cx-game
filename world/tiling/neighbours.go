package tiling

type Neighbour int
const (
	None Neighbour = iota
	Self
	Solid
)

type DetailedNeighbours struct {
	Left, Right, Up, Down Neighbour
	UpLeft, UpRight, DownLeft, DownRight Neighbour
}

func (d DetailedNeighbours) Simplify() Neighbours {
	return Neighbours {
		Left: d.Left == Self,
		Right: d.Right == Self,
		Up: d.Up == Self,
		Down: d.Down == Self,
		UpLeft: d.UpLeft == Self,
		UpRight: d.UpRight == Self,
		DownLeft: d.DownLeft == Self,
		DownRight: d.DownRight == Self,
	}
}

// https://www.tilesetter.org/docs/generating_tilesets
type Neighbours struct {
	Left, Right, Up, Down bool
	UpLeft, UpRight, DownLeft, DownRight bool
}

func NewSolidNeighbours() Neighbours {
	return Neighbours {
		Left: true, Right: true, Up: true, Down: true,
		UpLeft: true, UpRight: true, DownLeft: true, DownRight: true,
	}
}

func (n Neighbours) upLeftInnerCorner() bool {
	return n.Up && n.Left && !n.UpLeft
}

func (n Neighbours) upRightInnerCorner() bool {
	return n.Up && n.Right && !n.UpRight
}

func (n Neighbours) downLeftInnerCorner() bool {
	return n.Down && n.Left && !n.DownLeft
}

func (n Neighbours) downRightInnerCorner() bool {
	return n.Down && n.Right && !n.DownRight
}

func (n Neighbours) hasLeftInnerCorner() bool {
	return n.upLeftInnerCorner() || n.downLeftInnerCorner()
}

func (n Neighbours) hasRightInnerCorner() bool {
	return n.upRightInnerCorner() || n.downRightInnerCorner()
}

func (n Neighbours) hasUpInnerCorner() bool {
	return n.upLeftInnerCorner() || n.upRightInnerCorner()
}

func (n Neighbours) hasDownInnerCorner() bool {
	return n.downLeftInnerCorner() || n.downRightInnerCorner()
}

func (n Neighbours) countInnerCorners() int {
	return boolToInt(n.upLeftInnerCorner()) +
		boolToInt(n.upRightInnerCorner()) +
		boolToInt(n.downLeftInnerCorner()) +
		boolToInt(n.downRightInnerCorner())
}

func (n Neighbours) countEdges() int {
	return boolToInt(n.Left) + boolToInt(n.Up) +
		boolToInt(n.Right) + boolToInt(n.Down)
}

func boolToInt(x bool) int {
	if x { return 1 } else { return 0 }
}

// compute coords on the maps TileSetter provides in their "Blob" format
// TODO maybe alter internal packing for more efficient determining
func (n Neighbours) blobCoords() (x,y int) {
	innerCorners := n.countInnerCorners()
	// default to the solid square
	x = 1
	y = 1
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
		if n.downRightInnerCorner() { x = 9; y = 2 }
		if n.downLeftInnerCorner() { x = 10; y = 2 }
		if n.upRightInnerCorner() { x = 9; y = 3 }
		if n.upLeftInnerCorner() { x= 10; y = 3 }
	}
	if innerCorners == 4 {
		// tile at (8,4)
		x = 8
		y = 4
	}
	return x,y
}

func (n Neighbours) simpleBlobCoords() (x,y int) {
	// block from (0,0) to (3,3)
	x = 1 + boolToInt(n.Left) - boolToInt(n.Right)
	y = 1 + boolToInt(n.Up) - boolToInt(n.Down)
	if !n.Left && !n.Right { x = 3 }
	if !n.Up && !n.Down { y = 3 }
	return x,y
}
