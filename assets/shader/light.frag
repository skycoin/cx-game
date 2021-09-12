

out vec4 frag_colour;


uniform vec3 color;

void main(){
    frag_colour = vec4(color,1.0);
}