package cxmath

import (
	"math"
	"github.com/skycoin/cx-game/physics"
)

// http://playtechs.blogspot.com/2007/03/raytracing-on-grid.html
type Point struct {X,Y int}
type GridLine struct {
	increment physics.Vec2i
	n int
	next float64
	dt float64
}

func setupGridLine(x int,x0,x1,dx,dt_dx float64, axis physics.Vec2i) GridLine {
	if dx==0 {
		return GridLine {
			increment: physics.Vec2i{},
			next: dt_dx,
			n: 0,
			dt: dt_dx,
		}
	}
	if x1 > x0 {
		return GridLine {
			increment: axis,
			next: (math.Floor(x0) + 1 -x0) * dt_dx,
			n: int(math.Floor(x1)) - x,
			dt:dt_dx,
		}
	} else {
		return GridLine {
			increment: axis.Mult(-1),
			n: x-int(math.Floor(x1)),
			next: (x0 - math.Floor(x0)) * dt_dx,
			dt: dt_dx,
		}
	}

}

func getCloserGridLine(xLines, yLines *GridLine) *GridLine {
	if xLines.next < yLines.next {
		return xLines
	} else {
		return yLines
	}
}

func Raytrace(x0,y0, x1,y1 float64) []physics.Vec2i {
	dx := math.Abs(x1-x0)
	dy := math.Abs(y1-y0)

	x := int(math.Floor(x0))
	y := int(math.Floor(y0))

	dt_dx := 1.0 / dx
	dt_dy := 1.0 / dy


	xLines := setupGridLine(x,x0,x1,dx,dt_dx,physics.Vec2i{1,0})
	yLines := setupGridLine(y,y0,y1,dy,dt_dy,physics.Vec2i{0,1})

	n := 1 + xLines.n + yLines.n

	pos := physics.Vec2i{int32(x0),int32(y0)}
	points := make([]physics.Vec2i,n)
	for i:=0; i<n; i++ {
		points[i] = pos
		closerLine := getCloserGridLine(&xLines,&yLines)
		pos = pos.Add(closerLine.increment)
		closerLine.next += closerLine.dt
	}
	return points
}
