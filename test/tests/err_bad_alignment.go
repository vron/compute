package kernel

import (
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 8, local_size_y = 8, local_size_z = 1) in;
layout(rgba32f, binding = 0) uniform image2D img;
layout(std430, set = 0, binding = 0) buffer In {
	float din[];
};

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
	d := Data{
		Img: Image2Drgba32f{
			Data:  make([]float32, 4*nop*nop*8*8+1),
			Width: int32(nop * 8),
		},
		Din: make([]float32, 10),
	}
	dr := encodeData(d)
	dr.Din = dr.Din[1:] // should create alignemtn error

	k, _ := New(-1, -1)
	defer k.Free()
	if nil == k.DispatchRaw(dr, nop, nop, 1) {
		t.Error("expected error")
	}
}
