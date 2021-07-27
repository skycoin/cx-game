layout (location = 0) in vec2 aPos;
uniform mat4 uProjection;
uniform mat4 uWorld;

void main()
{
   gl_Position = uProjection * vec4(aPos,0, 1.0);
}
