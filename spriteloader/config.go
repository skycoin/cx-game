package spriteloader

import (
	"errors"
	"fmt"
	"os"
	"log"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
)

type SpriteConfig struct {
	Unit string `yaml:"unit"`
	Left int `yaml:"left"`
	Top int `yaml:"top"`
	Width int `yaml:"width"`
	Height int `yaml:"height"`
}

func (sprite SpriteConfig) Model() mgl32.Mat4 {
	return sprite.Scale().Mat4()
}

func (sprite SpriteConfig) Scale() mgl32.Mat3 {
	// fill defaults
	width := float32(sprite.Width)
	if width == 0 { width = 1 }

	height := float32(sprite.Height)
	if height == 0 { height = 1 }

	return mgl32.Scale2D( width, height )
}

func (sprite SpriteConfig) Transform() mgl32.Mat3 {

	translate := mgl32.Translate2D( float32(sprite.Left), float32(sprite.Top) )
	scale := sprite.Scale()

	return translate.Mul3(scale)
}

type SpriteSheetConfig struct {
	Name string `yaml:"name"`
	Width int `yaml:"width"`
	Height int `yaml:"height"`
	CellWidth int `yaml:"cellwidth"`
	CellHeight int `yaml:"cellheight"`
	SpriteConfigs map[string]SpriteConfig `yaml:"sprites"`
	Autoname string `yaml:"autoname"`
    ModelWidth int `yaml:"modelwidth"`
    ModelHeight int `yaml:"modelheight"`
}

func (config *SpriteSheetConfig) Rows() int {
	return config.Width / config.CellWidth
}

func (config *SpriteSheetConfig) Cols() int {
	return config.Height / config.CellHeight
}

func (spritesheetConfig *SpriteSheetConfig) spriteScale(
		spriteConfig SpriteConfig,
) mgl32.Mat3 {
	w := float32(spritesheetConfig.Width)
	h := float32(spritesheetConfig.Height)
	cw := float32(spritesheetConfig.CellWidth)
	ch := float32(spritesheetConfig.CellHeight)

	if spriteConfig.Unit == "grid" {
		return mgl32.Scale2D( cw/w, ch/h )
	}
	if spriteConfig.Unit == "pixel" {
		return mgl32.Scale2D ( 1/w, 1/h )
	}
	if spriteConfig.Unit == "full" {
		return mgl32.Scale2D( 1, 1)
	}

	log.Fatalf("unsupported unit type %v",spriteConfig.Unit)
	return mgl32.Ident3()
}

func (spritesheetConfig *SpriteSheetConfig) Sprite(
		name string, spriteConfig SpriteConfig,
) render.Sprite {
	scale := spritesheetConfig.spriteScale(spriteConfig)
	local := spriteConfig.Transform()
	transform := scale.Mul3(local)

	return render.Sprite {
		Name: name, Transform: transform,
		Model: spriteConfig.Model(),
	}
}

func (spritesheetConfig *SpriteSheetConfig) cells() int {
	cols := spritesheetConfig.Width / spritesheetConfig.CellWidth
	rows := spritesheetConfig.Height / spritesheetConfig.CellHeight
	return cols*rows
}

func (spritesheetConfig *SpriteSheetConfig) autoSprites() []render.Sprite {
	if spritesheetConfig.Autoname == "index" {
		sprites := make([]render.Sprite, spritesheetConfig.cells())
		x := 0; y := 0;
		downScale := mgl32.Scale2D(
			1 / float32(spritesheetConfig.Width),
			1 / float32(spritesheetConfig.Height),
		)
		upScale := mgl32.Scale2D(
			float32(spritesheetConfig.CellWidth),
			float32(spritesheetConfig.CellHeight),
		)
		for idx := range sprites {
			translate := mgl32.Translate2D(float32(x),float32(y))
			sprites[idx] = render.Sprite {
				Name: fmt.Sprintf("%v:%d", spritesheetConfig.Name, idx),
				Model: mgl32.Ident4(),
				Transform: downScale.Mul3(translate).Mul3(upScale),
			}
			x += spritesheetConfig.CellWidth
			if x >= spritesheetConfig.Width {
				x = 0
				y += spritesheetConfig.CellHeight
			}
		}
		return sprites
	}

	log.Fatalf("unrecognized autoname=%v",spritesheetConfig.Autoname)
	return nil
}

func (spritesheetConfig *SpriteSheetConfig) Sprites() []render.Sprite {
	if len(spritesheetConfig.SpriteConfigs)==0 &&
			spritesheetConfig.hasZeroDimensions() {
		// special case for 1:1 img to yaml
		return []render.Sprite{render.Sprite{
			Name: spritesheetConfig.Name,
			Transform: mgl32.Ident3(),
            Model: mgl32.Scale3D(
                float32(spritesheetConfig.ModelWidth),
                float32(spritesheetConfig.ModelHeight),
                1,
            ),
		}}
	}

	if spritesheetConfig.Autoname != "" {
		return spritesheetConfig.autoSprites()
	}

	sprites := make([]render.Sprite, 0, len(spritesheetConfig.SpriteConfigs))

	for name,spriteConfig := range spritesheetConfig.SpriteConfigs {
		sprites = append(sprites, spritesheetConfig.Sprite(name, spriteConfig))
	}
	return sprites
}

func (spritesheetConfig *SpriteSheetConfig) hasNonFullSprite() bool {
	for _,spriteConfig := range spritesheetConfig.SpriteConfigs {
		if spriteConfig.Unit != "full" { return true }
	}
	return false
}

func (config *SpriteSheetConfig) hasZeroDimensions() bool {
	return config.Width==0 && config.Height==0
}

func (config *SpriteSheetConfig) Validate() error {
	if config.hasNonFullSprite() && config.hasZeroDimensions() {
		return errors.New("texture dimensions not set")
	}

	return nil
}

func readSpriteSheetConfig(path string) (SpriteSheetConfig,error) {
	buf, err := os.ReadFile(path)
	if err != nil { return SpriteSheetConfig{},err }
	var config SpriteSheetConfig
	yaml.Unmarshal(buf,&config)
	return config,config.Validate()
}

func LoadSpriteSheetFromConfig(imgPath, cfgPath string) render.SpriteSheet {
	tex := LoadTextureFromFileToGPU(imgPath)
	config,err := readSpriteSheetConfig(cfgPath)
	if err != nil {
		log.Fatalf("load spritesheet from \n%s and %s\n\n%v",imgPath,cfgPath,err)
	}
	return render.SpriteSheet {
		Texture : render.Texture{ gl.TEXTURE_2D, tex.Gl },
		Sprites: config.Sprites(),
	}
}

func RegisterSpritesFromConfig(cfgPath string) []SpriteID {
	imgPath := strings.TrimSuffix(cfgPath, ".yaml") + ".png"
	sheet := LoadSpriteSheetFromConfig(imgPath, cfgPath)
	// TODO deprecate this
	spriteloadersheet := AddSpritesheetFromTexture(sheet.Texture.Texture)
	for _,sprite := range sheet.Sprites {
		render.RegisterSprite(sprite)
		offset := sprite.Transform.Col(2)
		xScale := sprite.Transform.At(0,0)
		yScale := sprite.Transform.At(1,1)
		worldXScale := sprite.Model.At(0,0)
		worldYScale := sprite.Model.At(1,1)

		LoadSpriteF(
			spriteloadersheet,
			sprite.Name, offset.X(), offset.Y(),
			xScale, yScale, worldXScale, worldYScale,
		)
	}
	return []SpriteID{}
}
