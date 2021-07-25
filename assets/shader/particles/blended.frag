#version 410 core

out vec4 frag_colour;
in vec2 texCoord;

uniform sampler2D particle_texture;
//used for dissappearing
uniform vec4 color;

void main(){
    frag_colour = texture(particle_texture, texCoord)*color;
    // frag_colour = vec4(1,0,0,1);
    // frag_colour = vec4(texCoord, 0, 1);
}