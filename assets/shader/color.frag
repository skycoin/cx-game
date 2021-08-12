out vec4 frag_colour;

uniform vec4 colors[NUM_INSTANCES];
flat in int instance;

void main() {
	frag_colour = colors[instance];
}
