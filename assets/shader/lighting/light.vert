

layout (location = 0) in vec2 aPos;

uniform mat4 projection;
uniform mat4 model;
uniform mat4 view;




void main(){
    gl_Position = projection * view * model * vec4(aPos, 0, 1.0);
}   