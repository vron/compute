
struct thread_data {
  pthread_barrier_t *barrier;
  pthread_t thread;
  kernel_comp *kernel;
};

void *thread_core(void *ptr) {
  struct thread_data *td = (struct thread_data *)(ptr);
  d_log("starting thread %p", td);

  td->kernel->thread = td;
  td->kernel->main();
  return NULL;
}

void kernel_comp::barrier() {
  d_trace("barrier call %p", this);
  pthread_barrier_wait(this->thread->barrier);
  d_trace("barrier release %p", this);
}

struct kernel {
  unsigned int num_threads;

  kernel(int32_t num_threads) { this->num_threads = (unsigned int)num_threads; }

  ~kernel() {
    // TODO: clean up threads etc. etc.
  }

  int ensure_alignments(cpt_data d) {
#include "align.hpp"
  }

  int dispatch(cpt_data d, int32_t nx, int32_t ny, int32_t nz) {
    int erc = this->ensure_alignments(d);
    if (erc)
      return erc; // TODO: Move over the error structuring to C world so not
                  // dependent on Go...

    struct thread_data *threads;
    threads = (struct thread_data *)malloc(sizeof(*threads) * _cpt_WG_SIZE_Z *
                                           _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X);
    pthread_barrier_t barrier;

    uint32_t numx = (uint32_t)nx;
    uint32_t numy = (uint32_t)ny;
    uint32_t numz = (uint32_t)nz;
    for (uint32_t gz = 0; gz < numz; ++gz) {
      for (uint32_t gy = 0; gy < numy; ++gy) {
        for (uint32_t gx = 0; gx < numx; ++gx) {

          int ercc = pthread_barrier_init(
              &barrier, NULL, _cpt_WG_SIZE_Z * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X);
          d_log("init %d %p", ercc, &barrier);
          if(ercc)
            return ercc + 150;

          auto sd = kernel_comp::create_shared_data();


          for (uint32_t lz = 0; lz < _cpt_WG_SIZE_Z; ++lz) {
            for (uint32_t ly = 0; ly < _cpt_WG_SIZE_Y; ++ly) {
              for (uint32_t lx = 0; lx < _cpt_WG_SIZE_X; ++lx) {
                long index = lx + ly * _cpt_WG_SIZE_X +
                             lz * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X;
                d_verbose("in loop %d %ld",
                        _cpt_WG_SIZE_Z * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X,
                        index);
                kernel_comp *k = new kernel_comp(); // TODO: re-use the
                                                    // allocation across calls
                                                    
                // TODO: can we avoid all these calcs?
                k->gl_NumWorkGroups.z = nz;
                k->gl_NumWorkGroups.y = ny;
                k->gl_NumWorkGroups.x = nx;
                k->gl_WorkGroupID.z = gz;
                k->gl_WorkGroupID.y = gy;
                k->gl_WorkGroupID.x = gx;
                k->gl_WorkGroupSize.z = _cpt_WG_SIZE_Z;
                k->gl_WorkGroupSize.y = _cpt_WG_SIZE_Y;
                k->gl_WorkGroupSize.x = _cpt_WG_SIZE_X;
                k->gl_LocalInvocationID.z = lz;
                k->gl_LocalInvocationID.y = ly;
                k->gl_LocalInvocationID.x = lx;
                k->gl_GlobalInvocationID = k->gl_WorkGroupID * make_uvec3(_cpt_WG_SIZE_X, _cpt_WG_SIZE_Y, _cpt_WG_SIZE_Z) + k->gl_LocalInvocationID;
                k->gl_LocalInvocationIndex = lx + ly*_cpt_WG_SIZE_X + lx*_cpt_WG_SIZE_X*_cpt_WG_SIZE_Y;
                
                k->set_shared_data(sd);
                int erno = k->set_data(d);
                if (erno != 0) 
                  return erno;
        
                threads[index].kernel = k;
                threads[index].barrier = &barrier;

                // for each one, create a thread that we will use:
                int errno;
                d_verbose("about to creae");
                errno = pthread_create(&threads[index].thread, NULL, thread_core,
                                       (void *)&threads[index]);
                d_verbose("thread create no %d", errno);
                if (errno)
                  return errno; // TODO: memory leas... (this entire functin...)
                  
              }
            }
          }

          // wait for all the threads to finish and join them
          for (uint32_t lz = 0; lz < _cpt_WG_SIZE_Z; ++lz) {
            for (uint32_t ly = 0; ly < _cpt_WG_SIZE_Y; ++ly) {
              for (uint32_t lx = 0; lx < _cpt_WG_SIZE_X; ++lx) {
                long index = lx + ly * _cpt_WG_SIZE_X +
                             lz * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X;
                int s = pthread_join(threads[index].thread, NULL);
                d_log("ERC %d", s);
                if (s != 0)
                  return s +1000;
                  //return error("error joining threads: %d", s);
                free(threads[index].kernel);
              }
            }
          }

          
          kernel_comp::free_shared_data(sd);

          // TODO: Free!
          pthread_barrier_destroy(&barrier);
        }
      }
    }

    free(threads);

    return 0;
  }
};
