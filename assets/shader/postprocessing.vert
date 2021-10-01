layout (location=0) in vec3 aPosition;
layout (location=1) in vec2 aTexCoord;
out vec2 texCoords;
void main() {
	gl_Position = vec4(aPosition,1);
	texCoords = vec2(aTexCoord.x,aTexCoord.y);
}
