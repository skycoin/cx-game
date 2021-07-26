#version 410
layout (location=0) in vec3 position;
layout (location=1) in vec2 texcoord;
out vec2 tCoord;
uniform mat4 projection;
uniform mat4 worlds[100];
uniform vec2 texScales[100];
uniform vec2 texOffsets[100];
void main() {
	mat4 world = worlds[gl_InstanceID];
	vec2 texScale = texScales[gl_InstanceID];
	vec2 texOffset = texOffsets[gl_InstanceID];
	gl_Position = projection  *  world * vec4(position, 1.0);
	//https://stackoverflow.com/questions/13901119/how-does-vectors-multiply-act-in-shader-language
	
	tCoord = (texcoord+texOffset) * texScale;
}
