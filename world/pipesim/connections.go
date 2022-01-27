package pipesim

import (
	"github.com/skycoin/cx-game/world/tiling"
	"github.com/skycoin/cx-game/cxmath"
)

type Connections struct { Up, Left, Right, Down bool }

func (c Connections) Bits() [4]bool {
	return [4]bool { c.Up, c.Left, c.Right, c.Down }
}

func (c Connections) Diff(newC Connections) ConnectionDiff {
	return ConnectionDiff {
		Up: computeBoolDiff(c.Up, newC.Up),
		Left: computeBoolDiff(c.Left, newC.Left),
		Right: computeBoolDiff(c.Right, newC.Right),
		Down: computeBoolDiff(c.Down, newC.Down),
	}
}

func (c Connections) ApplyDiff(diff ConnectionDiff) Connections {
	return Connections {
		Up: applyBoolDiff(c.Up, diff.Up),
		Left: applyBoolDiff(c.Left, diff.Left),
		Right: applyBoolDiff(c.Right, diff.Right),
		Down: applyBoolDiff(c.Down, diff.Down),
	}
}

func (c Connections) Valid(possible Connections) bool {
	isValid :=
		( possible.Up || !c.Up  ) &&
		( possible.Down || !c.Down  ) &&
		( possible.Left || !c.Left  ) &&
		( possible.Right || !c.Right  )
	return isValid
}

func ConnectionsFromNeighbours(n tiling.DetailedNeighbours) Connections {
	s := n.Simplify()
	return Connections { Up: s.Up, Left: s.Left, Right: s.Right, Down: s.Down }
}

func ConnectedNeighbours(
		connections Connections, neighbours tiling.DetailedNeighbours,
) tiling.DetailedNeighbours {
	connectedNeighbours := neighbours // copy
	// hide neighbours which we shouldn't see due to connections
	if !connections.Up { connectedNeighbours.Up = tiling.None }
	if !connections.Down { connectedNeighbours.Down = tiling.None }
	if !connections.Left { connectedNeighbours.Left = tiling.None }
	if !connections.Right { connectedNeighbours.Right = tiling.None }
	return connectedNeighbours
}

func composeBits(bits []bool) int {
	place := 1
	sum := 0
	for _,bit := range bits {
		if bit { sum += place }
		place *= 2
	}
	return sum
}

func decomposeBits(composed int, bits []bool) {
	place := 1
	for idx := range bits {
		bits[idx] = composed&place != 0
		place *= 2
	}
}

// given some current connection state, cycles to another connection state.
// loops over all possible states eventually
func (c Connections) Next(valid Connections) Connections {
	bits := []bool { c.Up, c.Left, c.Right, c.Down }
	i := composeBits(bits)
	for true {
		decomposed := [4]bool{}
		decomposeBits(i+1,decomposed[:])
		d := decomposed
		candidate := Connections { d[0], d[1], d[2], d[3] }
		if candidate.Valid(valid) { return candidate }
		i = (i+1) % 16
	}
	return Connections{} // unreachable anyway
}

func FindNewConnections(disp cxmath.Vec2i) (Connections,Connections) {
	if disp.X == 0 && disp.Y == 0 { return Connections{},Connections{} }
	if disp.X == 0 {
		if disp.Y > 0 {
			return Connections { Up: true}, Connections { Down: true }
		} else {
			return Connections { Down: true}, Connections { Up: true }
		}
	} else {
		if disp.X > 0 {
			return Connections { Right: true}, Connections { Left: true }
		} else {
			return Connections { Left: true}, Connections { Right: true }
		}
	}
	return Connections{},Connections{}
}

func (x Connections) OR(y Connections) Connections {
	return Connections {
		Up: x.Up || y.Up,
		Down: x.Down || y.Down,
		Left: x.Left || y.Left,
		Right: x.Right || y.Right,
	}
}

func (x Connections) AND(y Connections) Connections {
	return Connections {
		Up: x.Up && y.Up,
		Down: x.Down && y.Down,
		Left: x.Left && y.Left,
		Right: x.Right && y.Right,
	}
}

func (x Connections) NOT() Connections {
	return Connections {
		Up: !x.Up, Down: !x.Down, Left: !x.Left, Right: !x.Right,
	}
}
