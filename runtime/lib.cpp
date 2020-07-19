// export only the symbols in shared to allow efficient dce.
#pragma GCC visibility push(default)
extern "C" {
#include "./generated/shared.h"
}
#pragma GCC visibility pop

#include "debug.hpp"
#include "runtime.hpp"

extern "C" {

void *cpt_new_kernel(int32_t num_t) {
  cpt_log("cpt_new_kernel %d", num_t);
  Kernel *k = new Kernel(num_t);
  return static_cast<void *>(k);
}

struct cpt_error_t cpt_dispatch_kernel(void *k, cpt_data d, int32_t x, int32_t y,
                                   int32_t z) {
  cpt_log("cpt_dispatch_kernel %p %d %d %d", k, x, y, z);
  Kernel *kt = static_cast<Kernel *>(k);
  return kt->dispatch(d, x, y, z);
}

void cpt_free_kernel(void *k) {
  cpt_log("cpt_free_kernel %p", k);
  Kernel *kt = static_cast<Kernel *>(k);
  delete kt;
}
}