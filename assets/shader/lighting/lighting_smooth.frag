out vec4 frag_colour;

in vec2 texCoords;

uniform sampler2D u_lightmask;


struct Data {
    float scaleW;
    float scaleH;
    float offsetX;
    float offsetY;
};

uniform Data data;
const float gamma_factor = 1.0f / 2.2f;
const vec3 gamma_factor3 = vec3(gamma_factor);


void main(){
    // frag_colour = texture(u_lightmap, vec2(mvp*vec4(texCoords,1,1)));

    vec2 modifiedTexCoords = texCoords * vec2(data.scaleW, data.scaleH)
    +vec2(data.offsetX,data.offsetY);
    
    vec3 out_color = texture(u_lightmask, modifiedTexCoords).rgb;

    // out_color = pow(out_color, gamma_factor3);
    // vec4 gradColor = texture(gradTexture, out_color.z);
    frag_colour = vec4(out_color, 1.0);
}
