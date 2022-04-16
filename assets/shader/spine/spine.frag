in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D tex;

void main(){
    frag_colour = texture(tex,tCoord);
};