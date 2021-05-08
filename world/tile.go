package world

type TileType uint32

const (
	TileTypeNone TileType = iota
	TileTypeNormal
)

type Tile struct {
	SpriteID uint32
	TileType TileType
	Name     string
}
