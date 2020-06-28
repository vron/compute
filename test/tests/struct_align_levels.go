package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct element {
	vec3 notused;
	bool peter;
	float val;
	vec2 arr[5]; // so the size of element embedds the full array here! (Note: variable length not supported here)
};

layout(std430, set = 0, binding = 0) buffer From {
	int index;
	bool flag;
	vec3 pos;
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
		Array: make([]byte, 100*Element{}.Stride()),
	}

	ensureRun(t, 1, d, 1, 1, 1)

	// Decode the given element and verify
	el := &Element{}
	el.Decode(d.Array[el.Stride()*43:])
	if el.Val != 109 {
		t.Error("output not as expected")
	}
}
