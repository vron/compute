package kernel

import (
	"reflect"
	"testing"
	"unsafe"
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
	data := make([]float32, 100)
	// first columns of 1, 2nd of 2 etc.
	d := Data{
		Results: floatToByte(data),
		M2:      [4]float32{1, 1, 2, 2},
		M3:      [9]float32{1, 1, 1, 2, 2, 2, 3, 3, 3},
		M4:      [16]float32{1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4},
	}
	ensureRun(t, 1, d, 1, 1, 1)
	for i := range data {
		if data[i] != 0 {
			t.Error(i, "data not 0:", data[i])
		}
	}
}

func floatToByte(raw []float32) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= 4
	header.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
