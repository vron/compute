package kernel

import (
	"reflect"
	"testing"
	"unsafe"
)

var shader = `
#version 450

layout(local_size_x = 2, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 0) buffer Data {
	uint data[];
};

struct mstr {
	int a;
	int b;
};

shared uint shInt;
shared uint shIntArr[10];
shared vec2 shVec;
shared mstr shStr;

void main() {

	atomicAdd(shIntArr[0], 1);
	shIntArr[0]+= 1;
	uint i = gl_LocalInvocationIndex;
	if (i==0) {
		shInt = 0;
	}

	barrier();

	if (i == 1) {
		atomicAdd(shInt, 1);
	}

	barrier();

	data[i*10 + 0] = shInt;
}
`

func TestShader(t *testing.T) {
	data := make([]uint32, 20)
	d := Data{Data: intToByte(data)}
	ensureRun(t, -1, d, 1, 1, 1) // TODO: Testcase to ensure that the shared data is not shared between dispatches

	if data[0] != 1 {
		t.Error("not as expected")
	}
	if data[10] != 1 {
		t.Error("not as expected")
	}
}

func intToByte(raw []uint32) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= 4
	header.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
