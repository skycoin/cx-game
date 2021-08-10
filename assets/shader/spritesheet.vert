layout (location=0) in vec3 position;
layout (location=1) in vec2 texcoord;

out vec2 spriteCoord;
flat out int instance;
uniform mat4 projection;
uniform mat4 modelviews[NUM_INSTANCES];

void main() {
	gl_Position = 
		projection *
		modelviews[gl_InstanceID] * 
		vec4(position, 1.0) ;

	instance = gl_InstanceID;
	spriteCoord = texcoord;
}
