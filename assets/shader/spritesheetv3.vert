layout (location=0) in vec3 position;
layout (location=1) in vec2 texcoord;

out vec2 spriteCoord;
flat out int instance;
uniform mat4 projection;
uniform mat4 modelviews[NUM_INSTANCES];
uniform mat4 views[NUM_INSTANCES];

void main() {
	mat4 model = inverse(views[gl_InstanceID]) *modelviews[gl_InstanceID];

	vec4 pos = floor(model * vec4(position,1.0)*32)/32;

	gl_Position = 
		projection *
		views[gl_InstanceID] * pos;

	instance = gl_InstanceID;
	spriteCoord = texcoord;
}
