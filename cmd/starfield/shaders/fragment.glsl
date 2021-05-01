#version 410
in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D ourTexture;
uniform sampler1D texture_1d;
uniform float gradValue;
void main() {
    // frag_colour = mix(texture(texture_2d, tCoord), texture(texture_1d, gradValue),0.5);
    // frag_colour = mix(texture(ourTexture, tCoord), texture(texture_1d, gradValue), 0.5);
    frag_colour = texture(ourTexture, tCoord);
}
