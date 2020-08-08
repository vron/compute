#pragma once
#include "common.hpp"
#include "vec.hpp"
#include "vec_make.hpp"
#include <cstdint>

struct image2Drgba32f {
  float *data;
  char padd[20];
  int32_t width;
};

void always_inline imageStore(image2Drgba32f image, ivec2 P, vec4 data) {
  int32_t index = 4 * P[0] + 4 * P[1] * image.width;
  image.data[index + 0] = data[0];
  image.data[index + 1] = data[1];
  image.data[index + 2] = data[2];
  image.data[index + 3] = data[3];
}

vec4 always_inline imageLoad(image2Drgba32f image, ivec2 P) {
  int32_t index = 4 * P[0] + 4 * P[1] * image.width;
  return make_vec4(image.data[index + 0], image.data[index + 1],
                   image.data[index + 2], image.data[index + 3]);
}