package spriteloader

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

func LoadSpriteSheet(filepath string) (int, error) {
	spriteSheet, err := NewSpriteSheet(filepath)
	if err != nil {
		return 0, err
	}
}

func LoadSprite(spriteSheetId int, spriteName string, xpos, ypos int) {

}

func GetSpriteIdByName(spriteName string) (int, error) {

}
