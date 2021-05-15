#version 410
out vec4 frag_colour;

uniform sampler2D ourTexture;
uniform sampler1D texture_1d;
uniform float gradValue;

void main() {
    vec4 tex1d = texture(texture_1d, gradValue);
    frag_colour = tex1d;
}