package pipesim

// a BoolDiff describes the difference between two boolean states.
// one can apply this diff to force another boolean on or off,
// or leave it unchanged.
type BoolDiff int
const (
	NO_CHANGE BoolDiff = iota
	SET_ON
	SET_OFF
)


// compute the difference between two boolean values.
// typically, these would be the old and new states for some flag.
func computeBoolDiff(x,y bool) BoolDiff {
	if y && !x { return SET_ON }
	if !y && x { return SET_OFF }
	return NO_CHANGE
}

// apply a boolean diff.
// this either forces x ON or OFF,
// or leaves it unchanged
func applyBoolDiff(x bool, diff BoolDiff) bool {
	if diff == SET_OFF { return false }
	if diff == SET_ON { return true }
	return x
}

// a ConnectionDiff describes the differene between two connections.
// this is used to update neighbouring pipes on place/connect
// by mirroring connection changes
type ConnectionDiff struct { Up,Left,Right,Down BoolDiff }

// a PipeNeighbour describes a pipe, along with any connections/disconnections
// that must be made.
type PipeNeighbour struct {
	X,Y int // position
	ConnectionDiff ConnectionDiff
}

// returns a list of pipes which need updating,
// including connections that need to be modified.
func PipeNeighbours(
		x,y int, connections, oldConnections Connections,
) []PipeNeighbour {
	diff := oldConnections.Diff(connections)
	return []PipeNeighbour {
		PipeNeighbour {
			x-1, y,
			ConnectionDiff { Right: diff.Left },
		},
		PipeNeighbour {
			x+1, y,
			ConnectionDiff { Left: diff.Right },
		},
		PipeNeighbour {
			x, y-1,
			ConnectionDiff { Up: diff.Down },
		},
		PipeNeighbour {
			x, y+1,
			ConnectionDiff { Down: diff.Up },
		},
	}
}
