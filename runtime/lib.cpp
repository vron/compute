extern "C" {
#include "shared.h" // the definitions shared with go / lib user
}

// clang-format off
namespace { 
  struct thread_data;
  #include "debug.hpp"
  #include "threads.hpp" // wraps system threading libraries to support platforms
  #include "types.hpp"   // implementations of glsl types and functions
  #include "kernel.hpp"  // the kernel code generated from the glsl shader
  #include "runtime.hpp" // the runtime scheduling execution
}
// clang-format on

extern "C" {
// The exposed library wrapping C++ world

void *cpt_new_kernel(int32_t num_threads) {
  d_log("num_threads=%d", num_threads);
  init_debug(); // TODO: This is not thread safe!
  kernel *k = new kernel(num_threads);
  return static_cast<void *>(k);
}

void cpt_free_kernel(void *k) {
  d_log("free %p", k);
  kernel *kt = static_cast<kernel *>(k);
  delete kt;
}

int cpt_dispatch_kernel(void *k, cpt_data d, int32_t numx, int32_t numy,
                        int32_t numz) {
  d_log("dispatch %p %d %d %d", k, numx, numy, numz);
  kernel *kt = static_cast<kernel *>(k);
  return kt->dispatch(d, numx, numy, numz);
}


}