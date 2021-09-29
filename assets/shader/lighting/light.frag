

out vec4 frag_colour;


uniform sampler2DArray u_lightmap;

uniform float zoffset;

in vec2 texCoords;

const float gamma_factor = 1.0f / 2.2f;
const vec3 gamma_factor3 = vec3(gamma_factor);

void main(){
    vec3 out_color = texture(u_lightmap, vec3(0,0, zoffset)).rgb;
    // out_color = pow(out_color, gamma_factor3);
    frag_colour = vec4(out_color, 1.0);
}