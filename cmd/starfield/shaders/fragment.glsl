#version 410
in vec2 tCoord;
out vec4 frag_colour;


uniform sampler2D ourTexture;

uniform float mixvalue;

void main() {
    // frag_colour = mix(texture(ourTexture, vec2(tCoord.x, tCoord.y)), texture(texture_1d, 0.3), 0.2);
    // frag_colour = texture(texture_1d, gradvalue);
    // frag_colour = texture(texture_2d, tCoord);
    // tex1 = texture(texture_1d, gradvalue);
    frag_colour = texture(ourTexture, vec2(tCoord.x, tCoord.y));
}   