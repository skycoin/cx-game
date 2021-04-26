package mainmap

//Map structure
type Map struct {
	bounds Fulstrum
	tiles  [size][size]*MapTile
}

//Contains neccessary data for drawing a single tile of the map
type MapTile struct {
	x                int
	y                int
	show             int
	tileIdBackground int
	tileIdMid        int
	tileIdFront      int
}

//FulStrum is to define drawing bounds on the map
type Fulstrum struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}
