#version 410 
out vec4 FragColor;
in vec3 color;
in vec2 TexCoord;

// texture samplers
uniform sampler1D texture1;
uniform sampler2D texture2;
uniform float mixValue;
uniform float randColor;

void main()
{
    // linearly interpolate between both textures (80% container, 20% awesomeface)
    // FragColor = mix(texture(texture1, TexCoord), texture(texture2, TexCoord), mixValue);
    FragColor = texture(texture1, randColor);
    // FragColor = vec4(color,1.0);

}