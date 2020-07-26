package kernel

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

var shader = `
#version 450

layout(local_size_x = 2, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 0) buffer Data {
	uint da[];
	uint db[];
};

void main() {
	uint i = gl_WorkGroupID.x;

	// test atomics on global shared
	if (i %2 == 0) {
		for (int j = 0; atomicAdd(da[1],0) == 0; j++) {};
		atomicExchange(da[0], da[2]);
	} else  {
		atomicExchange(da[2], 2);
		atomicExchange(da[1], 1); 
	}
}
`

func TestShader(t *testing.T) {
	// run multiple times to try to fins scheduling problems
	for i := 0; i < 1000; i++ {
		fmt.Println(i)
		runTest(t)
	}
}

func runTest(t *testing.T) {
	da := make([]uint32, 2)
	db := make([]uint32, 2)
	d := Data{Da: intToByte(db), Db: intToByte(db)}

	ensureRun(t, -1, d, 200, 1, 1)

	if da[0] != 2 {
		t.Error("expected a 2: ", da[0])
	}
}

func intToByte(raw []uint32) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= 4
	header.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
