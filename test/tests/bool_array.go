package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct notused {
	vec3 aa[3];
	bool bb[2];
	float cc;
};

layout(std430, set = 0, binding = 1) buffer To {
	notused array[];
};

void main() {
	if (array[0].bb[0]) {
		array[0].cc = 33;
	}
	if (!array[0].bb[1]) {
		array[0].cc = 33;
	}
	array[1].bb[0] = false;
	array[1].bb[1] = true;
}
`

func TestShader(t *testing.T) {
	ensureRun(t, -1, 1, 1, 1, func() Data {
		return Data{
			Array: []Notused{{
				Bb: [2]Bool{{B: false}, {B: true}},
			}, {
				Bb: [2]Bool{{B: false}, {B: true}},
			}}}
	}, func(res Data) {
		if res.Array[0].Cc != 0.0 {
			t.Error("expected 0")
		}
		if res.Array[0].Bb[0].B {
			t.Error("expected false 1")
		}
		if !res.Array[0].Bb[1].B {
			t.Error("expected true 2")
		}
		if res.Array[1].Bb[0].B {
			t.Error("expected false 3")
		}
		if !res.Array[1].Bb[1].B {
			t.Error("expected true 4")
		}
	})
}
