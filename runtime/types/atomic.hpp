#pragma once
#include "common.hpp"
#include <atomic>
#include <cstdint>

/* barrier */
void always_inline memoryBarrier() {
  std::atomic_thread_fence(std::memory_order_seq_cst);
}

void always_inline memoryBarrierBuffer() { memoryBarrier(); }

void always_inline memoryBarrierImage() { memoryBarrier(); }

void always_inline memoryBarrierShared() { memoryBarrier(); }

/* atomic operations */

uint32_t always_inline atomicAdd(uint32_t *mem, uint32_t data) {
  return __atomic_add_fetch(mem, data, __ATOMIC_SEQ_CST);
}

uint32_t always_inline atomicAnd(uint32_t *mem, uint32_t data) {
  return __atomic_and_fetch(mem, data, __ATOMIC_SEQ_CST);
}

uint32_t always_inline atomicOr(uint32_t *mem, uint32_t data) {
  return __atomic_or_fetch(mem, data, __ATOMIC_SEQ_CST);
}

uint32_t always_inline atomicXor(uint32_t *mem, uint32_t data) {
  return __atomic_xor_fetch(mem, data, __ATOMIC_SEQ_CST);
}

uint32_t always_inline atomicMin(uint32_t *mem, uint32_t data) {
  return __atomic_fetch_min(mem, data, __ATOMIC_SEQ_CST);
}

uint32_t always_inline atomicMax(uint32_t *mem, uint32_t data) {
  return __atomic_fetch_max(mem, data, __ATOMIC_SEQ_CST);
}

uint32_t always_inline atomicExchange(uint32_t *mem, uint32_t data) {
  return __atomic_exchange_n(mem, data, __ATOMIC_SEQ_CST);
}

uint32_t always_inline atomicCompSwap(uint32_t *mem, uint32_t compare,
                                      uint32_t data) {
  __atomic_compare_exchange_n(mem, &compare, data, true, __ATOMIC_SEQ_CST,
                              __ATOMIC_SEQ_CST);
  return compare;
}

int32_t always_inline atomicAdd(int32_t *mem, int32_t data) {
  return __atomic_add_fetch(mem, data, __ATOMIC_SEQ_CST);
}

int32_t always_inline atomicAnd(int32_t *mem, int32_t data) {
  return __atomic_and_fetch(mem, data, __ATOMIC_SEQ_CST);
}

int32_t always_inline atomicOr(int32_t *mem, int32_t data) {
  return __atomic_or_fetch(mem, data, __ATOMIC_SEQ_CST);
}

int32_t always_inline atomicXor(int32_t *mem, int32_t data) {
  return __atomic_xor_fetch(mem, data, __ATOMIC_SEQ_CST);
}

int32_t always_inline atomicMin(int32_t *mem, int32_t data) {
  return __atomic_fetch_min(mem, data, __ATOMIC_SEQ_CST);
}

int32_t always_inline atomicMax(int32_t *mem, int32_t data) {
  return __atomic_fetch_max(mem, data, __ATOMIC_SEQ_CST);
}

int32_t always_inline atomicExchange(int32_t *mem, int32_t data) {
  return __atomic_exchange_n(mem, data, __ATOMIC_SEQ_CST);
}

int32_t always_inline atomicCompSwap(int32_t *mem, int32_t compare,
                                     int32_t data) {
  __atomic_compare_exchange_n(mem, &compare, data, true, __ATOMIC_SEQ_CST,
                              __ATOMIC_SEQ_CST);
  return compare;
}
