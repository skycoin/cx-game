in vec2 tCoord;
out vec4 frag_colour;


uniform sampler2D u_screenTexture;
void main() {
	frag_colour = texture(u_screenTexture, tCoord);
	// frag_colour = vec4(1,0,0,1);
}
