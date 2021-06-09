#version 410
layout (location=0) in vec3 position;
layout (location=1) in vec3 texCoord;
out vec2 tCoord;
uniform mat4 world;
uniform mat4 projection;
void main() {
	gl_Position = projection * world * vec4(position, 1.0);
	tCoord = vec2(texCoord.x,texCoord.y);
}
