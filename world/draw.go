package world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/mathi"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render"
)

type PositionedTile struct {
	Tile Tile
	Position cxmath.Vec2i
}

func (planet *Planet) DrawLayer(layer Layer, cam *camera.Camera) {
	//configureGlForPlanet()
	planet.program.Use()
	defer planet.program.StopUsing()

	w := int(planet.Width)
	// split up planet into 2 hemispheres to achieve wrap around
	// without calculating relative tile positions individually
	planet.DrawHemisphere(layer, cam, 0, w/2)
	planet.DrawHemisphere(layer, cam, w/2, w)
}

func configureGlForPlanet() {
	gl.BindVertexArray(render.QuadVao)
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
}

func (planet *Planet) DrawHemisphere(
	layer Layer, cam *camera.Camera, left,right int,
) {
	center := float32( (left+right)/2 )

	camToCenter := planet.ShortestDisplacement(
		mgl32.Vec2{ cam.X, cam.Y },
		mgl32.Vec2 { center, 0 } ) // to.y doesn't matter here

	projection := cam.GetProjectionMatrix()
	planet.program.
		SetMat4("projection",&projection )

	visible := planet.visibleTiles(layer, cam, left, right)
	bins := planet.binTilesBySpritesheet(visible)

	for tex,tiles := range bins {
		planet.drawSpritesheetBin(tex,tiles,camToCenter.X(),cam.Y, center)
	}
}

func (planet *Planet) Draw(cam *camera.Camera, layerID LayerID) {
	planet.DrawLayer(planet.Layers[layerID], cam)
}

func (planet *Planet) visibleTiles(
	layer Layer, cam *camera.Camera, left,right int,
) []PositionedTile {
	bottom := mathi.Max(cam.Frustum.Bottom,0)
	top := mathi.Min(cam.Frustum.Top, int(planet.Height))
	capacity := (top-bottom+1) * (right-left+1)
	positionedTiles := make([]PositionedTile,0,capacity)

	for y := bottom ; y <= top ; y++ {
		for x := left ; x <= right; x ++ {
			tileIdx := planet.GetTileIndex(x,y)
			tile := layer.Tiles[tileIdx]
			if tile.TileCategory != TileCategoryNone {
				positionedTiles = append(positionedTiles, PositionedTile {
					Position: cxmath.Vec2i { X: int32(x), Y: int32(y) },
					Tile: tile,
				})
			}
		}
	}
	return positionedTiles
}

func (planet *Planet) binTilesBySpritesheet(
		positionedTiles []PositionedTile,
) map[uint32][]PositionedTile {
	bins := make(map[uint32][]PositionedTile)
	for _,positionedTile := range positionedTiles {
		meta := spriteloader.GetSpriteMetadata(positionedTile.Tile.SpriteID)
		_,ok := bins[meta.GpuTex]
		if !ok {
			bins[meta.GpuTex] = []PositionedTile{}
		}
		bins[meta.GpuTex] = append(bins[meta.GpuTex], positionedTile)
	}
	return bins
}

func (planet *Planet) drawSpritesheetBin(
		tex uint32, tiles []PositionedTile,
		// should probably structure these further
		camToCenterX float32, camY float32, center float32,
) {
	var instance int32 = 0
	worlds := [100]mgl32.Mat4{}
	texScales := [100]mgl32.Vec2{}
	texOffsets := [100]mgl32.Vec2{}
	gl.BindTexture(gl.TEXTURE_2D, tex)
	planet.program.SetUint("ourTexture",tex)

	for _,positionedTile := range tiles {
		tile := positionedTile.Tile
		x := positionedTile.Position.X
		y := positionedTile.Position.Y
		meta := spriteloader.GetSpriteMetadata(tile.SpriteID)

		worlds[instance] = mgl32.Translate3D(
			camToCenterX+float32(x)-center, float32(y)-camY, 0,
		)
		texScales[instance] = mgl32.Vec2 { meta.ScaleX, meta.ScaleY }
		// hack - fix translate*scale vs scale*translate mismatch
		// between shader and spriteloader
		texOffsets[instance] = mgl32.Vec2 {
			float32(meta.PosX) / float32(meta.ScaleX),
			float32(meta.PosY) / float32(meta.ScaleY),
		}

		instance++
		if instance==100 {
			planet.program.SetMat4s("worlds",worlds[:])
			planet.program.SetVec2s("texScales",texScales[:])
			planet.program.SetVec2s("texOffsets",texOffsets[:])
			gl.DrawArraysInstanced(gl.TRIANGLES,0,6, instance)
			instance=0
		}
	}
	// draw leftovers
	planet.program.SetMat4s("worlds",worlds[:])
	planet.program.SetVec2s("texScales",texScales[:])
	planet.program.SetVec2s("texOffsets",texOffsets[:])
	gl.DrawArraysInstanced(gl.TRIANGLES,0,6, instance)
}
