in vec2 tCoord;
out vec4 frag_colour;
uniform sampler2D tex;
uniform vec2 texelSize;
uniform vec4 borderColor;
void main() {
	vec4 above = texture(tex,vec2(tCoord.x,tCoord.y-texelSize.y));
	vec4 below = texture(tex,vec2(tCoord.x,tCoord.y+texelSize.y));
	vec4 left = texture(tex,vec2(tCoord.x-texelSize.x,tCoord.y));
	vec4 right = texture(tex,vec2(tCoord.x+texelSize.x,tCoord.y));
	vec4 here = texture(tex,tCoord);

	float others = above.a+below.a+left.a+right.a;
	bool isBorder = here.a < 0.1 && others > 0;
	frag_colour = isBorder ? borderColor : here;
}
