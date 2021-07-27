
out vec4 frag_colour;

uniform sampler2D particle_texture;

in vec3 fragData;

void main(){
    frag_colour = texture(particle_texture, fragData.xy) * vec4(1,1,1,fragData.z);
}
