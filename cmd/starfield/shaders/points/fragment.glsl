#version 410
out vec4 frag_colour;

uniform sampler2D ourTexture;
uniform sampler1D texture_1d;
uniform float gradValue;
uniform float intensity;

void main() {
    vec4 tex1d = texture(texture_1d, gradValue);
    frag_colour = vec4(tex1d.xyz,intensity);
    // frag_colour = vec4(1.0,1.0,0.0, 1.0);
}