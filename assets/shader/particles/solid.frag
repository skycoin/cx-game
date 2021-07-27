
out vec4 frag_colour;
in vec2 texCoord;

uniform sampler2D particle_texture;
//used for dissappearing
uniform vec4 color;

void main(){
    frag_colour = texture(particle_texture, texCoord)*color;
}
