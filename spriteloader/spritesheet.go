package spriteloader

import (
	"image/png"
	"os"
)

type SpriteSheet struct {
	rows    int
	columns int
}


//NewSpriteSheet creates new SpriteSheet instance
func NewSpriteSheet(filepath string) (*SpriteSheet, error) {
	rows, cols, err := getDimensions(filepath)
	if err != nil {
		return nil, err
	}

	return &SpriteSheet{
		rows:    rows,
		columns: cols,
	}, nil
}

func getDimensions(spriteFile string) (int, int, error) {
	file, err := os.Open(spriteFile)
	if err != nil {
		return 0, 0, err
	}
	img, err := png.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Bounds().Dx() / 32, img.Bounds().Dy() / 32, nil
}