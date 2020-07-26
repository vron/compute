package kernel

import (
	"encoding/binary"
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
	d := Data{
		Data: make([]byte, 8),
	}
	ensureRun(t, 1, d, 1, 1, 1)

	f1 := binary.LittleEndian.Uint32(d.Data)
	if f1 != 3 {
		t.Error("0 should be 3, ", f1)
	}
}
