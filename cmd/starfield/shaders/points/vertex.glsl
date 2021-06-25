#version 410
layout (location=0) in vec3 aPos;
uniform mat4 projection;
uniform mat4 world;
uniform float pointSize;


void main() {
    gl_Position = projection* world * vec4(aPos, 1.0);
    gl_PointSize = 2;
}