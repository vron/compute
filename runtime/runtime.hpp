#pragma once

struct thread_data;

#include "debug.hpp"
#include <cstring>
#include <thread>
#include "./types/types.hpp"
#include "threads.hpp"
#include "./generated/shader.hpp"

struct thread_data {
  pthread_barrier_t *barrier;
  std::thread thread;
  shader *kernel;
};

void thread_core(struct thread_data *td) {
  td->kernel->thread = td;
  cpt_log("thread calling main %p", td);
  td->kernel->main();
  cpt_log("main called %p", td);
  return;
}

void shader::barrier() {
  cpt_trace("barrier call %p", this->thread->barrier);
  pthread_barrier_wait(this->thread->barrier);
  cpt_trace("barrier release %p", this);
}

struct kernel {
  unsigned int num_threads;

  char error_msg[1024];
  int error_no;

  kernel(int32_t num_t) {
    if (num_t < 1) {
      this->set_error(EINVAL, "must use at least 1 thread");
      return;
    }
    this->num_threads = (unsigned int)num_threads;
  }

  ~kernel() {}

private:
  bool set_error(int no, const char *msg);
  struct cpt_error_t error();
  bool ensure_alignments(cpt_data d);

public:
  struct cpt_error_t dispatch(cpt_data d, int32_t nx, int32_t ny, int32_t nz);
};

bool kernel::set_error(int no, const char *msg) {
  if (this->error_no != 0)
    return false;
  strncpy(&this->error_msg[0], msg, 1023);
  this->error_no = no;
  return false;
}

struct cpt_error_t kernel::error() {
  return (struct cpt_error_t){this->error_no, &this->error_msg[0]};
}

struct cpt_error_t kernel::dispatch(cpt_data d, int32_t nx, int32_t ny,
                                int32_t nz) {
  if (this->error_no)
    return this->error();
  if (!this->ensure_alignments(d))
    return this->error();

  struct thread_data threads[_cpt_WG_SIZE_Z * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X];
  pthread_barrier_t barrier;

  for (uint32_t gz = 0; gz < (uint32_t)nz; ++gz) {
    for (uint32_t gy = 0; gy < (uint32_t)ny; ++gy) {
      for (uint32_t gx = 0; gx < (uint32_t)nx; ++gx) {

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
              shader *k = new shader();
              
              k->gl_NumWorkGroups = make_uvec3(nx, ny, nz);
              k->gl_WorkGroupID = make_uvec3(gx, gy, gz);
              k->gl_WorkGroupSize = make_uvec3(_cpt_WG_SIZE_X, _cpt_WG_SIZE_Y, _cpt_WG_SIZE_Z);
              k->gl_LocalInvocationID = make_uvec3(lx, ly, lz);
              k->gl_GlobalInvocationID =
                  k->gl_WorkGroupID * k->gl_WorkGroupSize +
                  k->gl_LocalInvocationID;
              k->gl_LocalInvocationIndex = lx + ly * _cpt_WG_SIZE_X +
                                           lz * _cpt_WG_SIZE_X * _cpt_WG_SIZE_Y;

              k->set_shared_data(sd);
              k->set_data(d);

              threads[index].kernel = k;
              threads[index].barrier = &barrier;

              try {
                cpt_verbose("about to create %p", &threads[index]);
                threads[index].thread =
                    std::thread(thread_core, &threads[index]);
              } catch (...) {
                cpt_log("ERROR CRATING thread");
              }
              cpt_verbose("thread created now");
            }
          }
        }

        cpt_verbose("loop done");
        // wait for all the threads to finish and join them
        for (uint32_t lz = 0; lz < _cpt_WG_SIZE_Z; ++lz) {
          for (uint32_t ly = 0; ly < _cpt_WG_SIZE_Y; ++ly) {
            for (uint32_t lx = 0; lx < _cpt_WG_SIZE_X; ++lx) {
              long index = lx + ly * _cpt_WG_SIZE_X +
                           lz * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X;

              cpt_verbose("about to join %ld", index);
              try {
                threads[index].thread.join();
              } catch (...) {
                cpt_log("ERROR JOINING thread");
              }
              cpt_log("joined");
              free(threads[index].kernel);
            }
          }
        }
        cpt_log("wg done");
        shader::free_shared_data(sd);

        pthread_barrier_destroy(&barrier);
      }
    }
  }

  return this->error();
}

#include "./generated/align.hpp"
