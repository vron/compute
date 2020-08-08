package kernel

import (
	"testing"
)

var shader = `
#version 450

layout (local_size_x = 16, local_size_y = 1, local_size_z = 1) in;

layout(rgba32f) uniform image2D inImage;
layout(rgba32f) uniform image2D outImage;

shared vec4 shared_data[16];

void main() {
    ivec2 base = ivec2(gl_WorkGroupID.xy * gl_WorkGroupSize.xy);
	ivec2 my_index = base + ivec2(gl_LocalInvocationID.x,0);
	
    if (gl_LocalInvocationID.x == 0) {
        for (int i = 0; i < 16; i++) {
            shared_data[i] = imageLoad(inImage, base + ivec2(i,0));
        }
    }

    barrier();

    imageStore(outImage, my_index, shared_data[gl_LocalInvocationID.x]);
}
`

func TestShader(t *testing.T) {
	ensureRun(t, 1, 4, 4, 1, func() Data {
		img1 := make([]float32, 4*16*4*4)
		for i := range img1 {
			img1[i] = float32(i)
		}
		img2 := make([]float32, 4*16*4*4)
		return Data{
			InImage: Image2Drgba32f{
				Data:  (img1),
				Width: 16 * 4,
			},
			OutImage: Image2Drgba32f{
				Data:  (img2),
				Width: 16 * 4,
			},
		}
	}, func(res Data) {
		img1, img2 := res.InImage.Data, res.OutImage.Data
		for i := range img1 {
			if float32(i) != img2[i] {
				t.Error("not equal at", i, img1[i], "!=", img2[i])
			}
		}
	})
}
