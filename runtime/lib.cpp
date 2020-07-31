// export only the symbols in shared to allow efficient dce.
#pragma GCC visibility push(default)
extern "C" {
#include "./generated/shared.h"
}
#pragma GCC visibility pop

#include "runtime.hpp"

extern "C" {

void *cpt_new_kernel(int32_t num_t, int32_t stack_size) {
  Kernel *k = new Kernel(num_t, stack_size);
  return static_cast<void *>(k);
}

struct cpt_error_t cpt_dispatch_kernel(void *k, cpt_data d, int32_t x,
                                       int32_t y, int32_t z) {
  Kernel *kt = static_cast<Kernel *>(k);

  // TODO: here we actually must memcpy the data -but rely on this being optimized out
  // since we do it only for type punning
  // https://stackoverflow.com/questions/98650/what-is-the-strict-aliasing-rule
  
  return kt->dispatch(d, x, y, z);
}

void cpt_free_kernel(void *k) {
  Kernel *kt = static_cast<Kernel *>(k);
  delete kt;
}
}
