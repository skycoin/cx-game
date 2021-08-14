
layout (location=0) in vec4 posAndTexCoord;

out vec2 texCoord;

uniform mat4 projection;
uniform mat4 world;
uniform mat4 view;

void main(){
    gl_Position = projection *view* world * vec4(posAndTexCoord.xy,0,1);
    texCoord = posAndTexCoord.zw;
}
