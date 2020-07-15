#pragma once
#include "common.hpp"
#include "vec.hpp"
#include "vec_make.hpp"
#include "vec_math.hpp"

struct mat2 {
  vec2 c[2];


  vec2 &operator[](int index);
};

struct mat3 {
  vec3 c[3];


  vec3 &operator[](int index);
};

struct mat4 {
  vec4 c[4];


  vec4 &operator[](int index);
};

mat2 operator*(mat2 lhs, const mat2 &rhs) {
  for (int i = 0; i < 2; i++) {
    vec2 row = make_vec2(lhs.c[0][i], lhs.c[1][i]);
    lhs.c[0][i] = dot(row, rhs.c[0]);
    lhs.c[1][i] = dot(row, rhs.c[1]);
  }
  return lhs;
};

mat3 operator*(mat3 lhs, const mat3 &rhs) {
  for (int i = 0; i < 3; i++) {
    vec3 row = make_vec3(lhs.c[0][i], lhs.c[1][i], lhs.c[2][i]);
    lhs.c[0][i] = dot(row, rhs.c[0]);
    lhs.c[1][i] = dot(row, rhs.c[1]);
    lhs.c[2][i] = dot(row, rhs.c[2]);
  }
  return lhs;
};

mat4 operator*(mat4 lhs, const mat4 &rhs) {
  for (int i = 0; i < 4; i++) {
    vec4 row = make_vec4(lhs.c[0][i], lhs.c[1][i], lhs.c[2][i], lhs.c[3][i]);
    lhs.c[0][i] = dot(row, rhs.c[0]);
    lhs.c[1][i] = dot(row, rhs.c[1]);
    lhs.c[2][i] = dot(row, rhs.c[2]);
    lhs.c[3][i] = dot(row, rhs.c[3]);
  }
  return lhs;
};

vec2 operator*(mat2 lhs, const vec2 &rhs) {
  return make_vec2(lhs[0][0] * rhs[0] + lhs[1][0] * rhs[1],
                   lhs[0][1] * rhs[1] + lhs[1][1] * rhs[1]);
};

vec3 operator*(mat3 lhs, const vec3 &rhs) {
  return make_vec3(lhs[0][0] * rhs[0] + lhs[1][0] * rhs[1] + lhs[2][0] * rhs[2],
                   lhs[0][1] * rhs[0] + lhs[1][1] * rhs[1] + lhs[2][1] * rhs[2],
                   lhs[0][2] * rhs[0] + lhs[1][2] * rhs[1] +
                       lhs[2][2] * rhs[2]);
};

vec4 operator*(mat4 lhs, const vec4 &rhs) {
  return make_vec4(lhs[0][0] * rhs[0] + lhs[1][0] * rhs[1] +
                       lhs[2][0] * rhs[2] + lhs[3][0] * rhs[3],
                   lhs[0][1] * rhs[0] + lhs[1][1] * rhs[1] +
                       lhs[2][1] * rhs[2] + lhs[3][1] * rhs[3],
                   lhs[0][2] * rhs[0] + lhs[1][2] * rhs[1] +
                       lhs[2][2] * rhs[2] + lhs[3][2] * rhs[3],
                   lhs[0][3] * rhs[0] + lhs[1][3] * rhs[1] +
                       lhs[2][3] * rhs[2] + lhs[3][3] * rhs[3]);
};

vec2 &mat2::operator[](int index) {
  index = index % 2;
  switch (index) {
  case 0:
    return (this->c[0]);
  case 1:
    return (this->c[1]);
  }
  __builtin_unreachable();
}

vec3 &mat3::operator[](int index) {
  index = index % 3;
  switch (index) {
  case 0:
    return (this->c[0]);
  case 1:
    return (this->c[1]);
  case 2:
    return (this->c[2]);
  }
  __builtin_unreachable();
}

vec4 &mat4::operator[](int index) {
  index = index % 4;
  switch (index) {
  case 0:
    return (this->c[0]);
  case 1:
    return (this->c[1]);
  case 2:
    return (this->c[2]);
  case 3:
    return (this->c[3]);
  }
  __builtin_unreachable();
}

  void always_inline from_api(mat2 *me, float *arr) {
    me->c[0][0] = arr[0];
    me->c[0][1] = arr[1];
    me->c[1][0] = arr[2];
    me->c[1][1] = arr[3];
  };
  void always_inline from_api(mat3 *me, float *arr) {
    me->c[0][0] = arr[0];
    me->c[0][1] = arr[1];
    me->c[0][2] = arr[2];
    me->c[1][0] = arr[3];
    me->c[1][1] = arr[4];
    me->c[1][2] = arr[5];
    me->c[2][0] = arr[6];
    me->c[2][1] = arr[7];
    me->c[2][2] = arr[8];
  };
  void always_inline from_api(mat4 *me, float *arr) {
    me->c[0][0] = arr[0];
    me->c[0][1] = arr[1];
    me->c[0][2] = arr[2];
    me->c[0][3] = arr[3];
    me->c[1][0] = arr[4];
    me->c[1][1] = arr[5];
    me->c[1][2] = arr[6];
    me->c[1][3] = arr[7];
    me->c[2][0] = arr[8];
    me->c[2][1] = arr[9];
    me->c[2][2] = arr[10];
    me->c[2][3] = arr[11];
    me->c[3][0] = arr[12];
    me->c[3][1] = arr[13];
    me->c[3][2] = arr[14];
    me->c[3][3] = arr[15];
  };