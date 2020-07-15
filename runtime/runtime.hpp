#pragma once

#include "debug.hpp"
#include <cstring>

struct thread_data {
  pthread_barrier_t *barrier;
  pthread_t thread;
  shader *kernel;
};

void *thread_core(void *ptr) {
  struct thread_data *td = (struct thread_data *)(ptr);
  cpt_log("starting thread %p", td);

  td->kernel->thread = td;
  td->kernel->main();
  return NULL;
}

void shader::barrier() {
  cpt_trace("barrier call %p", this);
  pthread_barrier_wait(this->thread->barrier);
  cpt_trace("barrier release %p", this);
}

struct kernel {
  unsigned int num_threads;

  char error_msg[1024];
  int error_no;

  kernel(int32_t num_t) {
    if (num_t < 1) {
    }
    this->num_threads = (unsigned int)num_threads;
  }

  ~kernel() {}

private:
  bool set_error(int no, const char *msg);
  struct error_t error();
  bool ensure_alignments(cpt_data d);

public:
  struct error_t dispatch(cpt_data d, int32_t nx, int32_t ny, int32_t nz);
};

bool kernel::set_error(int no, const char *msg) {
  if (this->error_no != 0)
    return false;
  strncpy(&this->error_msg[0], msg, 1023);
  this->error_no = no;
  return false;
}

struct error_t kernel::error() {
  return (struct error_t){this->error_no, &this->error_msg[0]};
}

struct error_t kernel::dispatch(cpt_data d, int32_t nx, int32_t ny,
                                int32_t nz) {
  if (this->error_no)
    return this->error();
  if (!this->ensure_alignments(d))
    return this->error();

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
        cpt_log("init %d %p", ercc, &barrier);
        if (ercc) {
          this->set_error(ercc, "error creating thread barrier");
          return this->error();
        }

        auto sd = shader::create_shared_data();

        for (uint32_t lz = 0; lz < _cpt_WG_SIZE_Z; ++lz) {
          for (uint32_t ly = 0; ly < _cpt_WG_SIZE_Y; ++ly) {
            for (uint32_t lx = 0; lx < _cpt_WG_SIZE_X; ++lx) {
              long index = lx + ly * _cpt_WG_SIZE_X +
                           lz * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X;
              cpt_verbose("in loop %d %ld",
                          _cpt_WG_SIZE_Z * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X,
                          index);
              shader *k = new shader(); // TODO: re-use the
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
              k->gl_GlobalInvocationID =
                  k->gl_WorkGroupID * make_uvec3(_cpt_WG_SIZE_X, _cpt_WG_SIZE_Y,
                                                 _cpt_WG_SIZE_Z) +
                  k->gl_LocalInvocationID;
              k->gl_LocalInvocationIndex = lx + ly * _cpt_WG_SIZE_X +
                                           lz * _cpt_WG_SIZE_X * _cpt_WG_SIZE_Y;

              k->set_shared_data(sd);
              k->set_data(d);

              threads[index].kernel = k;
              threads[index].barrier = &barrier;

              // for each one, create a thread that we will use:
              int no;
              cpt_verbose("about to creae");
              no = pthread_create(&threads[index].thread, NULL, thread_core,
                                  (void *)&threads[index]);
              cpt_verbose("thread create no %d", no);
              if (no) {
                this->set_error(no, "error creating threads");
                return this->error();
                // TODO: memory leas... (this entire functin...)
              }
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
              cpt_log("ERC %d", s);
              if (s != 0) {
                this->set_error(s, "error joining threads");
                return this->error();
              }
              // return error("error joining threads: %d", s);
              free(threads[index].kernel);
            }
          }
        }

        shader::free_shared_data(sd);

        // TODO: Free!
        pthread_barrier_destroy(&barrier);
      }
    }
  }

  free(threads);

  return this->error();
}

#include "align.hpp"