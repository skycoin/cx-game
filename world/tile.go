package world

type TileType uint32

const (
	TileTypeNone TileType = iota
	TileTypeNormal
	TileTypeChild
)

func (tt TileType) ShouldRender() bool {
	return tt!=TileTypeNone
}

type Tile struct {
	SpriteID uint32
	TileType TileType
	Name     string
	OffsetX  int8
	OffsetY  int8
}

type MultiTile struct {
	Width     int
	Height    int
	TileType  TileType
	SpriteIDs []uint32
	Name      string
}
