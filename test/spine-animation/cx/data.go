package cx

import (
	//	"fmt"
	"image"
	"image/color"
	"image/draw"

	//"image/draw"
	"math"
)

type PictureData struct {
	Pix    []color.RGBA
	Stride int
	Rect   Rect
}

// Rect is a 2D rectangle aligned with the axes of the coordinate system. It is defined by two
// points, Min and Max.
//
// The invariant should hold, that Max's components are greater or equal than Min's components
// respectively.

type Rect struct {
	Min, Max Vec
}

// Vec is a 2D vector type with X and Y coordinates.
type Vec struct {
	X, Y float64
}

func PictureDataFromImage(img image.Image) *PictureData {
	var rgba *image.RGBA
	if rgbaImg, ok := img.(*image.RGBA); ok {
		rgba = rgbaImg
	} else {
		rgba = image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	}

	verticalFlip(rgba)

	pd := MakePictureData(R(
		float64(rgba.Bounds().Min.X),
		float64(rgba.Bounds().Min.Y),
		float64(rgba.Bounds().Max.X),
		float64(rgba.Bounds().Max.Y),
	))

	for i := range pd.Pix {
		pd.Pix[i].R = rgba.Pix[i*4+0]
		pd.Pix[i].G = rgba.Pix[i*4+1]
		pd.Pix[i].B = rgba.Pix[i*4+2]
		pd.Pix[i].A = rgba.Pix[i*4+3]
	}

	return pd
}

// MakePictureData creates a zero-initialized PictureData covering the given rectangle.
func MakePictureData(rect Rect) *PictureData {
	w := int(math.Ceil(rect.Max.X)) - int(math.Floor(rect.Min.X))
	h := int(math.Ceil(rect.Max.Y)) - int(math.Floor(rect.Min.Y))
	pd := &PictureData{
		Stride: w,
		Rect:   rect,
	}
	pd.Pix = make([]color.RGBA, w*h)
	return pd
}

func verticalFlip(rgba *image.RGBA) {
	bounds := rgba.Bounds()
	width := bounds.Dx()

	tmpRow := make([]uint8, width*4)
	for i, j := 0, bounds.Dy()-1; i < j; i, j = i+1, j-1 {
		iRow := rgba.Pix[i*rgba.Stride : i*rgba.Stride+width*4]
		jRow := rgba.Pix[j*rgba.Stride : j*rgba.Stride+width*4]

		copy(tmpRow, iRow)
		copy(iRow, jRow)
		copy(jRow, tmpRow)
	}
}

// R returns a new Rect with given the Min and Max coordinates.
//
// Note that the returned rectangle is not automatically normalized.
func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{
		Min: Vec{minX, minY},
		Max: Vec{maxX, maxY},
	}
}
