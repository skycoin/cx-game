#version 410
in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D tex;

uniform float left,right,top,bottom;
uniform float width,height;

// stretch the middle of a range
// but leave the left and right ends unaltered
mediump float nineclamp(float w,float left,float right,float x) {
	// return raw coord from left
	if (x<left) return x; 
	// return raw coord from right
	if (x>(w-right)) return (1-right)+(x-w+right); 

	float iw = (w-left-right); // inner width
	float ix = (x-left) / iw; // inner x (normalized)
	return left + ix*(1-left-right);
}

void main() {
	if (tCoord.x > width || tCoord.y > height) discard;
	frag_colour = texture(tex,vec2(
		nineclamp(width,left,right,tCoord.x),
		nineclamp(height,top,bottom,tCoord.y)
	));
}
