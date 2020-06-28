package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct element {
	vec3 notused;
	bool peter;
	int val;
};

struct notused {
	vec3 aa[3];
	bool bb[2];
	float cc;
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
	int idx = int(index);
	if (flag) {
		idx += int(pos.x + pos.y + pos.z);
	}
	array[index].val = idx;
}
`

func TestShader(t *testing.T) {
	p := [3]float32{11, 22, 33} // sum = 66
	d := Data{
		Index: 43,
		Flag:  true,
		Pos:   p,
		Array: make([]byte, 100*Element{}.Stride()),
	}

	ensureRun(t, 1, d, 1, 1, 1)

	// Decode the given element and verify
	el := &Element{}
	el.Decode(d.Array[el.Stride()*43:])
	if el.Val != 66+43 {
		for i := 0; i < 100; i++ {
			el := &Element{}
			el.Decode(d.Array[el.Stride()*i:])
			t.Log(el)
		}
		t.Log(d.Array[el.Stride()*43:])
		t.Error("output not as expected", el)
	}
}
