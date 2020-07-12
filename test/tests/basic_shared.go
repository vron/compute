package kernel

import (
	"reflect"
	"testing"
	"unsafe"
)

var shader = `
#version 450

layout (local_size_x = 2, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 0) buffer Out {
	int values[];
};

shared int shared_data[16];

void main() {
	// one thread writes
    if (gl_LocalInvocationID.x == 0) {
        shared_data[0] = 3;
	}

	barrier();

	// one reads
    if (gl_LocalInvocationID.x == 1) {
        values[0] = shared_data[0];
    }
}
`

func TestShader(t *testing.T) {
	// create the input data
	data := make([]int, 4)
	d := Data{
		Values: intToByte(data),
	}

	ensureRun(t, -1, d, 1, 1, 1)

	if data[0] != 3 {
		t.Error("data not as expected: ", data)
	}
}

func intToByte(raw []int) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= 4
	header.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
