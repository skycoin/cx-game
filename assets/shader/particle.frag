#version 410
in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D tex;
uniform vec4 color;
uniform vec2 scale;
uniform vec2 offset;
void main() {
	frag_colour = texture(tex,scale*tCoord+offset) * color;
	frag_colour = vec4(1,1,1,1);
}
