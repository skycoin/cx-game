
out vec4 frag_colour;
in vec2 texCoords;
uniform sampler2D u_screenTexture;

uniform float u_ratio;
// uniform bool offset;

struct ShockwaveData {
	vec2 center;
	float force;
	float size;
	float thickness;
};

uniform ShockwaveData data;
uniform bool u_shockwave_enabled;

void main(){

	vec2 tCoords = texCoords;
	if (u_shockwave_enabled){
		vec2 disp = normalize(texCoords-data.center)*data.force;
   		 float mask = (1-smoothstep(data.size-0.1, data.size,length(texCoords-data.center)))*smoothstep(data.size-data.thickness-0.1, data.size-data.thickness,length(texCoords-data.center));


		tCoords = tCoords - disp*mask;
	}
    // frag_colour = texture(u_texture, texCoords - vec2(0.5, 0.5));
    // frag_colour = vec4(texCoords-disp, 0,1);
    frag_colour = texture(u_screenTexture, tCoords);


    // frag_colour.rgb= vec3(mask);
    // frag_colour = vec4(texCoords, 0, 1);
    // frag_colour = vec4( texCoords, 0,1 );
    
}