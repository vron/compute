package kernel

import (
	"reflect"
	"testing"
	"unsafe"
)

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 2, local_size_z = 3) in;

layout(std430, set = 0, binding = 0) buffer In {
	int mysum[];
	int values[];
};

uint prod(uvec3 a) {
	return a.x*a.y*a.z;
}

uint id(uvec3 size, uvec3 index) {
	return index.x + index.y*size.x + index.z*size.x*size.y;
}

void main() {
	// each call adds one value
	uint wg_size = 1*2*3;
	uint gid = id(gl_NumWorkGroups, gl_WorkGroupID);
	uint index = gid*wg_size + gl_LocalInvocationIndex;
	atomicAdd(mysum[0], values[index]);
}
`

func TestShader(t *testing.T) {
	// create the input data
	values := make([]int32, 144)
	mysum := make([]int32, 2)
	for i := range values {
		values[i] = 1
	}
	d := Data{Values: intToByte(values), Mysum: intToByte(mysum)}
	ensureRun(t, -1, d, 2, 3, 4) // 6*4 = 24 * 6 = 144

	for i := range values {
		if values[i] != 1 {
			t.Error(i, values[i], "!=", 1)
		}
	}
	if mysum[0] != 144 {
		t.Error("sum not 144", mysum)
	}
}

func intToByte(raw []int32) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= 4
	header.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
