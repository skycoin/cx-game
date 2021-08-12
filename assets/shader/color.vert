// use typical layout even though we don't need texcoord
layout (location=0) in vec3 position;
layout (location=1) in vec3 texCoord;

uniform mat4 projection;
uniform mat4 modelviews[NUM_INSTANCES];

flat out int instance;

void main() {
	gl_Position = projection * modelviews[gl_InstanceID] * vec4(position, 1.0);
	instance = gl_InstanceID;
}
