package particle_draw

import (
	"fmt"
	"log"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/constants/particle_constants"
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
	RegisterShader(id)
}

func RegisterShader(id types.ParticleDrawHandlerId) {
	switch id {
	case particle_constants.DRAW_HANDLER_SOLID:
		shader := render.NewShader("./assets/shader/particles/solid.vert", "./assets/shader/particles/solid.frag")
		ParticleShaderList[id] = shader
	case particle_constants.DRAW_HANDLER_ALPHA_BLENDED:
		shader := render.NewShader("./assets/shader/particles/blended.vert",
			".assets/shader/particles/blended.frag")
		ParticleShaderList[id] = shader
	}
}

func GetShader(id types.ParticleDrawHandlerId) *render.Shader {
	shader, ok := ParticleShaderList[id]
	if !ok {
		log.Fatal("GET SHADER FAILED!", fmt.Sprintf("shader, look %v"), id)
	}
	return shader
}

func GetDrawHandler(id types.ParticleDrawHandlerId) ParticleDrawHandler {
	return ParticleDrawHandlerList[id]
}
