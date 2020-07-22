package kernel

import (
	"reflect"
	"testing"
	"unsafe"
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

func TestShader(t *testing.T) {
	// create the input data
	din := make([]byte, 2*2*2*4)
	dout := make([]byte, 2*2*2*4)
	d := Data{Din: din, Dout: dout}
	ensureRun(t, -1, d, 1, 1, 1)

	out := unsafeToFloat(dout)
	if len(out) != 8 {
		t.Error("bad conversion")
	}
	for i := range out {
		if out[i] != 200 {
			t.Error(i, out[i], "!=", 200.0)
		}
	}
}

func TestShader3(t *testing.T) {
	// create the input data
	din := make([]byte, 2*2*2*4)
	dout := make([]byte, 2*2*2*4)
	d := Data{Din: din, Dout: dout}
	ensureRun(t, 3, d, 1, 1, 1)

	out := unsafeToFloat(dout)
	if len(out) != 8 {
		t.Error("bad conversion")
	}
	for i := range out {
		if out[i] != 200 {
			t.Error(i, out[i], "!=", 200.0)
		}
	}
}
func unsafeToFloat(raw []byte) []float32 {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len /= 4
	header.Cap /= 4
	data := *(*[]float32)(unsafe.Pointer(&header))
	return data
}
