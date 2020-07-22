#pragma once

#include "generated/shader.hpp"
#include "routines/routines.hpp"
#include "types/types.hpp"

#define WG_SIZE (_cpt_WG_SIZE_Z * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X)

// wrapper for the main function in the computational shader
void thread_core(Invocation<struct shader*> *th) {
  auto td = th->get_argument();
  td->invocation = th;
  td->main();
  th->exit();
}

// the shader type generated from the compiled glsl shader has a barrier call,
// we need to forward that to the barrier call in the 
void shader::barrier() { 
  this->invocation->barrier();
}

class Kernel {
  WorkGroup<struct shader*> *wg;
  Invocation<struct shader*> *threads[WG_SIZE];
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
    this->wg = new WorkGroup<struct shader*>(static_cast<int>(num_t), 1024 * 128);
    for (int i = 0; i < WG_SIZE; i++) {
      threads[i] = wg->create_thread();
    }
    shared_data = shader::create_shared_data();
  }

  ~Kernel() {
    if (this->wg)
      delete this->wg;
    if (this->shared_data)
      shader::free_shared_data(shared_data);
    for (int i = 0; i < WG_SIZE; i++) {
      delete threads[i];
    }
  }

private:
  bool set_error(int no, const char *msg) {
    if (this->error_no != 0)
      return false;
    // manual string copy since conflicting deprecations etc. on platforms. we have a
    // constant length so this should be safe.
    for(int i = 0; i < 1023; i++) {
      this->error_msg[i] =  msg[i];
      if(msg[i] == 0)
       break;
    }
    this->error_no = no;
    return false;
  };

  struct cpt_error_t error() {
    return (struct cpt_error_t){this->error_no, &this->error_msg[0]};
  };

  bool ensure_alignments(cpt_data d);

  void dispatch_wg(uvec3 wgID, bool pause);

public:
  struct cpt_error_t dispatch(cpt_data d, int32_t nx, int32_t ny, int32_t nz);
};

void Kernel::dispatch_wg(uvec3 wgID, bool pause) {
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
        s->gl_LocalInvocationIndex = index;

        s->invocation = threads[index];
        threads[index]->set_function(&thread_core, &this->shaders[index]); // TODO: Why does this one need to be in here and not only done once outside global loop?
      }
    }
  }
  wg->run(WG_SIZE, this->threads, pause);
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
        dispatch_wg(make_uvec3(gx, gy, gz), gz == (uint32_t)nz-1 && gy == (uint32_t)ny-1 && gx == (uint32_t)nx-1);
      }
    }
  }

  return this->error();
}

#include "./generated/align.hpp"
