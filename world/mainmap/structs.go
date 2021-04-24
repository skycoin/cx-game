package mainmap

type Map struct {
	bounds Fullstrum
	tiles  []*MapTile
}

type MapTile struct {
	spriteId         int
	tileIdBackground int
	tileIdMid        int
	tileIdFront      int
	x                int
	y                int
	show             int
}

type Fullstrum struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}
