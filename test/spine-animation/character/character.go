package character

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/test/spine-animation/animation"
	"github.com/skycoin/cx-game/test/spine-animation/cx"

	//"github.com/hajimehoshi/ebiten"

	_ "image/png"
)

type Character struct {
	Time  float64
	Play  bool
	Speed float64

	// TODO: replace this with atlas
	ImagesPath string
	Images     map[string]*cx.PictureData

	Skeleton  *spine.Skeleton
	Animation *spine.Animation

	SkinIndex      int
	AnimationIndex int

	DebugCenter bool
	DebugBones  bool
}

func LoadCharacter(loc animation.Location) (*Character, error) {
	rd, err := os.Open(loc.JSON)
	if err != nil {
		return nil, err
	}

	data, err := spine.ReadJSON(rd)
	if err != nil {
		return nil, err
	}
	data.Name = loc.Name

	char := &Character{}

	char.ImagesPath = loc.Images

	char.Images = make(map[string]*cx.PictureData)
	fmt.Printf("%v", char.Images)
	char.Play = true
	char.DebugBones = false
	char.DebugCenter = false

	char.Speed = 1
	char.Skeleton = spine.NewSkeleton(data)
	char.Skeleton.Skin = char.Skeleton.Data.DefaultSkin
	char.Animation = char.Skeleton.Data.Animations[0]

	char.AnimationIndex = 0
	char.SkinIndex = 0

	char.Skeleton.FlipY = true

	char.Skeleton.UpdateAttachments()
	char.Skeleton.Update()
	fmt.Printf("%v", char)
	return char, nil

}

func (char *Character) Description() string {
	return char.Skeleton.Data.Name + " > " + char.Skeleton.Skin.Name + " > " + char.Animation.Name
}

func (char *Character) NextAnimation(offset int) {
	char.AnimationIndex += offset
	for char.AnimationIndex < 0 {
		char.AnimationIndex += len(char.Skeleton.Data.Animations)
	}
	char.AnimationIndex = char.AnimationIndex % len(char.Skeleton.Data.Animations)
	char.Animation = char.Skeleton.Data.Animations[char.AnimationIndex]
	char.Skeleton.SetToSetupPose()
	char.Skeleton.Update()
}

func (char *Character) NextSkin(offset int) {
	char.SkinIndex += offset
	for char.SkinIndex < 0 {
		char.SkinIndex += len(char.Skeleton.Data.Skins)
	}
	char.SkinIndex = char.SkinIndex % len(char.Skeleton.Data.Skins)
	char.Skeleton.Skin = char.Skeleton.Data.Skins[char.SkinIndex]
	char.Skeleton.SetToSetupPose()
	char.Skeleton.Update()
	char.Skeleton.UpdateAttachments()
}
func (char *Character) Update(dt float64, x, y float64) {
	if char.Play {
		char.Time += dt * char.Speed
	}

	char.Skeleton.Local.Translate.Set(float32(x), float32(y))
	// char.Skeleton.Local.Scale.Set(0.5, 0.5)
	char.Animation.Apply(char.Skeleton, float32(char.Time), true)
	char.Skeleton.Update()
}

func (char *Character) GetImage(attachment, path string) *cx.PictureData {
	if path != "" {
		attachment = path
	}
	if pd, ok := char.Images[attachment]; ok {
		return pd
	}
	fmt.Println("Loading " + attachment)

	fallback := func() *cx.PictureData {
		fmt.Println("missing: ", attachment)

		m := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for i := range m.Pix {
			m.Pix[i] = 0x80
		}
		pd := cx.PictureDataFromImage(m)
		char.Images[attachment] = pd
		return pd
	}

	fullpath := filepath.Join(char.ImagesPath, attachment+".png")
	file, err := os.Open(fullpath)
	if err != nil {
		return fallback()
	}

	m, _, err := image.Decode(file)
	if err != nil {
		return fallback()
	}
	pd := cx.PictureDataFromImage(m)

	char.Images[attachment] = pd

	return pd
}

// func (char *Character) GetImage(attachment, path string) *image.RGBA {
// if path != "" {
// 	attachment = path
// }
// if pd, ok := char.Images[attachment]; ok {
// 	return pd
// }
// fmt.Println("Loading " + attachment)

// fallback := func() *image.RGBA {
// 	fmt.Println("missing: ", attachment)

// 	m := image.NewRGBA(image.Rect(0, 0, 10, 10))
// 	for i := range m.Pix {
// 		m.Pix[i] = 0x80
// 	}

// //	pd, _ := image.NewRGBA(m, 0)
// //	char.Images[attachment] = pd
// 	return char.Images[attachment]
// }

// fullpath := filepath.Join(char.ImagesPath, attachment+".png")
// file, err := os.Open(fullpath)
// if err != nil {
// 	return fallback()
// }

// m, _, err := image.Decode(file)
// if err != nil {
// 	return fallback()
// }
// pd, _ := new(m, 0)

// //char.Images[attachment] = pd

//	return char.Images[attachment]
// }

// func (char *Character) Draw(target *ebiten.Image) {
// 	for _, slot := range char.Skeleton.Order {
// 		bone := slot.Bone
// 		switch attachment := slot.Attachment.(type) {
// 		case nil:
// 		case *spine.RegionAttachment:
// 			local := attachment.Local.Affine()
// 			final := bone.World.Mul(local)

// 			// BUG: inconvenient
// 			var geom ebiten.GeoM
// 			geom.SetElement(0, 0, float64(final.M00))
// 			geom.SetElement(0, 1, float64(final.M01))
// 			geom.SetElement(0, 2, float64(final.M02))
// 			geom.SetElement(1, 0, float64(final.M10))
// 			geom.SetElement(1, 1, float64(final.M11))
// 			geom.SetElement(1, 2, float64(final.M12))

// 			m := char.GetImage(attachment.Name, attachment.Path)
// 			box := m.Bounds()

// 			var draw ebiten.DrawImageOptions
// 			// flip image and set origin to center
// 			draw.GeoM.Translate(-float64(box.Dx())*0.5, -float64(box.Dy())*0.5)
// 			draw.GeoM.Scale(1, -1)
// 			draw.GeoM.Concat(geom)
// 			// tint the texture
// 			draw.ColorM.Scale(slot.Color.NRGBA64())

// 			// set blending mode
// 			// BUG: incorrect, should use blending mode not compositing mode
// 			switch slot.Data.Blend {
// 			case spine.Normal:
// 				// MISSING
// 			case spine.Additive:
// 				draw.CompositeMode = ebiten.CompositeModeLighter
// 			case spine.Multiply:
// 				// MISSING
// 			case spine.Screen:
// 				// MISSING
// 				draw.CompositeMode = ebiten.CompositeModeLighter
// 			}

// 			target.DrawImage(m, &draw)
// 		default:
// 			panic(fmt.Sprintf("unknown attachment %v", attachment))
// 		}
// 	}
// }
