package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct element {
	vec3 notused;
};

layout(std430, set = 0, binding = 1) buffer To {
	element array[];
};

void main() {}
`

func TestShader(t *testing.T) {
	d := DataRaw{
		Array: make([]byte, 100*Element{}.Sizeof()+4),
	}
	k, _ := New(-1, -1)
	defer k.Free()
	if nil == k.DispatchRaw(d, 1, 1, 1) {
		t.Error("expected error")
	}
}
