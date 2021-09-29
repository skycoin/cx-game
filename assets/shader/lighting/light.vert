

layout (location = 0) in vec4 aData;

uniform mat4 projection;
uniform mat4 model;
uniform mat4 view;


out vec2 texCoords;


void main(){
    gl_Position = projection * view * model * vec4(aData.xy, 0, 1.0);
    texCoords = aData.zw;
}   