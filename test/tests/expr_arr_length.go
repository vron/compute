package kernel

import (
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 0) buffer In {
	int d1[1<<8];
	int d2[1+2];
	int d3[1-(-1)];
	int d4[4*4];
	int d5[255>>7];
	int d6[11/2];
};

void main() {
	d1[255] = 5;
}
`

func TestShader(t *testing.T) {
	ensureRun(t, 1, 1, 1, 1, func() Data {
		return Data{
			D1: &[1 << 8]int32{},
			D2: &[1 + 2]int32{},
			D3: &[1 - (-1)]int32{},
			D4: &[4 * 4]int32{},
			D5: &[255 >> 7]int32{},
			D6: &[11 / 2]int32{},
		}
	}, func(res Data) {
		if res.D1[255] != 5 {
			t.Error("should be 5", res.D1[255])
		}
	})
}
