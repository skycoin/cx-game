package mainmap

type Map struct {
	bounds Fulstrum
	tiles  [size][size]*MapTile
}

type MapTile struct {
	x                int
	y                int
	show             int
	tileIdBackground int
	tileIdMid        int
	tileIdFront      int
}

type Fulstrum struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}
