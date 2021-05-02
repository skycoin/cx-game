#version 410

in vec2 tCoord;
out vec4 frag_colour;

uniform sampler2D nebulaText;
uniform sampler2D gradientText;

void main() {
	vec4 texel = texture(nebulaText, tCoord);
	ivec2 size = textureSize(gradientText, 0);
	vec2 uv = vec2(texel.g, 0.5);
	frag_colour = texture(gradientText, uv);
}

