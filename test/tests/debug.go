package kernel

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

var shader = `
#version 450

layout(local_size_x = 4, local_size_y = 4, local_size_z = 4) in;

layout(std430) buffer Out {
	int data[];
};

void main() {
	uint index = gl_LocalInvocationIndex;
	int value = int(index);

	if (index % 2 == 0) {
		value *= 2;
	}

	barrier();

	if (index < 1024) {
		value = 1;
	}
	for (int i = 0; i < 4; i++) {
		if (index % 2 == 0) {
			value += 1;
		}

		barrier();
	}

	barrier();

	data[index] = value;
}
`

func TestShader(t *testing.T) {
	for i := 0; i < 10; i++ {
		test(t)
	}
}

func test(t *testing.T) {
	data := make([]int32, 64)
	d := Data{Data: intToByte(data)}
	k, err := New(runtime.GOMAXPROCS(-1), 1024*1024)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer k.Free()

	for i := 0; i < 10; i++ {
		err := k.Dispatch(d, 1, 1, 1) // TODO: 8,8,8
		if err != nil {
			t.Error(err)
		}
	}

	for i := range data {
		ex := int32(1)
		if i%2 == 0 {
			ex = 5
		}
		if data[i] != ex {
			t.Error(i, "expected value: ", ex, "got", data[i])
		}
	}
}

func intToByte(raw []int32) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= 4
	header.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
