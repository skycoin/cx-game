
//xy - xpos ypos, zw - texcoords
layout (location=0) in vec4 aQuad;

//xy =xoffset yoffset, z - scale, w - alpha
layout (location=1) in vec4 aData;


out vec3 fragData;

uniform mat4 projection,view;

void main(){
    gl_Position = projection * view * vec4(aQuad.xy*aData.z + aData.xy, 0, 1);
    fragData = vec3(aQuad.zw,aData.w);
}
