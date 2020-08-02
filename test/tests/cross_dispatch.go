package kernel

import (
	"runtime"
	"testing"
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
	d := Data{Data: make([]int32, 64)}
	k, err := New(runtime.GOMAXPROCS(-1), 1024*1024)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer k.Free()

	for i := 0; i < 10; i++ {
		err := k.Dispatch(d, 8, 8, 8)
		if err != nil {
			t.Error(err)
		}
	}

	for i := range d.Data {
		ex := int32(1)
		if i%2 == 0 {
			ex = 5
		}
		if d.Data[i] != ex {
			t.Error(i, "expected value: ", ex, "got", d.Data[i])
		}
	}
}
