package tiling



type Tiling interface {
	// how many tiles does this tiling have?
	Count() int
	// what index of the loaded sprites should the tiling use,
	// based on the given neighbours?
	Index(DetailedNeighbours) int
}
