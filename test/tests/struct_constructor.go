package kernel

import (
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

struct cc {
	uint i;
};

struct element {
	vec3 a;
	bool b;
	cc c;
};

layout(std430) buffer To {
	element array[];
};

void main() {
	float a = 7.7;
	float b = 1.0;
	array[0] = element(vec3(1.0,1.0,1.0), true, cc(2));
}
`

func TestShader(t *testing.T) {
	ensureRun(t, 1, 1, 1, 1,
		func() Data {
			return Data{
				Array: []Element{{}},
			}
		},
		func(res Data) {
			e1 := res.Array[0]
			if e1.A != (Vec3{1, 1, 1}) {
				t.Error("bad vec", e1.A)
			}
			if !e1.B.B {
				t.Error("expected true", e1.B)
			}
			if e1.C.I != 2 {
				t.Error("expected 2", e1.C.I)
			}
		})
}
