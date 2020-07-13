package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 64, local_size_y = 1, local_size_z = 1) in;

struct triangle {
	vec2 vertices[3];
};

struct polygon {
	triangle triangles[2];
};

layout(std430) buffer In {
	polygon polygons[];
};

void main() {
}
`

func TestShader(t *testing.T) {
	// This test simply tests that we manage to generate the order of dedinitions correctly in C,
	// so it might seam empty but do not delete it since it tests that the c code actually compliles
}
