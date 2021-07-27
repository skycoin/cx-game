in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D tex;
uniform vec4 color;
uniform vec2 scale;
uniform vec2 offset;

uniform float value;

float pi() { return radians(180); }

void main() {
	float angle = atan(tCoord.y-0.5,tCoord.x-0.5);
	if (angle<0) angle += 2*pi();
	if (angle>value) discard;
	frag_colour = texture(tex,scale*tCoord+offset) * color;
	//frag_colour = vec4(angle,0,0,1);
}
