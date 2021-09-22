out vec4 frag_colour;

in vec2 texCoords;

uniform sampler2D u_lightmap;

struct Data {
    float scaleW;
    float scaleH;
    float offsetX;
    float offsetY;
};

uniform Data data;
uniform sampler1D gradTexture;

uniform bool blue_out;

void main(){
    // frag_colour = texture(u_lightmap, vec2(mvp*vec4(texCoords,1,1)));

    vec2 modifiedTexCoords = texCoords * vec2(data.scaleW, data.scaleH)
    +vec2(data.offsetX,data.offsetY);
    
    vec4 out_color = texture(u_lightmap, modifiedTexCoords);
    // vec4 gradColor = texture(gradTexture, out_color.z);
    frag_colour = out_color;
    
}
