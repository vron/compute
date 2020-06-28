extern "C" {
  #include "shared.h" // the definitions shared with go / lib user
}
namespace {
  #include "types.hpp"    // implementations of glsl types and functions
  #include "kernel.hpp"  // the kernel code generated from the glsl shader
  #include "threads.hpp" // wraps system threading libraries to support platforms
  #include "runtime.hpp" // the runtime scheduling execution
}

extern "C" {
// This defines the C api that will be exposed by
// the final library. These methods simply wrap the
// C++ implementation.

void *cpt_new_kernel(int32_t num_threads) {
  kernel *k = new kernel(num_threads);
  return static_cast<void *>(k);
}

void cpt_free_kernel(void *k) {
  kernel *kt = static_cast<kernel *>(k);
  delete kt;
}

int cpt_dispatch_kernel(void *k, cpt_data d, int32_t numx, int32_t numy, int32_t numz) {
  kernel *kt = static_cast<kernel *>(k);
  return kt->dispatch(d, numx, numy, numz);
}
}