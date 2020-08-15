package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 64, local_size_y = 1, local_size_z = 1) in;

layout(std430) buffer In {
	uint data;
};

void main() {
	data = 0;
	do {
		data++;
	} while (data < 10);
}
`

func TestShader(t *testing.T) {
	// this tests so that casts in glsl are translated to c types
	ensureRun(t, 1, 1, 1, 1, func() Data {
		var a uint32
		return Data{
			Data: &a,
		}
	}, func(res Data) {
		if *res.Data != 10 {
			t.Error("should be 10 got", *res.Data)
		}
	})
}
