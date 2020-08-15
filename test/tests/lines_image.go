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
  ivec2 pixel_coords = ivec2(gl_GlobalInvocationID.xy);
  if(pixel_coords.x%2 == 0) {
    pixel = vec4(1.0f, 1.0f, 1.0f, 1.0f);
  } 
  imageStore(img, pixel_coords, pixel);
}
`

func TestShader(t *testing.T) {
	// create the input data
	nop := 2
	ensureRun(t, 1, nop, nop, 1, func() Data {
		return Data{
			Img: Image2Drgba32f{
				Data:  make([]float32, nop*nop*8*8*4),
				Width: int32(nop * 8),
			},
		}
	}, func(res Data) {
		img := res.Img.Data
		if img[0] != 1 ||
			img[1] != 1 || img[2] != 1 || img[3] != 1 || img[4] != 0 || img[5] != 0 || img[6] != 0 || img[7] != 1 {
			t.Log(img[0], img[1], img[2], img[3], img[4], img[5], img[6], img[7], img[8], img[9], img[10], img[11])
			t.Error("output not as expected")
		}
	})
}
