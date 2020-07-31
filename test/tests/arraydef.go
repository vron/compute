package kernel

import (
	"encoding/binary"
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(std430, set = 0, binding = 0) buffer Out {
	uint[] dout; // or lie this...
};

struct A {
	int a[1][2];
};

struct B {
	vec3[2] b;
};

layout(std430, set = 0, binding = 0) buffer In {
	uint a1;
	uint[1] a2;
	uint[2] a3;
	uint[1][2] a4;
	uint[1][2] a5[3];
	A[2] a6;
	B[3] b1;
};

// TODO: actually test to access this one
layout(std430, set = 0, binding = 0) buffer In2 {
	B[][2] dInVec; // or lie this...
};

shared uint al1[2];
shared uint[2] al2;

shared uint mlt1[1][2][3];
shared uint[1][2][3] mlt2;
shared uint[3] mlt3[1][2];


void main() {
	uint i = -1;

	/* First test all the shared variables */

	al1[1] = 1;
	dout[++i] = al1[1];

	al2[0] = 1;
	dout[++i] = al2[0];

	mlt1[0][1][2] = 1;
	dout[++i] = mlt1[0][1][2];

	mlt2[0][1][2] = 1;
	dout[++i] = mlt2[0][1][2];

	mlt3[0][1][2] = 1;
	dout[++i] = mlt3[0][1][2];

	/* Then test that we correctly map from input */
	dout[++i] = a1;
	dout[++i] = a2[0];
	dout[++i] = a3[1];
	dout[++i] = a4[0][1];
	dout[++i] = a5[2][0][1];
	dout[++i] = a6[1].a[0][1];


	++i;
	dout[i] = i;
}
`

func TestShader(t *testing.T) {
	nox := 5 + 6
	d := Data{
		Dout: make([]byte, (1+nox)*4),
	}
	d.A1 = 1
	d.A2[0] = 1
	d.A3[1] = 1
	d.A4[0][1] = 1
	d.A5[2][0][1] = 1
	d.A6[1].A[0][1] = 1

	ensureRun(t, 1, d, 1, 1, 1)

	for i := 0; i < nox; i++ {
		v := binary.LittleEndian.Uint32(d.Dout[i*4:])
		if v != 1 {
			t.Error(i, "should be 1 but found", v)
		}
	}
	v := int(binary.LittleEndian.Uint32(d.Dout[nox*4:]))
	if v != nox {
		t.Error("wrong number of chcs: ", v, nox)
	}
}
