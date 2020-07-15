#pragma once
#include "vec.hpp"
#include "vec_make.hpp"
#include <cmath>

float always_inline cross(vec2 a, vec2 b) { return a[0] * b[1] - a[1] * b[0]; };
float always_inline dot(vec2 a, vec2 b) { return a[0] * b[0] + a[1] * b[1]; };
float always_inline length(vec2 a) { return sqrt(dot(a, a)); };
float always_inline distance(vec2 a, vec2 b) { return length(a - b); };
vec2 always_inline normalize(vec2 a) { return a / length(a); };

vec3 always_inline cross(vec3 a, vec3 b) {
  return make_vec3(a[1] * b[2] - a[2] * b[1], a[2] * b[0] - a[0] * b[2],
                   a[0] * b[1] - a[1] * b[0]);
}
float always_inline dot(vec3 a, vec3 b) {
  return a[0] * b[0] + a[1] * b[1] + a[2] * b[2];
};
float always_inline length(vec3 a) { return sqrt(dot(a, a)); };
float always_inline distance(vec3 a, vec3 b) { return length(a - b); };
vec3 always_inline normalize(vec3 a) { return a / length(a); };

float always_inline dot(vec4 a, vec4 b) {
  return a[0] * b[0] + a[1] * b[1] + a[2] * b[2] + a[3] * b[3];
};
float always_inline length(vec4 a) { return sqrt(dot(a, a)); };
float always_inline distance(vec4 a, vec4 b) { return length(a - b); };
vec4 always_inline normalize(vec4 a) { return a / length(a); };
