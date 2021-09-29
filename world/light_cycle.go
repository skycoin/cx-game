package world

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/cxmath"
)

const (
	LIGHT_FALLOFF_FACTOR = 0.75
)

type LightTextureGenerator struct {
	//for smooth lighting, todo better approach
	lightingCurveImage *image.RGBA
	lightingCurveTex   uint32
}

func NewLightTextureGenerator() LightTextureGenerator {
	return LightTextureGenerator{
		lightingCurveImage: image.NewRGBA(image.Rect(0, 0, 16, 16)),
	}
}

func (ltg *LightTextureGenerator) GenerateLightTexture(lightValue float32) uint32 {
	if ltg.lightingCurveTex == 0 {
		gl.GenTextures(1, &ltg.lightingCurveTex)
		gl.BindTexture(gl.TEXTURE_2D_ARRAY, ltg.lightingCurveTex)

		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

		gl.TexImage3D(gl.TEXTURE_2D_ARRAY, 0, gl.RGBA, 1, 1, 16*16, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	}

	d2 := mgl32.Vec3{0, 1, 1}

	L1 := [16]mgl32.Vec3{} //skylight
	L2 := [16]mgl32.Vec3{} //envlight

	for i := 0; i < 16; i++ {
		factor := falloff(15-i, LIGHT_FALLOFF_FACTOR)
		L1[i] = cxmath.Vec3ScalarMult(getTwist(i, lightValue), factor)
	}

	for i := 0; i < 16; i++ {
		factor := falloff(15-i, LIGHT_FALLOFF_FACTOR)
		L2[i] = cxmath.Vec3ScalarMult(d2, factor)
	}

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			t3 := mgl32.Vec3{}

			// _j := 15 - j
			t3[0] = lightValue*L1[i].X() + L2[j].X()
			t3[1] = lightValue*L1[i].Y() + L2[j].X()
			t3[2] = lightValue*L1[i].Z() + L2[j].Z()

			r := 255 * t3[0]
			g := 255 * t3[1]
			b := 255 * t3[2]

			if r > 255 {
				r = 255
			}
			if g > 255 {
				g = 255
			}
			if b > 255 {
				b = 255
			}
			ltg.lightingCurveImage.Set(
				j, i,
				color.RGBA{
					uint8(r),
					uint8(g),
					uint8(b),
					255,
				},
			)

		}
	}

	ltg.lightingCurveImage.Set(
		0, 0, color.RGBA{5, 5, 8, 255},
	)

	// if int(glfw.GetTime())%3 == 0 {
	// 	ltg.LightTextureIntoFile()
	// }
	/*


	   for (int i=0; i<dim; i++)
	   {

	       for (int j=0; j<dim; j++)
	       {
	           //int _i = i;
	           const int _j = 15-j;

	           struct Vec3 t3;
	           t3.x = lightv*L1[i].x + L2[_j].x;
	           t3.y = lightv*L1[i].y + L2[_j].y;
	           t3.z = lightv*L1[i].z + L2[_j].z;

	           values[3*(dim*_j+i)+0] = t3.x;
	           values[3*(dim*_j+i)+1] = t3.y;
	           values[3*(dim*_j+i)+2] = t3.z;
	       }
	   }

	*/

	//upload to texture array
	tilesX, tilesY := 16, 16
	tileSizeX := 4 // r,g,b,a
	rowLen := tileSizeX * tilesX
	for iy := 0; iy < tilesY; iy++ {
		var result []uint8
		i := iy * 16
		result = ltg.lightingCurveImage.Pix[rowLen*iy : rowLen*(iy+1)]
		gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, 0, 0, 0, int32(i), 1, 1, 16, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(result))
	}

	return ltg.lightingCurveTex
}

func falloff(it int, factor float32) float32 {
	result := float32(1.0)
	for i := 0; i < it; i++ {
		result *= factor
	}
	return result
}

func getTwist(i int, lightv float32) mgl32.Vec3 {
	//white light
	b := mgl32.Vec3{1, 1, 1}
	//gamma danger twist
	a := mgl32.Vec3{1.6, 0.4, 0.4}

	twist_start := float32(0.5)
	if lightv < twist_start {
		return b
	}
	lightv -= twist_start
	lightv /= (1 - twist_start)

	if lightv <= 0.0-0.001 || lightv >= 1.0+0.001 {
		fmt.Printf("ERROR: lightv = %v\n", lightv)
	}

	if i <= 10 {
		return b
	}
	if i == 11 {
		return cxmath.Vec3Mix(b, cxmath.Vec3Mix(b, a, 0.2), lightv)
	}
	if i == 12 {
		return cxmath.Vec3Mix(b, cxmath.Vec3Mix(b, a, 0.4), lightv)
	}
	if i == 13 {
		return cxmath.Vec3Mix(b, cxmath.Vec3Mix(b, a, 0.6), lightv)
	}
	if i == 14 {
		return cxmath.Vec3Mix(b, cxmath.Vec3Mix(b, a, 0.8), lightv)
	}
	if i == 15 {
		return cxmath.Vec3Mix(b, cxmath.Vec3Mix(b, a, 1.0), lightv)
	}
	return b
}

//get z offset for texture-2d-array
func (ltg *LightTextureGenerator) GetZOffset(skyLight, envLight uint8) float32 {
	return float32(skyLight)*16 + float32(envLight)
}

func (ltg *LightTextureGenerator) LightTextureIntoFile() {
	f, err := os.Create("daylight.png")
	defer f.Close()
	if err != nil {
		return
	}
	png.Encode(f, ltg.lightingCurveImage)
}
