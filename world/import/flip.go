package worldimport

import (
	"github.com/lafriks/go-tiled"
)

func scaleFromFlipFlags(tile *tiled.LayerTile) (int,int) {
	x := 1
	y := 1
	if tile.HorizontalFlip {
		x *= -1
	}
	if tile.VerticalFlip {
		y *= -1
	}
	if tile.DiagonalFlip {
		x *= -1
		y *= -1
	}

	return x,y
}
