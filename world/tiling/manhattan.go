package tiling

const (
	manhattanTilingWidth int = 4
	manhattanTilingHeight int = 4
)

type ManhattanTiling struct {}

func (t ManhattanTiling) Count() int {
	return manhattanTilingWidth * manhattanTilingHeight
}

func (t ManhattanTiling) Index(detailed DetailedNeighbours) int {
	n := detailed.Simplify()
	// block from (0,0) to (3,3)
	x := 1 + boolToInt(n.Left) - boolToInt(n.Right)
	y := 1 + boolToInt(n.Up) - boolToInt(n.Down)
	if !n.Left && !n.Right { x = 3 }
	if !n.Up && !n.Down { y = 3 }

	return manhattanTilingWidth*y + x
}
