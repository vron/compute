package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct element {
	vec3 notused;
	bool peter;
	float val;
};

layout(std430, set = 0, binding = 0) buffer From {
	int index;
	bool flag;
	vec3 pos;
	vec3 posar[3];
};

layout(std430, set = 0, binding = 1) buffer To {
	element array[];
};

void main() {
	float idx = float(index);
	if (flag) {
		idx += pos.x + pos.y + pos.z;
	}
	array[index].val = idx;
}
`

func TestShader(t *testing.T) {
	d := Data{
		Index: 43,
		Flag:  true,
		Pos:   [3]float32{11, 22, 33},
		Array: make([]byte, 100*Element{}.Stride()+4),
	}

	if run(t, 1, d, 1, 1, 1) == nil {
		t.Error("expected error showing bad array length - got none")
	}
}

// TODO: Add test case to ensure we chec specific length when we have the info available
