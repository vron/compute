package kernel

import (
	"testing"
)

// TODO: Also add a test case with vectors and matrices as part of a slice - to chec that we are picing up
// the data at the correct location...
// Actually! run this same test case twice in order to achieve that!

var shader = `
#version 450

layout(local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

layout(std430) buffer Data {
	mat2 m2;
	mat3 m3;
	mat4 m4;
	float results[]; // expect all of them to be 0!
};

void main() {
	int i = -1;
	// Try to access components
	results[++i] = 8.0f+dot(-2.0f*m2[0], m2[1]); // this should equal 0 if everything is correct.

	// Try matrix*matrix
	mat2 p = m2*m2; // 3 3 6 6
	results[++i] = p[0][0] - 3.0f;
	results[++i] = p[0][1] - 3.0f;
	results[++i] = p[1][0] - 6.0f;
	results[++i] = p[1][1] - 6.0f;

	// try chaning the value in a matrix
	m2[0][0] = 0;
	m2[1][1] = m2[0][0];
	results[++i] = m2[1][1];

	// try multiplying a matrix with a vector
	results[++i] = dot(m3*vec3(-1, -1, 1), vec3(1));
}
`

func TestShader(t *testing.T) {
	ensureRun(t, 1, 1, 1, 1, func() Data {
		return Data{
			Results: make([]float32, 100),
			M2:      &Mat2{Vec2{1, 1}, Vec2{2, 2}},
			M3:      &Mat3{Vec3{1, 1, 1}, Vec3{2, 2, 2}, Vec3{3, 3, 3}},
			M4:      &Mat4{Vec4{1, 1, 1, 1}, Vec4{2, 2, 2, 2}, Vec4{3, 3, 3, 3}, Vec4{4, 4, 4, 4}},
		}
	}, func(res Data) {
		for i := range res.Results {
			if res.Results[i] != 0 {
				t.Error(i, "data not 0:", res.Results[i])
			}
		}
	})
}
