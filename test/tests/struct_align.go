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
	ensureRun(t, -1, 1, 1, 1, func() Data {
		i := int32(43)
		b := True
		return Data{
			Index: &i,
			Flag:  &b,
			Pos:   &Vec3{11, 22, 33}, // sum = 66,
			Array: make([]Element, 100),
		}
	}, func(res Data) {
		if res.Array[43].Val != 66+43 {
			for i := 0; i < 100; i++ {
				t.Log(res.Array[i])
			}
			t.Log(res.Array[43].Val)
			t.Error("output not as expected")
		}
	})
}
