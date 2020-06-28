// This file will be generated as part of the build, this is an example only

#define _cpt_WG_SIZE_X 8
#define _cpt_WG_SIZE_Y 8
#define _cpt_WG_SIZE_Z 1

struct kernel_comp {
  uvec3 gl_GlobalInvocationID;
  uvec3 gl_LocalInvocationID;
  image2DRGBA32F image;
  void main() {
    vec4_scalar pixel = make_vec4(0.f, 0.f, 0.f, 1.f);
    uvec2 pixel_coords;
    pixel_coords = (gl_GlobalInvocationID).sel(_ind_X, _ind_Y);
    pixel_coords *= 8;
    pixel_coords += (gl_LocalInvocationID).sel(_ind_X, _ind_Y);
    imageStore(image, make_ivec2(pixel_coords), pixel);
  }

  int set_data(kernel_data d) {
	#include "setdata.hpp"
  }
};
