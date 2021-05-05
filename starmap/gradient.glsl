#version 410

in vec2 tCoord;
out vec4 frag_colour;

uniform sampler2D nebulaTexture;
uniform sampler2D gradientTexture;

void main() {
	vec4 texel = texture(nebulaTexture, tCoord);
	ivec2 size = textureSize(gradientTexture, 0);
	vec2 uv = vec2(texel.g, 0.5);
	frag_colour = texture(gradientTexture, uv);
}

