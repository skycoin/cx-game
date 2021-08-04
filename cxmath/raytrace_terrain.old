package cxmath

import (
	"github.com/skycoin/cx-game/world/worldcollider"
)

type RayTraceData struct {
	Interval       float32
	CollisionPoint Vec2i
}

func (data *RayTraceData) SetCollisionPoint(vec Vec2) {
	data.CollisionPoint = vec
}

func raytrace_terrain(start, end Vec2, data *RayTraceData, collider worldcollider.WorldCollider) bool {
	data.Interval = 1

	if start.Equal(end) {
		if collider.TileIsSolid(int(start.X), int(start.Y)) {
			data.SetCollisionPoint(end)
			return true
		}
		return false
	}

	dx := Abs(end.X - start.X)
	dy := Abs(end.Y - start.Y)

	x := int(Floor(start.X))
	y := int(Floor(start.Y))

	dt_dx := 1.0 / dx
	dt_dy := 1.0 / dy

	t := float32(0)

	n := 1
	x_inc := 0
	y_inc := 0
	t_next_x := float32(0)
	t_next_y := float32(0)

	side := make([]int, 2)

	if dx == 0 {
		x_inc = 0
		t_next_x = dt_dx
	} else if end.X > start.X {
		x_inc = 1
		n += int(end.X) - x
	} else {
		x_inc = -1
		n += x - int(end.X)
		t_next_x = (start.X - Floor(start.X)) * dt_dx
	}

	if dy == 0 {
		y_inc = 0
		t_next_y = dt_dy //infinity
	} else if end.Y > start.Y {
		y_inc = 1
		n += int(end.Y) - y
	} else {
		y_inc = -1
		n += y - int(end.Y)
		t_next_y = (start.Y - Floor(start.Y)) * dt_dy
	}

	if t_next_x < t_next_y && t_next_x < t_next_y {
		side[0] = -x_inc
	}else if t_next_y < t_ne
}
