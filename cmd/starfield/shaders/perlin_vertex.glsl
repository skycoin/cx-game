#version 410
layout (location = 0) in vec3 aPos;
layout (location = 1) in float aSize;
layout (location = 2) in float aGradient;
out float gradientValue;
uniform mat4 projection;

void main(){
    gl_PointSize = aSize;
    gl_Position = projection*vec4(aPos, 1.0);
    gradientValue = aGradient;
}