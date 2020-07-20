#pragma once

#include "debug.hpp"
#include "generated/shader.hpp"
#include "routines/routines.hpp"
#include "types/types.hpp"
#include <cstring>

#define WG_SIZE (_cpt_WG_SIZE_Z * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X)

// TODO: Can we remove this one and move to just shader?
class WorkGroupArg {
public:
  shader *shader;
  WorkGroupArg(){};
  WorkGroupArg(struct shader *s) : shader(s){};
};

void thread_core(Invocation<WorkGroupArg> *th) {
  auto td = th->get_argument();
  td.shader->thread = th;
  td.shader->main();
  th->exit();
}

void shader::barrier() { 
  this->thread->barrier();
  }

class Kernel {
  WorkGroup<WorkGroupArg> *wg;
  Invocation<WorkGroupArg> *threads[WG_SIZE];
  shader shaders[WG_SIZE];
  shared_data_t *shared_data;

  char error_msg[1024];
  int error_no;

public:
  Kernel(int32_t num_t) {
    if (num_t < 1) {
      this->set_error(EINVAL, "must use at least 1 thread");
      return;
    }
    // use 128B stack as default, this is likely much larger then ever needed
    this->wg = new WorkGroup<WorkGroupArg>(static_cast<int>(num_t), 1024 * 128);
    for (int i = 0; i < WG_SIZE; i++) {
      threads[i] = wg->create_thread(); // TODO: Are we leaing threads here that
                                        // we are not collecting?
    }
    shared_data = shader::create_shared_data();
  }

  ~Kernel() {
    if (this->wg)
      delete this->wg;
    if (this->shared_data)
      shader::free_shared_data(shared_data);
  }

private:
  bool set_error(int no, const char *msg) {
    if (this->error_no != 0)
      return false;
    strncpy(&this->error_msg[0], msg, 1023);
    this->error_no = no;
    return false;
  };

  struct cpt_error_t error() {
    return (struct cpt_error_t){this->error_no, &this->error_msg[0]};
  };

  bool ensure_alignments(cpt_data d);

  void dispatch_wg(uvec3 wgID);

public:
  struct cpt_error_t dispatch(cpt_data d, int32_t nx, int32_t ny, int32_t nz);
};

void Kernel::dispatch_wg(uvec3 wgID) {
  int index = -1;
  for (uint32_t lz = 0; lz < _cpt_WG_SIZE_Z; ++lz) {
    for (uint32_t ly = 0; ly < _cpt_WG_SIZE_Y; ++ly) {
      for (uint32_t lx = 0; lx < _cpt_WG_SIZE_X; ++lx) {
        index++;
        shader *s = &this->shaders[index];
        s->gl_WorkGroupID = wgID;
        s->gl_LocalInvocationID = make_uvec3(lx, ly, lz);
        s->gl_GlobalInvocationID =
            s->gl_WorkGroupID * s->gl_WorkGroupSize + s->gl_LocalInvocationID;
        s->gl_LocalInvocationIndex =
            lx + ly * _cpt_WG_SIZE_X +
            lz * _cpt_WG_SIZE_X * _cpt_WG_SIZE_Y; // TODO: replace with index

        s->thread = threads[index];
        threads[index]->set_function(&thread_core, WorkGroupArg(&this->shaders[index]));
      }
    }
  }
  wg->run(WG_SIZE, this->threads);
}

struct cpt_error_t Kernel::dispatch(cpt_data d, int32_t nx, int32_t ny,
                                    int32_t nz) {
  if (this->error_no)
    return this->error();
  if (!this->ensure_alignments(d))
    return this->error();

  // set up constants that will be same for all
  uvec3 nwg = make_uvec3(nx, ny, nz);
  uvec3 wgs = make_uvec3(_cpt_WG_SIZE_X, _cpt_WG_SIZE_Y, _cpt_WG_SIZE_Z);
  for (int i = 0; i < WG_SIZE; i++) {
    this->shaders[i].gl_NumWorkGroups = nwg;
    this->shaders[i].gl_WorkGroupSize = wgs;
    this->shaders[i].set_data(d);
    this->shaders[i].set_shared_data(shared_data);
  }

  for (uint32_t gz = 0; gz < (uint32_t)nz; ++gz) {
    for (uint32_t gy = 0; gy < (uint32_t)ny; ++gy) {
      for (uint32_t gx = 0; gx < (uint32_t)nx; ++gx) {
        dispatch_wg(make_uvec3(gx, gy, gz));
      }
    }
  }

  return this->error();
}

#include "./generated/align.hpp"
