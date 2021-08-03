in vec2 spriteCoord;
flat in int instance; 

out vec4 frag_colour;

uniform sampler2D ourTexture;
uniform mat3  texTransforms[NUM_INSTANCES];

void main() {
		vec2 texCoord =
			 vec2(texTransforms[instance] * vec3(spriteCoord,1) );

		frag_colour = texture(ourTexture, texCoord);
}
