#version 410
out vec4 frag_colour;
in float gradientValue;

uniform sampler1D texture_1d;

void main(){
    frag_colour = texture(texture_1d, gradientValue);
}