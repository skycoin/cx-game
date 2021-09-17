package world

import (
	"log"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/render"
)

type LightValue uint8

//upper 4 bits
func (l LightValue) GetSkyLight() uint8 {
	return uint8((l & 0xf0) >> 4)
}

func (l *LightValue) SetSkyLight(value uint8) {
	*l = (*l & 0xf) | LightValue(value)<<4
}

//lower 4 bits
func (l LightValue) GetEnvLight() uint8 {
	return uint8(l & 0xf)
}
func (l *LightValue) SetEnvLight(value uint8) {
	*l = (*l & 0xf0) | LightValue(value)
}

func (l LightValue) MaxLightValue() uint8 {
	if l.GetSkyLight() > l.GetEnvLight() {
		return l.GetSkyLight()
	}
	return l.GetEnvLight()
}

type tilePos struct {
	X int
	Y int
}

var lightQuad render.Program
var lightVao uint32
var lightShader render.Program

var skyLightUpdateQueue []tilePos = make([]tilePos, slLengthMax)
var slStartIndex int = 0
var slNum int = 0
var slLengthMax int = 10000
var neighbourCount int = 4

func getLightAttenuation(tile *Tile) types.LightAttenuationType {
	switch tile.Name {
	case "":
		return constants.ATTENUATION_DEFAULT
	default:
		return constants.ATTENUATION_SOLID
	}
}
func isSolid(tile *Tile) bool {
	return tile.Name != ""
}

func (planet *Planet) InitLighting() {
	planet.InitSkyLight()

	lightConfig := render.NewShaderConfig("./assets/shader/light.vert", "./assets/shader/light.frag")
	lightShader = lightConfig.Compile()

	// var lightVao uint32
	gl.GenVertexArrays(1, &lightVao)
	gl.BindVertexArray(lightVao)

	var lightVbo uint32
	gl.GenBuffers(1, &lightVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, lightVbo)
	vertices := []float32{
		-0.5, -0.5,
		-0.5, 0.5,
		0.5, -0.5,

		-0.5, 0.5,
		0.5, -0.5,
		0.5, 0.5,
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

//assume light attenuation for all block is 1
func (planet *Planet) InitSkyLight() {
	//init skylight
	for x := 0; x < int(planet.Width); x++ {
		y := int(planet.Height) - 1
		for ; y >= 0; y-- {
			tile := planet.GetTile(x, y, TopLayer)
			//tile not empty
			if tile.Name != "" {
				break
			}
			idx := planet.GetTileIndex(x, y)
			planet.LightingValues[idx].SetSkyLight(15)
			// planet.PushSkyLightUpdate(x, y)
		}
		for ; y >= 0; y-- {
			// tile := planet.GetTile(x, y, TopLayer)
			//tile not empty
			idx := planet.GetTileIndex(x, y)

			topTileIdx := planet.GetTileIndex(x, y+1)
			topTileValue := planet.LightingValues[topTileIdx]
			lightTile := planet.GetTile(x, y, TopLayer)

			attenuation := uint8(getLightAttenuation(lightTile))
			//avoid overflow
			if topTileValue.GetSkyLight() >= attenuation {
				planet.LightingValues[idx].SetSkyLight(topTileValue.GetSkyLight() - attenuation)
			} else {
				planet.LightingValues[idx].SetSkyLight(topTileValue.GetSkyLight())
			}
			planet.PushSkylight(x, y)
		}

	}

	startTimer := time.Now()
	for slNum != 0 && time.Since(startTimer).Seconds() < 1 {
		planet.UpdateSkyLight(64 * 1024)
	}
	if time.Since(startTimer).Seconds() > 1 {
		log.Println("SKYLIGHT UPDATE TOOK TOO LONG")
	}

}

func (planet *Planet) LightUpdateBlock(xtile, yTile int) {

	idx := planet.GetTileIndex(xtile, yTile)
	if idx == -1 {
		log.Println("ERROR: PLACED AT WRONG COORDINATES")
		return
	}

	tile := planet.GetTile(xtile, yTile, TopLayer)
	tileType, ok := GetTileTypeByID(tile.TileTypeID)
	if ok && tile.LightSource {
		if tileType.LightAmount > 15 {
			log.Fatalf("ERROR: max light value is 15")
		}
		planet.LightingValues[idx].SetEnvLight(tileType.LightAmount)
	}

	planet.PushSkylight(xtile, yTile)
	planet.PushEnvLight(xtile, yTile)
}

func (planet *Planet) PushSkylight(xTile, yTile int) {
	if slNum == slLengthMax {
		return
	}
	//adds to queue
	skyLightUpdateQueue[(slStartIndex+slNum)%slLengthMax] = tilePos{xTile, yTile}
	slNum++
}

var neighboursOffsets = []int{
	0, 1, //n
	-1, 0, //w
	1, 0, //e
	0, -1, //s
	-1, 1, //nw
	1, 1, //ne
	-1, -1, //sw
	1, -1, //se
}

func (planet *Planet) UpdateSkyLight(iterations int) {
	//only top tile is accounted in calculations
	//do update tile logic on tiles in update queue
	if slNum == 0 {
		return
	} else if slNum < 0 {
		log.Fatalln("Skylight update error 1")
	}
	//update logic
	for i := 0; i < iterations; i++ {
		if slNum == 0 {
			return
		}

		pos := skyLightUpdateQueue[slStartIndex]

		slStartIndex = (slStartIndex + 1) % slLengthMax
		slNum--
		idx := planet.GetTileIndex(pos.X, pos.Y)
		if idx == -1 {
			continue

		}
		lightVal := planet.LightingValues[idx]
		lightTile := planet.GetTile(pos.X, pos.Y, TopLayer)
		lightSkyLightValue := lightVal.GetSkyLight()

		topTileIdx := planet.GetTileIndex(pos.X, pos.Y+1)
		if topTileIdx == -1 {
			continue
		}

		topTileLightValue := planet.LightingValues[topTileIdx]
		topTile := planet.GetTile(pos.X, pos.Y+1, TopLayer)

		attenuation := uint8(getLightAttenuation(lightTile))
		//placed solid block below sunlight
		if !isSolid(topTile) && topTileLightValue.GetSkyLight() == 15 {
			//toptile is skylight
			if isSolid(lightTile) && lightVal.GetSkyLight() != 15-2 {
				//if checked tile is solid
				planet.LightingValues[idx].SetSkyLight(15 - attenuation)
				for i := 0; i < neighbourCount; i++ {
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
				planet.PushSkylight(pos.X, pos.Y)
			} else if !isSolid(lightTile) && lightVal.GetSkyLight() != 15 {
				//removed block, waiting for sunlight instead
				planet.LightingValues[idx].SetSkyLight(15)
				for i := 0; i < neighbourCount; i++ {
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
				planet.PushSkylight(pos.X, pos.Y)
			} else {
			}
			continue
		}
		// continue
		neighbourIndexes := make([]int, neighbourCount)
		for i := 0; i < neighbourCount; i++ {
			idx := planet.GetTileIndex(
				pos.X+neighboursOffsets[i*2],
				pos.Y+neighboursOffsets[i*2+1],
			)
			neighbourIndexes[i] = idx
		}

		var maxSkylightValue uint8 = 0

		//determine brightest neighbour out of 4
		for i := 0; i < neighbourCount; i++ {
			if neighbourIndexes[i] == -1 {
				continue
			}
			neighbourLightValue := planet.LightingValues[neighbourIndexes[i]]
			if neighbourLightValue.GetSkyLight() > maxSkylightValue {
				maxSkylightValue = neighbourLightValue.GetSkyLight()
			}
		}

		if maxSkylightValue >= attenuation {
			if lightSkyLightValue != maxSkylightValue-attenuation {
				planet.LightingValues[idx].SetSkyLight(maxSkylightValue - attenuation)
				for i := 0; i < neighbourCount; i++ {
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
			}
		} else {
			if lightSkyLightValue != 0 {
				planet.LightingValues[idx].SetSkyLight(0)
				for i := 0; i < neighbourCount; i++ {
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
			}
		}

	}
}

func SwitchNeighbourCount(planet *Planet) {
	if neighbourCount == 4 {
		log.Println("LIGHTING: 8 neighbours")
		neighbourCount = 8
	} else {
		log.Println("LIGHTING: 4 neighbours")
		neighbourCount = 4
	}
	planet.InitSkyLight()
}

var envLightUpdateQueue []tilePos = make([]tilePos, elLengthMax)
var elStartIndex int = 0
var elNum int = 0
var elLengthMax int = 20000

func (planet *Planet) PushEnvLight(xTile, yTile int) {
	if elNum == elLengthMax {
		return
	}
	//adds to queue
	envLightUpdateQueue[(elStartIndex+elNum)%elLengthMax] = tilePos{xTile, yTile}
	elNum++
}

func (planet *Planet) UpdateEnvLight(iterations int) {
	if elNum == 0 {
		return
	} else if elNum < 0 {
		log.Fatalln("EnvLight update error 1")
	}
	// fmt.Println(elNum, "    elnum")
	//update logic
	for i := 0; i < iterations; i++ {
		if elNum == 0 {
			return
		}

		pos := envLightUpdateQueue[elStartIndex]

		elStartIndex = (elStartIndex + 1) % elLengthMax
		elNum--
		idx := planet.GetTileIndex(pos.X, pos.Y)

		if idx == -1 {
			continue

		}
		lightVal := planet.LightingValues[idx]
		lightTile := planet.GetTile(pos.X, pos.Y, TopLayer)

		envLightValue := lightVal.GetEnvLight()

		attenuation := uint8(getLightAttenuation(lightTile))

		neighbourIndexes := make([]int, neighbourCount)
		for i := 0; i < neighbourCount; i++ {
			idx := planet.GetTileIndex(
				pos.X+neighboursOffsets[i*2],
				pos.Y+neighboursOffsets[i*2+1],
			)
			neighbourIndexes[i] = idx
		}

		var maxEnvLightValue uint8 = 0

		//determine brightest neighbour out of 4
		for i := 0; i < neighbourCount; i++ {
			if neighbourIndexes[i] == -1 {
				continue
			}
			neighbourLightValue := planet.LightingValues[neighbourIndexes[i]]
			if neighbourLightValue.GetEnvLight() > maxEnvLightValue {
				maxEnvLightValue = neighbourLightValue.GetEnvLight()
			}
		}

		if lightTile.LightSource {
			for i := 0; i < neighbourCount; i++ {
				if planet.LightingValues[neighbourIndexes[i]].GetEnvLight() < envLightValue {
					planet.PushEnvLight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
			}
			continue
		}
		if maxEnvLightValue >= attenuation {
			if envLightValue != maxEnvLightValue-attenuation {
				planet.LightingValues[idx].SetEnvLight(maxEnvLightValue - attenuation)
				for i := 0; i < neighbourCount; i++ {
					planet.PushEnvLight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
			}
		} else {
			if envLightValue != 0 {
				planet.LightingValues[idx].SetEnvLight(0)
				for i := 0; i < neighbourCount; i++ {
					planet.PushEnvLight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
			}
		}
	}
}

var lightMaskOn bool = false

func (planet *Planet) DrawLightMap(cam *camera.Camera) {
	if !lightMaskOn {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.DST_COLOR, gl.ZERO)
	}
	defer gl.Disable(gl.BLEND)
	lightShader.Use()

	for x := cam.Frustum.Left; x < cam.Frustum.Right; x++ {
		for y := cam.Frustum.Bottom; y < cam.Frustum.Top; y++ {
			if y < 0 {
				continue
			}
			wrappedPos := planet.WrapAround(mgl32.Vec2{float32(x), float32(y)})
			idx := planet.GetTileIndex(int(wrappedPos[0]), int(wrappedPos[1]))
			lightValue := planet.LightingValues[idx]

			lightShader.SetMat4("projection", &render.Projection)
			view := cam.GetViewMatrix()
			lightShader.SetMat4("view", &view)
			model := mgl32.Translate3D(float32(x), float32(y), 0)
			// model = model.Mul4(model.)
			lightShader.SetMat4("model", &model)
			lightShader.SetVec3("color", &mgl32.Vec3{
				// float32(lightValue.GetSkyLight()) / 15,
				// float32(lightValue.GetSkyLight()) / 15,
				// float32(lightValue.GetSkyLight()) / 15,
				float32(lightValue.MaxLightValue()) / 15,
				float32(lightValue.MaxLightValue()) / 15,
				float32(lightValue.MaxLightValue()) / 15,
			})
			gl.BindVertexArray(lightVao)
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
		}
	}

}

func (planet *Planet) UpdateLighting() {
	planet.UpdateSkyLight(1000)
	planet.UpdateEnvLight(1000)
}
