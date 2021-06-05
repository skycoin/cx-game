#version 410
in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D tex;
uniform vec4 color;
uniform float stop;
void main() {
	if (tCoord.x > stop) discard;
	frag_colour = texture(tex,tCoord) * color;
}
