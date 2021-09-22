out vec4 frag_colour;

in vec2 texCoords;

uniform sampler2D u_lightmap;

void main(){
    frag_colour = texture(u_lightmap, texCoords);
}