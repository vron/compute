package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 1) buffer To {
	vec2 array[10];
};

void main() {}
`

func TestShader(t *testing.T) {
	d := DataRaw{
		Array: make([]byte, 9*Vec2{}.Sizeof()),
	}
	k, _ := New(-1, -1)
	defer k.Free()
	if nil == k.DispatchRaw(d, 1, 1, 1) {
		t.Error("expected error")
	}
}
