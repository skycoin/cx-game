layout (location=0) in vec3 position;
layout (location=1) in vec2 texcoord;

out vec2 spriteCoord;
flat out int instance;
uniform mat4 projection;
uniform mat4 models[NUM_INSTANCES];
uniform mat4 views[NUM_INSTANCES];

void main() {
	vec4 pos =round(models[gl_InstanceID] * vec4(position,1.0)*64)/64;

	gl_Position = 
		projection *
		views[gl_InstanceID] * pos;

	instance = gl_InstanceID;
	spriteCoord = texcoord;
}
