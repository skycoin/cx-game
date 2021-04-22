package spriteloader

import (
	"errors"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

// Make library
// - spriteloader

// The library should
// - take in one or more sprite files, 32x32 pixels per sprite (make so we can change on init of struct)
// - the sprite files are .png
// - there should be 8x8 sprites per sheet
// - but if there are not, divide x pixels by 32 to get number of sprites in X and divide Y by 32 to get number of sprites in Y
// - each sprite should be given a text name and an x,y position (can be in golang for now, but will read from yaml file later)

// Make an internal spritesheet
// - 32x32 sprites of 32x32
// - copy the sprites from the input files to the internal sprite sheet
// - assign an integer ID for the internal sprite sheet

// - have a map from sprite_name to internal sprite id. And a function
// - have a map from sprite_id to internal sprite id on the sheet. And a function

// Then

// So example

// "spriteSheetId = LoadSpriteSheet("./assets/tiles_01.png")
// "LoadSprite(spriteSheetId, "sprite_name", 1, 2)""

// "id := GetSpriteIdByName("sprite_name")"

// ---

// Make a function that
// - draws a "quad" at position x,y (centered)
// - quad width width, quad height, height
// - with texture id from the internal sheet

// DrawSpriteQuad(xpos, ypos, xwidth, yheight, spriteId)

//for LoadSpriteSheet
var spriteSheetLastId uint32 = 0
var spriteSheetIdMap map[uint32]*image.RGBA

//for sprite to internalsheet map
var spriteLastId uint32 = 0
var spriteIdToInternalSheetIdMap map[uint32]uint32

var spriteNameToInternalSheetIdMap map[string]uint32

var internalSheetLastId uint32 = 0
var internalSheetsMap map[uint32]*InternalSpriteSheet

func LoadSpriteSheet(filepath string) (uint32, error) {
	//assume every sprite_file is 8x8
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}

	im, err := png.Decode(file)
	if err != nil {
		return 0, err
	}

	img := image.NewRGBA(im.Bounds())
	draw.Draw(img, img.Bounds(), im, image.Pt(0, 0), draw.Src)

	//increment id and add to map
	addSpriteSheetToMap(img)
	return spriteSheetLastId, nil
}

func LoadSprite(spriteSheetId uint32, spriteName string, xpos, ypos int) error {
	spriteSheetImage, ok := spriteSheetIdMap[spriteSheetLastId]
	if !ok {
		return errors.New("Could not load sprite into internal sheet")
	}

	tileImage := extractTileFromSpriteSheet(spriteSheetImage, xpos, ypos)

	spriteLastId += 1
	sprite := &Sprite{
		name:      spriteName,
		id:        spriteLastId,
		xpos:      xpos,
		ypos:      ypos,
		imageInfo: tileImage,
	}

	addToInternalSpriteSheet(sprite)
	return nil

}

func GetSpriteIdByName(spriteName string) (uint32, error) {
	id, ok := spriteNameToInternalSheetIdMap[spriteName]
	if !ok {
		return 0, errors.New("no sprite with such name!")
	}
	return id, nil
}
func GetSpriteIdById(spriteId uint32) (uint32, error) {
	id, ok := spriteIdToInternalSheetIdMap[spriteId]
	if !ok {
		return 0, errors.New("no sprite with such id!")
	}
	return id, nil
}

func DrawSpriteQuad(xpos, ypos, xwidth, yheight int, spriteId uint32) {

}

func extractTileFromSpriteSheet(spriteSheetBitmap *image.RGBA, xpos, ypos int) *image.RGBA {
	tileImage := spriteSheetBitmap.SubImage(image.Rect(xpos*32, ypos*32, (xpos+1)*32, (ypos+1)*32))
	return tileImage.(*image.RGBA)
}

func addSpriteSheetToMap(img *image.RGBA) {
	spriteSheetLastId += 1
	spriteSheetIdMap[spriteSheetLastId] = img
}

func addToInternalSpriteSheet(sprite *Sprite) {

}
func getSpriteAt(xpos, ypos int) *image.RGBA {
	file, _ := os.Open("./assets/starfield/stars/planets.png")
	im, _ := png.Decode(file)
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	draw.Draw(img, img.Bounds(), im, image.Pt(0, 0), draw.Src)

	// bitmap := img.Pix

	// bmpReader := bytes.NewReader(bitmap)

	file, err := os.Create("newsprite.png")
	if err != nil {
		log.Fatal(err)
	}

	return img.SubImage(image.Rect(ypos*32, xpos*32, ypos*32+32, xpos*32+32)).(*image.RGBA)
}
