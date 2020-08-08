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
	ensureRun(t, 1, 1, 1, 1, func() Data {
		ft := int32(43)
		tb := True
		return Data{
			Index: &ft,
			Flag:  &tb,
			Pos:   &Vec3{11, 22, 33},
			Array: make([]Element, 100),
		}
	}, func(res Data) {
		if res.Array[43].Val != 109 {
			t.Error("output not as expected")
		}
	})
}
