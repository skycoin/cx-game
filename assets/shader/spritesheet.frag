in vec2 spriteCoord;
flat in int instance; 

out vec4 frag_colour;

uniform sampler2D tex;
uniform mat3 uvtransforms[NUM_INSTANCES];

void main() {
		vec2 texCoord =
			 vec2(uvtransforms[instance] * vec3(spriteCoord,1) );

		frag_colour = texture(tex, texCoord);
		if (frag_colour.a < 0.001) { discard; }
}
