layout (location=0) in vec3 position;
layout (location=1) in vec2 texCoord;
out vec2 tCoord;
uniform mat4 projection;
void main() {
	gl_Position = projection * vec4(position,1);
	tCoord = vec2(texCoord.x,texCoord.y);
}
