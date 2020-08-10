#pragma once
#include <algorithm>
#include <cmath>
#include <cstdint>
#include <cstring>

/* Bit twiddling functions */

uint32_t floatBitsToUint(float f) {
  uint32_t u;
  std::memcpy(&u, &f, sizeof f);
  return u;
}

float uintBitsToFloat(uint32_t u) {
  float f;
  std::memcpy(&f, &u, sizeof f);
  return f;
}

/* min-max functions */
float min(float x, float y) { return std::min(x, y); }
vec2 min(vec2 x, vec2 y) {
  return make_vec2(std::min(x[0], y[0]), std::min(x[1], y[1]));
};
vec3 min(vec3 x, vec3 y) {
  return make_vec3(std::min(x[0], y[0]), std::min(x[1], y[1]),
                   std::min(x[2], y[2]));
};
vec4 min(vec4 x, vec4 y) {
  return make_vec4(std::min(x[0], y[0]), std::min(x[1], y[1]),
                   std::min(x[2], y[2]), std::min(x[3], y[3]));
};
vec2 min(vec2 x, float y) {
  return make_vec2(std::min(x[0], y), std::min(x[1], y));
};
vec3 min(vec3 x, float y) {
  return make_vec3(std::min(x[0], y), std::min(x[1], y), std::min(x[2], y));
};
vec4 min(vec4 x, float y) {
  return make_vec4(std::min(x[0], y), std::min(x[1], y), std::min(x[2], y),
                   std::min(x[3], y));
};
int32_t min(int32_t x, int32_t y) { return std::min(x, y); }
ivec2 min(ivec2 x, ivec2 y) {
  return make_ivec2(std::min(x[0], y[0]), std::min(x[1], y[1]));
};
ivec3 min(ivec3 x, ivec3 y) {
  return make_ivec3(std::min(x[0], y[0]), std::min(x[1], y[1]),
                    std::min(x[2], y[2]));
};
ivec4 min(ivec4 x, ivec4 y) {
  return make_ivec4(std::min(x[0], y[0]), std::min(x[1], y[1]),
                    std::min(x[2], y[2]), std::min(x[3], y[3]));
};
ivec2 min(ivec2 x, int32_t y) {
  return make_ivec2(std::min(x[0], y), std::min(x[1], y));
};
ivec3 min(ivec3 x, int32_t y) {
  return make_ivec3(std::min(x[0], y), std::min(x[1], y), std::min(x[2], y));
};
ivec4 min(ivec4 x, int32_t y) {
  return make_ivec4(std::min(x[0], y), std::min(x[1], y), std::min(x[2], y),
                    std::min(x[3], y));
};
uint32_t min(uint32_t x, uint32_t y) { return std::min(x, y); }
uvec2 min(uvec2 x, uvec2 y) {
  return make_uvec2(std::min(x[0], y[0]), std::min(x[1], y[1]));
};
uvec3 min(uvec3 x, uvec3 y) {
  return make_uvec3(std::min(x[0], y[0]), std::min(x[1], y[1]),
                    std::min(x[2], y[2]));
};
uvec4 min(uvec4 x, uvec4 y) {
  return make_uvec4(std::min(x[0], y[0]), std::min(x[1], y[1]),
                    std::min(x[2], y[2]), std::min(x[3], y[3]));
};
uvec2 min(uvec2 x, uint32_t y) {
  return make_uvec2(std::min(x[0], y), std::min(x[1], y));
};
uvec3 min(uvec3 x, uint32_t y) {
  return make_uvec3(std::min(x[0], y), std::min(x[1], y), std::min(x[2], y));
};
uvec4 min(uvec4 x, uint32_t y) {
  return make_uvec4(std::min(x[0], y), std::min(x[1], y), std::min(x[2], y),
                    std::min(x[3], y));
};

float max(float x, float y) { return std::max(x, y); }
vec2 max(vec2 x, vec2 y) {
  return make_vec2(std::max(x[0], y[0]), std::max(x[1], y[1]));
};
vec3 max(vec3 x, vec3 y) {
  return make_vec3(std::max(x[0], y[0]), std::max(x[1], y[1]),
                   std::max(x[2], y[2]));
};
vec4 max(vec4 x, vec4 y) {
  return make_vec4(std::max(x[0], y[0]), std::max(x[1], y[1]),
                   std::max(x[2], y[2]), std::max(x[3], y[3]));
};
vec2 max(vec2 x, float y) {
  return make_vec2(std::max(x[0], y), std::max(x[1], y));
};
vec3 max(vec3 x, float y) {
  return make_vec3(std::max(x[0], y), std::max(x[1], y), std::max(x[2], y));
};
vec4 max(vec4 x, float y) {
  return make_vec4(std::max(x[0], y), std::max(x[1], y), std::max(x[2], y),
                   std::max(x[3], y));
};
int32_t max(int32_t x, int32_t y) { return std::max(x, y); }
ivec2 max(ivec2 x, ivec2 y) {
  return make_ivec2(std::max(x[0], y[0]), std::max(x[1], y[1]));
};
ivec3 max(ivec3 x, ivec3 y) {
  return make_ivec3(std::max(x[0], y[0]), std::max(x[1], y[1]),
                    std::max(x[2], y[2]));
};
ivec4 max(ivec4 x, ivec4 y) {
  return make_ivec4(std::max(x[0], y[0]), std::max(x[1], y[1]),
                    std::max(x[2], y[2]), std::max(x[3], y[3]));
};
ivec2 max(ivec2 x, int32_t y) {
  return make_ivec2(std::max(x[0], y), std::max(x[1], y));
};
ivec3 max(ivec3 x, int32_t y) {
  return make_ivec3(std::max(x[0], y), std::max(x[1], y), std::max(x[2], y));
};
ivec4 max(ivec4 x, int32_t y) {
  return make_ivec4(std::max(x[0], y), std::max(x[1], y), std::max(x[2], y),
                    std::max(x[3], y));
};
uint32_t max(uint32_t x, uint32_t y) { return std::max(x, y); }
uvec2 max(uvec2 x, uvec2 y) {
  return make_uvec2(std::max(x[0], y[0]), std::max(x[1], y[1]));
};
uvec3 max(uvec3 x, uvec3 y) {
  return make_uvec3(std::max(x[0], y[0]), std::max(x[1], y[1]),
                    std::max(x[2], y[2]));
};
uvec4 max(uvec4 x, uvec4 y) {
  return make_uvec4(std::max(x[0], y[0]), std::max(x[1], y[1]),
                    std::max(x[2], y[2]), std::max(x[3], y[3]));
};
uvec2 max(uvec2 x, uint32_t y) {
  return make_uvec2(std::max(x[0], y), std::max(x[1], y));
};
uvec3 max(uvec3 x, uint32_t y) {
  return make_uvec3(std::max(x[0], y), std::max(x[1], y), std::max(x[2], y));
};
uvec4 max(uvec4 x, uint32_t y) {
  return make_uvec4(std::max(x[0], y), std::max(x[1], y), std::max(x[2], y),
                    std::max(x[3], y));
};

/* mix functions */
float mix(float x, float y, float a) { return x * (1.0 - a) + y * a; }
vec2 mix(vec2 x, vec2 y, vec2 a) {
  return make_vec2(x[0] * (1.0f - a[0]) + y[0] * a[0],
                   x[1] * (1.0f - a[1]) + y[1] * a[1]);
};
vec3 mix(vec3 x, vec3 y, vec3 a) {
  return make_vec3(x[0] * (1.0f - a[0]) + y[0] * a[0],
                   x[1] * (1.0f - a[1]) + y[1] * a[1],
                   x[2] * (1.0f - a[2]) + y[2] * a[2]);
};
vec4 mix(vec4 x, vec4 y, vec4 a) {
  return make_vec4(
      x[0] * (1.0f - a[0]) + y[0] * a[0], x[1] * (1.0f - a[1]) + y[1] * a[1],
      x[2] * (1.0f - a[2]) + y[2] * a[2], x[3] * (1.0f - a[3]) + y[3] * a[3]);
};
vec2 mix(vec2 x, vec2 y, float a) {
  return make_vec2(x[0] * (1.0f - a) + y[0] * a, x[1] * (1.0f - a) + y[1] * a);
};
vec3 mix(vec3 x, vec3 y, float a) {
  return make_vec3(x[0] * (1.0f - a) + y[0] * a, x[1] * (1.0f - a) + y[1] * a,
                   x[2] * (1.0f - a) + y[2] * a);
};
vec4 mix(vec4 x, vec4 y, float a) {
  return make_vec4(x[0] * (1.0f - a) + y[0] * a, x[1] * (1.0f - a) + y[1] * a,
                   x[2] * (1.0f - a) + y[2] * a, x[3] * (1.0f - a) + y[3] * a);
};

/* clamp functions */
float clamp(float x, float mi, float ma) { return min(max(x, mi), ma); }
vec2 clamp(vec2 x, vec2 mi, vec2 ma) {
  return make_vec2(clamp(x[0], mi[0], ma[0]), clamp(x[1], mi[1], ma[1]));
}
vec3 clamp(vec3 x, vec3 mi, vec3 ma) {
  return make_vec3(clamp(x[0], mi[0], ma[0]), clamp(x[1], mi[1], ma[1]),
                   clamp(x[2], mi[2], ma[2]));
}
vec4 clamp(vec4 x, vec4 mi, vec4 ma) {
  return make_vec4(clamp(x[0], mi[0], ma[0]), clamp(x[1], mi[1], ma[1]),
                   clamp(x[2], mi[2], ma[2]), clamp(x[3], mi[3], ma[3]));
}
vec2 clamp(vec2 x, float mi, float ma) {
  return make_vec2(clamp(x[0], mi, ma), clamp(x[1], mi, ma));
}
vec3 clamp(vec3 x, float mi, float ma) {
  return make_vec3(clamp(x[0], mi, ma), clamp(x[1], mi, ma),
                   clamp(x[2], mi, ma));
}
vec4 clamp(vec4 x, float mi, float ma) {
  return make_vec4(clamp(x[0], mi, ma), clamp(x[1], mi, ma),
                   clamp(x[2], mi, ma), clamp(x[3], mi, ma));
}

/* bit count functions */

int32_t bitCount(int32_t value) {
    return __builtin_popcount(value);
}
 
int32_t bitCount(uint32_t value) {
    return __builtin_popcount(value);
}