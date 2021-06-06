#version 410
in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D tex;
uniform vec4 color;
void main() {
	frag_colour = texture(tex,tCoord) * color;
}
