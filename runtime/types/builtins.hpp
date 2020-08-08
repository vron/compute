#pragma once
#include <cstdint>
#include <cstring>

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
