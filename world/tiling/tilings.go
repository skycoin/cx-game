package tiling

type TilingID int
const (
	FullTilingID TilingID = iota
	ManhattanTilingID
	PlatformTilingID

	NUM_TILINGS // DO NOT SET MANUALLY
)

var tilings [NUM_TILINGS]Tiling
func init() {
	tilings[FullTilingID] = FullTiling{}
	tilings[ManhattanTilingID] = ManhattanTiling{}
	tilings[PlatformTilingID] = PlatformTiling{}
}

func ApplyTiling(id TilingID, neighbours DetailedNeighbours) int {
	return tilings[id].Index(neighbours)
}

func ByName(name string) (TilingID,bool) {
	if name == "full" { return FullTilingID,true }
	if name == "manhattan" { return ManhattanTilingID,true }
	if name == "platform" { return PlatformTilingID,true }
	return -1,false
}

func ( id TilingID ) Get() Tiling {
	return tilings[id]
}
