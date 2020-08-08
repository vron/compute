package kernel

import (
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
	ensureRun(t, 1, 1, 1, 1, func() Data {
		return Data{
			Data: make([]int32, 2),
		}
	}, func(res Data) {
		if res.Data[0] != 1 {
			t.Error("0 should be 1, ", res.Data[0])
		}
		if res.Data[1] != 2 {
			t.Error("1 should be 2, ", res.Data[1])
		}
	})
}
