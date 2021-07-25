package particle_draw

import (
	"fmt"
	"log"

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
