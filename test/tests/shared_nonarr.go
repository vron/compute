package kernel

import (
	"testing"
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

	ensureRun(t, -1, 1, 1, 1, func() Data {
		return Data{
			Data: make([]uint32, 20),
		}
	}, func(res Data) {
		if res.Data[0] != 1 {
			t.Error("not as expected")
		}
		if res.Data[10] != 1 {
			t.Error("not as expected")
		}
	})
}
