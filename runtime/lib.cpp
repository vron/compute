// export only the symbols in shared to allow efficient dce.
#pragma GCC visibility push(default)
extern "C" {
#include "shared.h"
}
#pragma GCC visibility pop

// clang-format off
struct thread_data;
#include "debug.hpp"
#include "threads.hpp" // wraps system threading libraries to support platforms
#include "types.hpp"   // implementations of glsl types and functions
#include "kernel.hpp"  // the kernel code generated from the glsl shader
#include "runtime.hpp" // the runtime scheduling execution
// clang-format on

extern "C" {

void *cpt_new_kernel(int32_t num_t) {
  kernel *k = new kernel(num_t);
  return static_cast<void *>(k);
}

struct error_t cpt_dispatch_kernel(void *k, cpt_data d, int32_t x, int32_t y,
                                   int32_t z) {
  cpt_log("dispatch %p %d %d %d", k, x, y, z);
  kernel *kt = static_cast<kernel *>(k);
  return kt->dispatch(d, x, y, z);
}

void cpt_free_kernel(void *k) {
  cpt_log("free %p", k);
  kernel *kt = static_cast<kernel *>(k);
  delete kt;
}
}