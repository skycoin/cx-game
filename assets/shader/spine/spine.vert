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