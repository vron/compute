#pragma once

#include "co/routines.hpp"
#include "generated/shader.hpp"
#include "types/types.hpp"
#include "waitgroup.hpp"
#include "workpiece.hpp"
#include "workqueue.hpp"
#include "workthread.hpp"

class Kernel {
  WorkThread<struct shader *> **work_threads;
  WorkQueue<WorkPiece> queue;
  WaitGroup<int32_t> wg;

  int num_t;
  int generation;
  char error_msg[1024];
  int error_no;

public:
  Kernel(int32_t num_t, int32_t stack_size)
      : work_threads(nullptr), num_t(num_t), generation(0), error_no(0) {
    if (num_t < 1) {
      set_error(EINVAL, "must use at least 1 thread");
      return;
    }
    if (stack_size < 0) {
      stack_size = 1024 * 16;
    }
    work_threads = (WorkThread<struct shader *> **)malloc(
        sizeof(WorkThread<struct shader *> *) * num_t);
    for (int i = 0; i < num_t; i++) {
      work_threads[i] =
          new WorkThread<struct shader *>(&queue, &wg, stack_size);
    }
  }

  ~Kernel() {
    if (work_threads == nullptr)
      return;
    for (int i = 0; i < num_t; i++) {
      queue.send(
          WorkPiece(true));
    }
    for (int i = 0; i < num_t; i++) {
      delete work_threads[i];
    }
    free(work_threads);
  }

private:
  bool set_error(int no, const char *msg) {
    if (error_no != 0)
      return false;
    // manual string copy since conflicting deprecations etc. on platforms. we
    // have a constant length so this should be safe.
    for (int i = 0; i < 1023; i++) {
      error_msg[i] = msg[i];
      if (msg[i] == 0)
        break;
    }
    error_no = no;
    return false;
  };

  struct cpt_error_t error() {
    return (struct cpt_error_t){error_no, &error_msg[0]};
  };

  // TODO: Make this one take reference instead.
  bool ensure_alignments(cpt_data d);

public:
  struct cpt_error_t dispatch(cpt_data d, int32_t nx, int32_t ny, int32_t nz) {
    generation++;
    if (this->error_no)
      return this->error();
    if (!this->ensure_alignments(d))
      return this->error();

    wg.add(nx * ny * nz);
    uvec3 nwg = make_uvec3(nx, ny, nz);
    for (uint32_t gz = 0; gz < (uint32_t)nz; ++gz) {
      for (uint32_t gy = 0; gy < (uint32_t)ny; ++gy) {
        for (uint32_t gx = 0; gx < (uint32_t)nx; ++gx) {
          queue.send(WorkPiece(nwg, make_uvec3(gx, gy, gz), &d, generation));
        }
      }
    }
    wg.wait();


    return this->error();
  }
};

#include "./generated/align.hpp"
