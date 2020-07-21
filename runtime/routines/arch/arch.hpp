#pragma once
#include <cstddef>
#include <cstdint>

#if defined(__x86_64__) && (defined(__linux__) || defined(__APPLE__))
// System V AMD64 ABI here
#define ARCH_reg_state 10
// the first argument to the func should also be set to to, here implicitly
// handled since we call this function with it..
void ARCH_switch(void *to, void *from) __asm__("amd64_nix_switch");
void ARCH_set_register_state(void **registers, void *fp, void *stack,
                             size_t stack_size) {
  // on x86 the stack grows downwards, so we want to find a pointer to the
  // end of the stack instead, but we want one that is aligned on a 16 byte
  // boundary
  uintptr_t u_p = (uintptr_t)(stack) + (uintptr_t)(stack_size)-1;
  u_p = u_p & (~0xF);
  // store the function pointer where it will be restored from the assembly
  registers[6] = fp;
  registers[7] = (void *)(u_p - 24);
};

#elif (defined(_WIN64) && defined(__x86_64__))
#define ARCH_reg_state 32
void ARCH_switch(void *to, void *from) __asm__("amd64_win_switch");
void ARCH_set_register_state(void **registers, void *fp, void *stack,
                             size_t stack_size) {
  uintptr_t u_p = (uintptr_t)(stack) + (uintptr_t)(stack_size)-1;
  u_p = u_p & (~0xF);
  registers[9] = fp;
  registers[4] = (void *)(u_p - 24);
};
#else
#error "unsupported (arch, ABI) combination encountered"
#endif
