package particle_draw

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/render"
)

type ParticleDrawHandler func([]*particles.Particle, *camera.Camera)

var ParticleDrawHandlerList [constants.NUM_PARTICLE_DRAW_HANDLERS]ParticleDrawHandler
var ParticleShaderList map[types.ParticleDrawHandlerId]*render.Shader = make(map[types.ParticleDrawHandlerId]*render.Shader)

func AssertAllDrawHandlersRegistered() {
	for _, handler := range ParticleDrawHandlerList {
		if handler == nil {
			log.Fatalln("Did not initialize all particle draw handlers")
		}
	}
}

func RegisterDrawHandler(id types.ParticleDrawHandlerId, handler ParticleDrawHandler) {
	ParticleDrawHandlerList[id] = handler
	//register shaders
	fmt.Println("THIS HAS BEEN CALLED WITH ", id)
	RegisterShader(id)
}

func RegisterShader(id types.ParticleDrawHandlerId) {
	switch id {
	case constants.PARTICLE_DRAW_HANDLER_NULL:
		ParticleShaderList[id] = nil
	case constants.PARTICLE_DRAW_HANDLER_SOLID:
		shader := render.NewShader("./assets/shader/particles/solid.vert", "./assets/shader/particles/solid.frag")
		ParticleShaderList[id] = shader
	case constants.PARTICLE_DRAW_HANDLER_TRANSPARENT:
		shader := render.NewShader("./assets/shader/particles/blended.vert",
			"./assets/shader/particles/blended.frag")
		ParticleShaderList[id] = shader
	case constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED:
		shader := render.NewShader("./assets/shader/particles/blended_instanced.vert", "./assets/shader/particles/blended_instanced.frag")
		ParticleShaderList[id] = shader
	}
}

func GetShader(id types.ParticleDrawHandlerId) *render.Shader {
	shader, ok := ParticleShaderList[id]

	if !ok {
		log.Fatal("GET SHADER FAILED!", fmt.Sprintf("shader, look %v", id))
	}
	return shader
}

func GetDrawHandler(id types.ParticleDrawHandlerId) ParticleDrawHandler {
	return ParticleDrawHandlerList[id]
}

func makeQuadVao() uint32 {
	vertices := []float32{
		-0.5, -0.5, 0, 1,
		-0.5, 0.5, 0, 0,
		0.5, -0.5, 1, 1,

		-0.5, 0.5, 0, 0,
		0.5, -0.5, 1, 1,
		0.5, 0.5, 1, 0,
	}
	gl.GenVertexArrays(1, &quad_vao)
	gl.BindVertexArray(quad_vao)

	var vao uint32
	gl.GenBuffers(1, &vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vao)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))

	return vao
}
