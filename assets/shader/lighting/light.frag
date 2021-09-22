

out vec4 frag_colour;


uniform vec3 color;
uniform sampler1D gradTexture;

void main(){
    vec4 out_color = texture(gradTexture, color.x);
    frag_colour = out_color;
}