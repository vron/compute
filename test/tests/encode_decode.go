package kernel

import (
	"reflect"
	"testing"
)

var shader = `
#version 450

struct s1 {
	bool a;
};
struct s2 {
	s1 s;
};
struct s3 {
	s1 s[1];
};
struct s4 {
	int p;
	ivec2 d;
};
struct s5 {
	s4 a;
	bool b;
	s4 c;
};
struct s6 {
	mat3 m;
};

struct outer {
	bool a;
	s1 s1;
	mat2 b;
	s2 s2;
	s3 s3;
	vec2 c[2];
	s4 s4;
	s5 s5;
	float d;
	mat2 e[3];
	s6 s6;
};

layout(std430) buffer Dummy {
	outer l;
	outer data[];
};

layout (local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

void main() {
	uint index = gl_WorkGroupID.x;

	outer d;
	d.a = true;
	d.s1.a = false;
	d.b = mat2(vec2(1,2), vec2(3,4));
	d.s2.s.a = true;
	d.s3.s[0].a = true;
	d.c[0] = vec2(5,6);
	d.s4.p = 12;
	d.s4.d = ivec2(99,89);
	d.s5.a.p = 121;
	d.s5.a.d = ivec2(44,55);
	d.s5.b = true;
	d.s5.c.p = 122;
	d.s5.c.d = ivec2(22,2);
	d.d = 3.3f;
	d.e[0] = mat2(-1, -2, -3, -4);
	d.e[1] = mat2(-5, -6, -7, -8);
	d.e[2] = mat2(-9, -10, -11, -12);
	d.s6.m = mat3(vec3(10,11,12), vec3(13, 14, 15), vec3(16,17,18));

	if (index%2 ==0 ) {
		data[index] = d;
	} else {
		data[index] = l;
	}
}
`

// Todo, change the above to also copy over one struct from the input to the output

var ref Outer = Outer{
	A:  true,
	S1: S1{A: false},
	B:  Mat2{Vec2{1, 2}, Vec2{3, 4}},
	S2: S2{S1{A: true}},
	S3: S3{[1]S1{S1{A: true}}},
	C:  [2]Vec2{{5, 6}, {0, 0}},
	S4: S4{P: 12, D: Ivec2{99, 89}},
	S5: S5{A: S4{P: 121, D: Ivec2{44, 55}}, B: true, C: S4{P: 122, D: Ivec2{22, 2}}},
	D:  3.3,
	E:  [3]Mat2{Mat2{Vec2{-1, -2}, Vec2{-3, -4}}, Mat2{Vec2{-5, -6}, Vec2{-7, -8}}, Mat2{Vec2{-9, -10}, Vec2{-11, -12}}},
	S6: S6{M: Mat3{Vec3{10, 11, 12}, Vec3{13, 14, 15}, Vec3{16, 17, 18}}},
}

func TestEncodeDecode(t *testing.T) {
	// create a struct that we will encode and decode to ensure it..
	buf := make([]byte, ref.Sizeof())
	ref.Encode(buf)
	var d Outer
	(&d).Decode(buf)

	if !reflect.DeepEqual(ref, d) {
		t.Error("encoded and decoded struct data not equal...")
		t.Log(ref)
		t.Log(d)
	}
}

func TestDecode(t *testing.T) {
	ensureRun(t, -1, 2, 1, 1, func() Data {
		r := ref
		return Data{
			Data: make([]Outer, 2),
			L:    &r,
		}
	}, func(res Data) {
		for i := 0; i < 2; i++ {
			if !reflect.DeepEqual(ref, res.Data[i]) {
				t.Error("encoded and decoded struct data not equal...", i)
				t.Log(ref)
				t.Log(res.Data[i])
			}
		}
	})
}
