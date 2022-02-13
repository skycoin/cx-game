#shader vertex
#version 330 core

layout(location = 0) in vec4 position;
layout(location = 1) in vec2 texCoord;
layout(location = 2) in float texIndex;


uniform mat4 u_MVP;

out vec2 v_TexCoord;
out float v_TexIndex;

void main(){
    	gl_Position = u_MVP * position;
        v_TexCoord = texCoord;
        v_TexIndex = texIndex;
};
#shader fragment
#version 330 core

layout(location = 0) out vec4 color;

in vec2 v_TexCoord;
in float v_TexIndex;

uniform vec4 u_Color;
uniform sampler2D u_Texture[2];

void main(){
    int index = int(v_TexIndex);
    vec4 texColor = texture(u_Texture[index], v_TexCoord);
    color = texColor;
};