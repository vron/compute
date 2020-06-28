
struct kernel {
  unsigned int num_threads;

  kernel(int32_t num_threads) {
      this->num_threads = (unsigned int)num_threads;

  }

  ~kernel() {
    // TODO: clean up threads etc. etc.
  }

  int ensure_alignments(cpt_data d) {
      // static and dynamic chc of alignments and data sizes
      // TODO: Do not do on every call all of it.
      #include "align.hpp"
  }

  int dispatch(cpt_data d, int32_t nx, int32_t ny, int32_t nz) {
      int erc = this->ensure_alignments(d);
      if (erc) {
          return erc;
      }
      // important we allocate this one on heap and not stac!
    kernel_comp* k = new kernel_comp(); // TODO: re-use the allocation across calls

    int erno = k->set_data(d);
    if (erno != 0) {
        delete k;
        return erno;
    }

    // loop over global invocations
    uint32_t numx = (uint32_t)nx;
    uint32_t numy = (uint32_t)ny;
    uint32_t numz = (uint32_t)nz;
    for (uint32_t gz = 0; gz < numz; ++gz) {
      k->gl_GlobalInvocationID.z = gz;
      for (uint32_t gy = 0; gy < numy; ++gy) {
        k->gl_GlobalInvocationID.y = gy;
        for (uint32_t gx = 0; gx < numx; ++gx) {
          k->gl_GlobalInvocationID.x = gx;
          // loop over the local work group
          for (uint32_t lz = 0; lz < _cpt_WG_SIZE_Z; ++lz) {
            k->gl_LocalInvocationID.z = lz;
            for (uint32_t ly = 0; ly < _cpt_WG_SIZE_Y; ++ly) {
              k->gl_LocalInvocationID.y = ly;
              for (uint32_t lx = 0; lx < _cpt_WG_SIZE_X; ++lx) {
                k->gl_LocalInvocationID.x = lx;
                k->main();
              }
            }
          }
        }
      }
    }
    delete k;
    return 0;
  }
};
