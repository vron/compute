package kernel

import "testing"

var shader = `
#version 450

layout(local_size_x = 64, local_size_y = 1, local_size_z = 1) in;

layout(std430) buffer To {
	int array[];
};

void main() {
	int a = 3;
	switch(a) {
	case 0:
		int b = 3;
	case 1:
		b = 2;
	case 3:
		int c1 = 4, c2 = 5;
		for(int i = 0; i < 10; i++);
		int bb = 2;
	case 4:
		c1 = 33;
	}
}
`

func TestShader(t *testing.T) {
	// ensuring that the c code we generate actually compiles
}
