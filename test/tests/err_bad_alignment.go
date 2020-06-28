package kernel

import (
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 8, local_size_y = 8, local_size_z = 1) in;
layout(rgba32f, binding = 0) uniform image2D img;

void main() {
  vec4 pixel = vec4(0.0f, 0.0f, 0.0f, 1.0f);
  ivec2 pixel_coords = ivec2(gl_GlobalInvocationID.xy)*8;
  pixel_coords += ivec2(gl_LocalInvocationID.xy);
  if(pixel_coords.x%2 == 0) {
    pixel = vec4(1.0f, 1.0f, 1.0f, 1.0f);
  } 
  imageStore(img, pixel_coords, pixel);
}
`

func TestShader(t *testing.T) {
	// create the input data
	nop := 2
	img := make([]byte, 4*nop*nop*8*8*4+1)
	d := Data{
		Img: Image2Drgba32f{
			// offset the data by 1 => should create alignment error
			Data:  img[1:],
			Width: int32(nop * 8),
		},
	}

	if run(t, 1, d, nop, nop, 1) == nil {
		t.Error("expected alignment error, got no error")
	}
}
