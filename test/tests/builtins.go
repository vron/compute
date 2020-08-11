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

layout(std430) buffer Out3 {
	int[] i;
};

void main() {
	u[0] = floatBitsToUint(f[0]);
	f[1] = uintBitsToFloat(u[1]);
	f[2] = min(1.1, 2.2);
	f[3] = max(1.1, 2.2);
	f[4] = mix(vec2(8,0), vec2(10,0), 0.5).x;
	f[5] = clamp(5.0, 3.0, 4.0);
	u[3] = bitCount(3);
	i[0] = findLSB(0);
	i[1] = findLSB(1);
	i[2] = findLSB(2);
	f[6] = round(vec3(1.3, 1.4, 1.2)).x;
}
`

func TestShader(t *testing.T) {
	ensureRun(t, 1, 1, 1, 1,
		func() Data {
			return Data{
				F: []float32{1.1, 1.2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				U: []uint32{0, 3, 0, 0, 0},
				I: []int32{0, 0, 0},
			}
		},
		func(res Data) {
			if math.Float32bits(res.F[1]) != 3 {
				t.Error("expected 3", res.F[1], math.Float32bits(res.F[1]))
			}
			if math.Float32frombits(res.U[0]) != 1.1 {
				t.Error("expected 1.1", res.U[0], math.Float32frombits(res.U[0]))
			}
			if res.F[2] != 1.1 {
				t.Error("expected 1.1 for min", res.F[2])
			}
			if res.F[3] != 2.2 {
				t.Error("expected 2.2 for max", res.F[3])
			}
			if res.F[4] != 9 { // mix
				t.Error("expected 9 for mix", res.F[4])
			}
			if res.F[5] != 4 { // clamp
				t.Error("expected 4 for clamp", res.F[5])
			}
			if res.U[3] != 2 {
				t.Error("expected 2 for bitcount", res.U[3])
			}
			if res.I[0] != -1 {
				t.Error("expected -1 for bitcount", res.I[0])
			}
			if res.I[1] != 0 {
				t.Error("expected 0 for bitcount", res.I[1])
			}
			if res.I[2] != 1 {
				t.Error("expected 1 for bitcount", res.I[2])
			}
			if res.F[6] != 1 {
				t.Error("expected 1 for round", res.F[6])
			}
		})
}
