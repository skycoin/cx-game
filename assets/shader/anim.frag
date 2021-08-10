in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D tex;
uniform mat3 texTransform;
void main() {
	vec2 coord  = vec2(texTransform*vec3(tCoord,1));
	frag_colour = texture(tex,coord);
	if (frag_colour.a<0.1) discard;
}
