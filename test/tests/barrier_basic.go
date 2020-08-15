package kernel

import (
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 2, local_size_y = 2, local_size_z = 2) in;

layout(std430, set = 0, binding = 0) buffer In {
	float din[];
};

layout(std430, set = 0, binding = 0) buffer Out {
	float dout[];
};

void main() {
	uint index = gl_LocalInvocationID.x + gl_LocalInvocationID.y*2 + gl_LocalInvocationID.z*2*2;
	float d = din[index] + 1.0f;
	float val = 0.0f;
	for(int i = 0; i < 100; i++) {
		val += d;
	}
	dout[index] = val + din[index];

	// Wait for all in this wg to finish before continuing.
	barrier();
	
	// now read from another index, ensuing no overflow. If the barrier is not woring correctly
	// we should not get the same value out in all indices..
	index = (index + 33) % (2*2*2);
	dout[index] *= 2;
}
`

func test(t *testing.T, nt int) {
	ensureRun(t, nt, 1, 1, 1, func() Data {
		return Data{
			Din:  make([]float32, 8),
			Dout: make([]float32, 8),
		}
	}, func(res Data) {
		for i := range res.Dout {
			if res.Dout[i] != 200 {
				t.Error(i, res.Dout[i], "!=", 200.0)
			}
		}
	})
}

func TestShader(t *testing.T) {
	test(t, 1)
}

func TestShader3(t *testing.T) {
	test(t, 3)
}
