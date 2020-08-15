#pragma once

#include "co/routines.hpp"
#include "generated/shader.hpp"
#include "types/types.hpp"
#include "workpiece.hpp"
#include <thread>

cptc_data global_uninit;

void invocation_core(co::Routine<struct shader *> *invocation) {
  auto shader = invocation->get_argument();
  shader->invocation = invocation;
  shader->main();
  invocation->exit();
}

void shader::barrier() { this->invocation->yield(); }

template <class T> class WorkThread {
public:
  co::State state;
  WorkQueue<WorkPiece> *queue;
  WaitGroup<int32_t> *wg;
  std::thread t;
  int stack_size;
  int generation;

  co::Routine<T> *routines[_cpt_WG_SIZE]; // tODO: non pointers?
  shared_data_t shared_data;
  shader *shaders[_cpt_WG_SIZE]; // TODO: can we somehow remove the need for one
                                // object per invocation? a lot of vectors eat
                                // space here... At least eep it on the
                                // invocation's stac?

public:
  WorkThread(WorkQueue<WorkPiece> *queue, WaitGroup<int32_t> *wg,
             int stack_size)
      : state(co::State(0)), queue(queue), wg(wg), stack_size(stack_size),
        generation(0), shaders{_cpt_REPEAT_WG_SIZE(nullptr)} {
    this->t = std::thread(&WorkThread::run, this);
  };

  ~WorkThread() {
    this->t.join();
  };

private:
  void setup() {
    for (int i = _cpt_WG_SIZE - 1; i >= 0; i--) {
      routines[i] = (new co::Routine<T>(
          stack_size,
          &state)); // TODO: THis will obviously leak memory? or will it?
    }
  }

  void finish() {
    for (int i = _cpt_WG_SIZE - 1; i >= 0; i--) {
          delete routines[i];
    }
  }

  void run() {
    setup();

    for (WorkPiece wp = queue->receive(); !wp.quit;
         wp = queue->receive()) {
      assert(wp.data != nullptr);
      for (int i = 0; i < _cpt_WG_SIZE; i++) {
        shaders[i] = new shader(wp.data, &shared_data);
        shaders[i]->gl_NumWorkGroups = wp.nwg;
        shaders[i]->invocation = routines[i];
        shaders[i]->gl_WorkGroupSize = make_uvec3(_cpt_WG_SIZE_X, _cpt_WG_SIZE_Y, _cpt_WG_SIZE_Z);
      }

      int index = -1;
      for (uint32_t lz = 0; lz < _cpt_WG_SIZE_Z; ++lz) {
        for (uint32_t ly = 0; ly < _cpt_WG_SIZE_Y; ++ly) {
          for (uint32_t lx = 0; lx < _cpt_WG_SIZE_X; ++lx) {
            // TODO: can we limit the amount of data that needs calculation
            // here? Increment them sequentially? - Actually all of these are
            // constants between threads and routines - Do this one and not
            // multiple times?
            // Also! we can parse if they are used or not, maybe most of them
            // need not be created at all (nor included in the shader class)
            index++;
            shaders[index]->gl_WorkGroupID = wp.gid;
            shaders[index]->gl_LocalInvocationID = make_uvec3(lx, ly, lz);
            shaders[index]->gl_GlobalInvocationID =
                shaders[index]->gl_WorkGroupID *
                    shaders[index]->gl_WorkGroupSize +
                shaders[index]->gl_LocalInvocationID;
            shaders[index]->gl_LocalInvocationIndex = index;
          }
        }
      }

      // must be re-set since it needs to start from beg.
      for (int i = _cpt_WG_SIZE - 1; i >= 0; i--) {
        routines[i]->set(&invocation_core, shaders[i]);
      }
      // note that as per the glsl compute specification we can reply on all
      // routines reaching the same barrier calls in the same order. If this is
      // no ensured by the author we will get strange results here...
      do {
        for (int i = 0; i < _cpt_WG_SIZE; i++) {
          // TODO: Think about: this call effectively makes it impossible to
          // vectorize beween kernels - can we SIMD them, ref webrender?
          routines[i]->resume();
        }
      } while (!routines[0]->is_finished());


      for (int i = 0; i < _cpt_WG_SIZE; i++) {
        delete(shaders[i]);
      }

      wg->done();
    }
    finish();
  }
};
