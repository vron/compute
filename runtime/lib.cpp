// export only the symbols in shared to allow efficient dce.
#pragma GCC visibility push(default)
extern "C" {
#include "./generated/shared.h"
}
#pragma GCC visibility pop

#include "runtime.hpp"
#include "generated/usertypes.hpp"
#include "./generated/align.hpp"

extern "C" {

void *cpt_new_kernel(int32_t num_t, int32_t stack_size) {
  Kernel *k = new Kernel(num_t, stack_size);
  return static_cast<void *>(k);
}

struct cpt_error_t cpt_dispatch_kernel(void *k, cpt_data d, int32_t x,
                                       int32_t y, int32_t z) {
  Kernel *kt = static_cast<Kernel *>(k);

  // TODO: we are using undefined behaviour here, but see no other way to avoid 2 memcopies...  https://stackoverflow.com/questions/98650/what-is-the-strict-aliasing-rule
  cptc_data *data = (cptc_data*)(&d);
  return kt->dispatch(data, x, y, z);
}

void cpt_free_kernel(void *k) {
  Kernel *kt = static_cast<Kernel *>(k);
  delete kt;
}
}
