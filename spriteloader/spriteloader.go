package spriteloader

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
)

type SPRITE_LOADER struct {
	SSheet               Sprite_Sheet
	SpriteID             int
	SpriteNameToSpriteID map[string]int
}

func (s *SPRITE_LOADER) Assign_Sheet_ID() {
	s.SSheet.Sheet_ID = 1
}

//For current example sprite.png 256x256. The per sprite in sheet will be 64x64
// func main() {
// 	var s SPRITE_LOADER
// 	s.INIT()
// 	s.Load_Sprite_File("./sprite.png")
// 	sprite, err := s.GetSprite(1, "SPRITE_1 X 2", 1, 2)
// 	fmt.Println(sprite)
// 	fmt.Println(err)
// 	// fmt.Println(s.SSheet.Sheet_ID)
// }
func (s *SPRITE_LOADER) INIT() {
	s.SpriteNameToSpriteID = map[string]int{}
}

func (s *SPRITE_LOADER) Load_Sprite_File(sheetPath string) (int, error) {
	var file *os.File
	var err error
	defer func() {
		defer file.Close()
	}()
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	file, err = os.Open(sheetPath) //"./sprite.png"
	if err != nil {
		fmt.Println("Error: File could not be opened")
		return 0, err
	}

	pixels, err1 := GetPixels(file)
	if err1 != nil {
		fmt.Println("Error While getting image pixels")
		return 0, err1
	}

	per_sprite_dimensions := Dimensions{32, 32}
	if len(pixels)%32 != 0 || len(pixels[0])%32 != 0 {
		return 0, errors.New("Sprite Sheet doesn't have a valid sprite dimension i.e 32")
	}

	sprite_sheet_dimensions := Get_Sheet_Dimensions_By_Sprite(pixels, per_sprite_dimensions)

	s.Get_Internal_Sprite_Sheet(pixels, sprite_sheet_dimensions, 32)
	s.Assign_Sheet_ID()
	return s.SSheet.Sheet_ID, nil
	// fmt.Println(sprite_file_set[0].Name())

}

func (s *SPRITE_LOADER) Get_Internal_Sprite_Sheet(_sourceSheet [][]Pixel, _spriteSheetDimension Dimensions, _perSpriteDimension int) {

	for row := 0; row < len(_sourceSheet); row += _perSpriteDimension {
		for col := 0; col < len(_sourceSheet[row]); col += _perSpriteDimension {
			//fmt.Println("row: " + strconv.Itoa(row) + " col: " + strconv.Itoa(col))
			var tempSprite [][]Pixel
			for pixel := 0; pixel < _perSpriteDimension; pixel++ {
				temp := _sourceSheet[row+pixel][col : col+_perSpriteDimension]
				// fmt.Println(temp)
				tempSprite = append(tempSprite, temp)
			}
			s.Add_Sprite_To_Sheet(tempSprite)

		}
	}
}

func Get_Sheet_Dimensions_By_Sprite(_sourceImage [][]Pixel, _pixelPerSprite Dimensions) Dimensions {
	if len(_sourceImage) < 1 {
		return Dimensions{}
	}

	XDimensions := len(_sourceImage) / _pixelPerSprite.Row
	YDimensions := len(_sourceImage[0]) / _pixelPerSprite.Column

	return Dimensions{XDimensions, YDimensions}
}

func (s *SPRITE_LOADER) GetSprite(_spriteSheetId int, _spriteName string, _xpos int, _ypos int) (Sprite, error) {
	if s.SSheet.Sheet_ID != _spriteSheetId {
		return Sprite{}, errors.New("Sprite Sheet Not Found")
	}
	if len(s.SSheet.Sprite_Sheet) <= _xpos || len(s.SSheet.Sprite_Sheet[0]) <= _ypos {
		return Sprite{}, errors.New("Sprite Position out of bound of Sprite Sheet")
	}
	sprite := s.SSheet.Sprite_Sheet[_xpos][_ypos]
	if sprite.Name != _spriteName {
		return sprite, errors.New("Sprite Name Mismatched. Sprite is returned according to sprite position")
	}
	return sprite, nil
}

func (s *SPRITE_LOADER) GetSpriteIdByName(spritename string) (int, error) {
	value, exists := s.SpriteNameToSpriteID[spritename]
	if exists {
		return value, nil
	} else {
		return 0, errors.New("Given Sprite Name not found")
	}
}

func (s *SPRITE_LOADER) Add_Sprite_To_Sheet(_sprite [][]Pixel) {
	sheetPosition, sheetId, sheetSpaceAvaiable := Check_SpriteSheet_Space(s.SSheet)
	if !sheetSpaceAvaiable {
		//Create a NEW SHEET
		// sheetId = Generate_Int_ID()
		// s.Sprite_Sheets = append(s.Sprite_Sheets, Sprite_Sheet{sheetId, [32][32]Sprite{}})
		// sheetPosition = Dimensions{0, 0}
		panic("Internal Sheet Space Not Available")

	}

	if s.SSheet.Sheet_ID == sheetId {
		spriteName := "SPRITE_" + strconv.Itoa(sheetPosition.Row) + " X " + strconv.Itoa(sheetPosition.Column)
		spriteId := s.Generate_Int_ID()
		s.SSheet.Sprite_Sheet[sheetPosition.Row][sheetPosition.Column] = Sprite{spriteId, spriteName, sheetPosition.Row, sheetPosition.Column, _sprite}
		s.SpriteNameToSpriteID[spriteName] = spriteId
	}

}

func (s *SPRITE_LOADER) Generate_Int_ID() int {
	s.SpriteID = s.SpriteID + 1
	return s.SpriteID

}

func Check_SpriteSheet_Space(currentSheet Sprite_Sheet) (Dimensions, int, bool) {

	temp_sheet := currentSheet.Sprite_Sheet

	for row := range temp_sheet {
		for col := range temp_sheet[row] {
			if temp_sheet[row][col].Sprite_ID == 0 {
				return Dimensions{row, col}, currentSheet.Sheet_ID, true
			}
		}
	}

	return Dimensions{}, 0, false
}

func GetPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			// fmt.Println(RgbaToPixel(img.At(x, y).RGBA()))
			row = append(row, RgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func RgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)}
}

// classes
type Sprite_Sheet struct {
	Sheet_ID     int
	Sprite_Sheet [32][32]Sprite
}

type Sprite struct {
	Sprite_ID int
	Name      string
	XPos      int
	YPos      int
	Data      [][]Pixel
}

type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Dimensions struct {
	Row    int
	Column int
}
