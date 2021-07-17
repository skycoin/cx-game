#version 330 core
layout (location = 0) in vec3 aPos;
uniform mat4 uProjection;
uniform mat4 uWorld;

void main()
{
   gl_Position = uProjection * vec4(aPos, 1.0);
}
