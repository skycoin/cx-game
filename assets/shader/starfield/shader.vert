layout (location=0) in vec3 position;
layout (location=1) in vec2 texcoord;
out vec2 tCoord;
uniform mat4 projection;
uniform mat4 world;
uniform vec2 texScale;
uniform vec2 texOffset;
void main() {
    gl_Position = projection * world* vec4(position, 1.0);
    tCoord = (texcoord+texOffset) * texScale;
}
