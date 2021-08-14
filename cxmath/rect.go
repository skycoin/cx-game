package cxmath

import (
	"sort"
)

type Rect struct {
	Origin Vec2i
	Size Vec2i
}


func (r Rect) Area() int32 {
	return r.Size.X * r.Size.Y
}

func (r Rect) Right() int32 {
	return r.Origin.X + r.Size.X
}

func (r Rect) Bottom() int32 {
	return r.Origin.Y + r.Size.Y
}

func (r Rect) Top() int32 { return r.Origin.Y }
func (r Rect) Left() int32 { return r.Origin.X }

func (r Rect) Contains(x,y int32) bool {
	return x >= r.Left() && x < r.Right() &&
			y >= r.Top() && y < r.Bottom()
}

func (this Rect) Intersects(other Rect) bool {
	verticalIntersect :=
		!( this.Top() >= other.Bottom() || this.Bottom() <= other.Top() )
	horizontalIntersect :=
		!( this.Left() >= other.Right() || this.Right() <= other.Left() )
	return verticalIntersect && horizontalIntersect
}

func (r Rect) Cells() []Vec2i {
	cells := make([]Vec2i,r.Area())
	idx := 0
	for y:=int32(0); y<r.Size.Y; y++ {
		for x:=int32(0); x<r.Size.X; x++ {
			cells[idx] = Vec2i { X: r.Origin.X + x, Y: r.Origin.Y + y }
			idx++
		}
	}
	return cells
}

func (r Rect) Neighbours() []Vec2i {
	neighbours := make([]Vec2i, 0, r.Size.X*2 + r.Size.Y*2 - 4)
	// above/below neighbours ( including corners )
	for x := int32(-1) ; x <= r.Size.X ; x++ {
		neighbours = append( neighbours,
			Vec2i { r.Origin.X + x, r.Origin.Y-1 },
			Vec2i { r.Origin.X + x, r.Origin.Y+r.Size.Y },
		)
	}
	// left/right neighbours ( excluding corners )
	for y := int32(0) ; y < r.Size.Y ; y++ {
		neighbours = append( neighbours,
			Vec2i { r.Origin.X-1, r.Origin.Y+y },
			Vec2i { r.Origin.X+r.Size.X, r.Origin.Y+y },
		)
	}
	return neighbours
}

type BinaryGrid struct {
	Width,Height int
	Occupied []bool
}

func NewBinaryGrid(width,height int) BinaryGrid {
	return BinaryGrid {
		Width: width, Height: height,
		Occupied: make([]bool,width*height),
	}
}

func (g *BinaryGrid) At(x,y int) *bool {
	return &g.Occupied[y*g.Width+x]
}

func (g *BinaryGrid) MarkRect(r Rect) {
	for _,cell := range r.Cells() { g.MarkPoint( int(cell.X), int(cell.Y) ) }
}

func (g *BinaryGrid) MarkPoint(x,y int) {
	*g.At(x,y) = true
}

// x,y are top left
func (g *BinaryGrid) RectFits(rect Rect) bool {
	for _,cell := range rect.Cells() {
		if *g.At(int(cell.X),int(cell.Y)) { return false }
	}
	return true
}

func (g *BinaryGrid) PlaceRect(rect *Rect) {
	// move rect until it fits
	for !g.RectFits(*rect) {
		if int(rect.Right()) < g.Width {
			rect.Origin.X++
		} else  {
			rect.Origin.X=0; rect.Origin.Y++
		}
	}
	g.MarkRect(*rect)
}

// for packed height worst case (maximum),
// assume that all rectangles are stacked vertically
func maxHeight(sizes []Vec2i) int {
	height := 0
	for _,size := range sizes { height += int(size.Y) }
	return height
}

// pack boxes into a fixed width, minimizing height
func PackRectangles(width int, sizes []Vec2i) []Rect {
	grid := NewBinaryGrid(width, maxHeight(sizes))
	// assemble rectangles
	rects := make([]Rect,len(sizes))
	for idx,size := range sizes { rects[idx] = Rect { Size: size } }

	// sort pointers to rects by descending area
	rectsByArea := make([]*Rect,len(rects))
	for idx,_ := range rects { rectsByArea[idx] = &rects[idx] }
	sort.Slice(rectsByArea, func(i,j int) bool {
		return rectsByArea[i].Area() > rectsByArea[j].Area()
	})

	// place rectangles in grid
	for _,rect := range rectsByArea { grid.PlaceRect(rect) }

	return rects
}
