package kernel

import (
	"math"
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(std430) buffer Out1 {
	uint[] u;
};

layout(std430) buffer Out2 {
	float[] f;
};

void main() {
	u[0] = floatBitsToUint(f[0]);
	f[1] = uintBitsToFloat(u[1]);
}
`

func TestShader(t *testing.T) {
	ensureRun(t, 1, 1, 1, 1,
		func() Data {
			return Data{
				F: []float32{1.1, 1.2, 0},
				U: []uint32{0, 3, 0},
			}
		},
		func(res Data) {
			if math.Float32bits(res.F[1]) != 3 {
				t.Error("expected 3", res.F[1], math.Float32bits(res.F[1]))
			}
			if math.Float32frombits(res.U[0]) != 1.1 {
				t.Error("expected 1.1", res.U[0], math.Float32frombits(res.U[0]))
			}
		})
}
