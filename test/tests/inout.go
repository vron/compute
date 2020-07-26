package kernel

import (
	"encoding/binary"
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 0) buffer In {
	int data[];
};

void f(int a, inout int b, out int c) {
	b = a;
	c = a;
}

void main() {
	int a = 1;
	int b = 2;
	int c = 2;
	f(a, b, c);
  	data[0] = b; // should now be 1
  	data[1] = b; // should now be 1
}
`

func TestShader(t *testing.T) {
	data := make([]byte, 8)
	d := Data{
		Data: data,
	}

	ensureRun(t, 1, d, 1, 1, 1)

	f1 := binary.LittleEndian.Uint32(data)
	f2 := binary.LittleEndian.Uint32(data[4:])
	if f1 != 1 {
		t.Error("0 should be 1, ", f1)
	}
	if f2 != 1 {
		t.Error("0 should be 1, ", f1)
	}
}
