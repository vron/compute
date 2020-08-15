#pragma once
#include "common.hpp"
#include "vec.hpp"
#include "vec_make.hpp"
#include <cstdint>

// ref https://www.khronos.org/opengl/wiki/Image_Format
// ref https://www.khronos.org/opengl/wiki/Normalized_Integer

struct image2Drgba32f {
  float *data;
  char padd[20];
  int32_t width;
};

struct image2Drgba8 {
  uint8_t *data;
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

void always_inline imageStore(image2Drgba8 image, ivec2 P, vec4 data) {
  int32_t index = 4 * P[0] + 4 * P[1] * image.width;
  image.data[index + 0] = uint8_t(data[0] * 255.0f);
  image.data[index + 1] = uint8_t(data[1] * 255.0f);
  image.data[index + 2] = uint8_t(data[2] * 255.0f);
  image.data[index + 3] = uint8_t(data[3] * 255.0f);
}

vec4 always_inline imageLoad(image2Drgba8 image, ivec2 P) {
  int32_t index = 4 * P[0] + 4 * P[1] * image.width;
  return make_vec4((float)image.data[index + 0] / 255.0f, (float)image.data[index + 1]/ 255.0f,
                   (float)image.data[index + 2]/ 255.0f, (float)image.data[index + 3]/ 255.0f);
}