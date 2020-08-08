package kernel

import (
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 0) buffer In {
	uint data[];
};

struct A {int a;};
struct B {int a;};
struct C {int a;};
struct D {int a;};
struct E {int a;};
struct F {int a;};
struct H {int a;};
struct I {int a;};
struct J {int a;};
struct L {int a;};
struct N {int a;};
struct O {int a;};
struct P {int a;};
struct Q {int a;};

struct M {
	int a;
};

struct G {
	M g;	
};


shared G m;

void main() {
	m.g.a = 3;
	data[0] = m.g.a;
}
`

func TestShader(t *testing.T) {
	ensureRun(t, 1, 1, 1, 1, func() Data {
		return Data{
			Data: make([]uint32, 2),
		}
	}, func(res Data) {
		if res.Data[0] != 3 {
			t.Error("0 should be 3, ", res.Data[0])
		}
	})
}
