out vec4 frag_colour;

in vec2 texCoords;

uniform sampler2D u_lightmap;
uniform mat4 mvp;

struct Data {
    float scaleW;
    float scaleH;
    float offsetX;
    float offsetY;
};

uniform Data data;

void main(){
    // frag_colour = texture(u_lightmap, vec2(mvp*vec4(texCoords,1,1)));
    vec2 modifiedTexCoords = texCoords * vec2(data.scaleW, data.scaleH)+vec2(data.offsetX,data.offsetY);
    frag_colour = texture(u_lightmap, modifiedTexCoords);
}