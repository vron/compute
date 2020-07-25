package kernel

import (
	"encoding/binary"
	"testing"
)

var shader = `
#version 450

#define VAL 1
#define MUL(x) (x*2)

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 0) buffer In {
	int data[];
};

void main() {
  data[0] = VAL;
  data[1] = MUL(VAL);
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
	if f2 != 2 {
		t.Error("1 should be 2, ", f2)
	}
}
