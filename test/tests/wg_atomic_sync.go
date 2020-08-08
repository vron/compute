package kernel

import (
	"testing"
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
	ensureRun(t, -1, 2, 3, 4, func() Data {
		values := make([]int32, 144)
		mysum := make([]int32, 2)
		for i := range values {
			values[i] = 1
		}
		return Data{Values: values, Mysum: mysum}
	}, func(res Data) {
		for i := range res.Values {
			if res.Values[i] != 1 {
				t.Error(i, res.Values[i], "!=", 1)
			}
		}
		if res.Mysum[0] != 144 {
			t.Error("sum not 144", res.Mysum)
		}
	})
}
