in vec2 tCoord;
out vec4 frag_colour;

uniform sampler2D ourTexture;
uniform sampler1D texture_1d;
uniform float gradValue;
uniform float intensity;

void main() {
    vec4 tex2d = texture(ourTexture, vec2(tCoord.x, tCoord.y));
    vec4 tex1d = texture(texture_1d, gradValue);
    if (tex2d.a < 0.5)
        discard;
    vec4 result = mix(tex2d, tex1d, 0.5);
    frag_colour = vec4(result.xyz + vec3(0.5), intensity);
}
