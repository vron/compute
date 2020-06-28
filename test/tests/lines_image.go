package test

import "testing"

var shader = `
#version 450

layout(local_size_x = 8, local_size_y = 8, local_size_z = 1) in;
layout(rgba32f, binding = 0) uniform image2D img;

void main() {
  vec4 pixel = vec4(0.0f, 0.0f, 0.0f, 1.0f);

  ivec2 pixel_coords;
  pixel_coords.x = int(gl_GlobalInvocationID.x);
  pixel_coords.y = int(gl_GlobalInvocationID.y);
  pixel_coords *= 8;
  pixel_coords.x += int(gl_LocalInvocationID.x);
  pixel_coords.y += int(gl_LocalInvocationID.y);
  if(pixel_coords.x%2 == 0) {
    pixel = vec4(1.0f, 1.0f, 1.0f, 1.0f);
  } 
  imageStore(img, pixel_coords, pixel);
}
`

func run(t *testing.T, nt int, d Data, numx, numy, numz int) {
	k, err := New(nt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer k.Free()
	err = k.Dispatch(d, numx, numy, numz)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestShader(t *testing.T) {
	// create the input data
	nop := 2
	img := make([]float32, nop*nop*8*8*4)
	d := Data{
		imgData:  &(img[0]),
		imgWidth: uint(nop * 8),
	}

	run(t, 1, d, nop, nop, 1)
	if img[0] != 1 ||
		img[1] != 1 || img[2] != 1 || img[3] != 1 || img[4] != 0 || img[5] != 0 || img[6] != 0 || img[7] != 1 {
		t.Log(img[0], img[1], img[2], img[3], img[4], img[5], img[6], img[7])
		t.Error("output not as expected")
	}
}
