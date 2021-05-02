#version 410
in vec2 tCoord;
out vec4 frag_colour;

uniform float gradvalue;
uniform sampler2D ourTexture;
uniform sampler1D texture_1d;
uniform sampler2D texture_2d;


void main() {
    frag_colour = mix(texture(ourTexture, vec2(tCoord.x, tCoord.y)), texture(texture_1d, 0.3), 0.2);
    // frag_colour = texture(texture_1d, gradvalue);
    // frag_colour = texture(texture_2d, tCoord);
    // frag_colour = vec4(1.0, 0.0, 0.0, 1.0);
    // frag_colour = texture(ourTexture, vec2(tCoord.x, tCoord.y));
}
