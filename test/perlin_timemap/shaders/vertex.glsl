#version 410
in vec2 aPos;

uniform mat4 ortho;
uniform mat4 world;
void main()
{
    gl_Position = ortho*world*vec4(aPos,0.0,1.0);
}