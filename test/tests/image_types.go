package kernel

import (
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(rgba8) uniform writeonly image2D i8;
layout(rgba32f) uniform image2D i32f; // this one is tested in other test

void main() {
	imageStore(i8, ivec2(0,0), vec4(0, 0.5, 1.0, 0.5));
}
`

func TestShader(t *testing.T) {
	ensureRun(t, 1, 1, 1, 1, func() Data {
		return Data{
			I8: Image2Drgba8{
				Data:  make([]byte, 4*8*8),
				Width: 8,
			},
			I32f: Image2Drgba32f{
				Data:  make([]float32, 4*8*8),
				Width: 8,
			},
		}
	}, func(res Data) {
		if res.I8.Data[0] != 0 {
			t.Error("0 should be 0, ", res.I8.Data[0])
		}
		if res.I8.Data[1] != 127 {
			t.Error("1 should be 127, ", res.I8.Data[1])
		}
		if res.I8.Data[2] != 255 {
			t.Error("2 should be 255, ", res.I8.Data[2])
		}
	})
}
