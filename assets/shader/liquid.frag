// temporary effect to demonstrate that 
// required parameters are being passed
in vec2 tCoord;
in mat4 world;
out vec4 frag_colour;
uniform sampler2D ourTexture;
uniform vec4 color;
uniform float time;
void main() {
		float s = mod(time,1);
		vec4 shimmer = vec4(s,s,s,1);
		//frag_colour = texture(ourTexture, tCoord) * color * shimmer;
		float x = world[3][0];
		frag_colour = vec4(mod(time,1),mod(x/4,1),1,1);
}
