// use typical layout even though we don't need texcoord
layout (location=0) in vec3 position;
layout (location=1) in vec3 texCoord;
uniform mat4 mvp;
void main() {
	gl_Position = mvp * vec4(position, 1.0);
}
