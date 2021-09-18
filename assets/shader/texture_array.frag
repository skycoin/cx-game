out vec4 frag_colour;

in vec2 spriteCoord;

uniform sampler2DArray u_texture;
uniform mat3 uvtransforms[NUM_INSTANCES];

flat in int instance;

void main(){
    float w = round(1/uvtransforms[instance][0][0]);
    float h = round(1/uvtransforms[instance][1][1]);
    float x = round(uvtransforms[instance][2][0]/uvtransforms[instance][0][0]);
    float y = round(uvtransforms[instance][2][1]/uvtransforms[instance][1][1]);
    float offset = y*w + x;
    frag_colour = texture(u_texture, vec3(spriteCoord,offset));
    if (frag_colour.a < 0.1) { discard; }
}


/*

0.09 0.18

0.2 0.4
*/

/*

		w := 1 / transform.At(0, 0)
		// h := 1 / transform.At(1, 1)

		xtransform := transform.At(0, 2)
		ytransform := transform.At(1, 2)

		x := math32.Round(xtransform / transform.At(0, 0))
		y := math32.Round(ytransform / transform.At(1, 1))
		// if xtransform == 0 {
		// 	x = 0
		// }
		// if ytransform == 0 {
		// 	y = 0
		// }

		offset := y*w + x


*/