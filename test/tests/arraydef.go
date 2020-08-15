package kernel

import (
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
	int a1[1][2];
};

struct B {
	uvec3[2] b;
};

layout(std430, set = 0, binding = 0) buffer In {
	uint a1;
	uint[1] a2;
	uint[2] a3;
	uint[1][2] a4;
	uint[1][2] a5[3];
	A[2] a6;
	B[3] b1;
	B b2;
};

layout(std430, set = 0, binding = 0) buffer In2 {
	B[][2] dInVec;
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
	dout[++i] = b1[1].b[0].y;
	dout[++i] = b2.b[0].y;


	/* test some complicated inputs handling */
	dout[++i] = dInVec[1][0].b[0].x;

	++i;
	dout[i] = i;
}
`

func TestShader(t *testing.T) {
	nox := 5 + 8 + 1

	ensureRun(t, 1, 1, 1, 1,
		func() Data {
			ai := uint32(1)
			return Data{
				Dout:   make([]uint32, (1 + nox)),
				A1:     &ai,
				A2:     &[1]uint32{1},
				A3:     &[2]uint32{0, 1},
				A4:     &[1][2]uint32{{0, 1}},
				A5:     &[3][1][2]uint32{{}, {}, {{0, 1}}},
				A6:     &[2]A{{}, {A: [1][2]int32{{0, 1}}}},
				B1:     &[3]B{{}, {B: [2]Uvec3{{0, 1, 0}}}},
				B2:     &B{B: [2]Uvec3{{0, 1, 0}}},
				DInVec: [][2]B{{}, {{B: [2]Uvec3{{1, 0, 0}}}}},
			}
		},
		func(res Data) {
			for i := 0; i < nox; i++ {
				v := res.Dout[i]
				if v != 1 {
					t.Error(i, "should be 1 but found", v)
				}
			}
			v := res.Dout[nox]
			if v != uint32(nox) {
				t.Error("wrong number of chcs: ", v, nox)
			}
		})
}
