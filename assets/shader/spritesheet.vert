layout (location=0) in vec3 position;
layout (location=1) in vec2 texcoord;

out vec2 spriteCoord;
flat out int instance;
uniform mat4 MVPs[NUM_INSTANCES];
uniform mat4 SpriteModels[NUM_INSTANCES];

void main() {
	gl_Position = 
		MVPs[gl_InstanceID] * 
		SpriteModels[gl_InstanceID] * 
		vec4(position, 1.0) ;

	instance = gl_InstanceID;
	spriteCoord = texcoord;
}
