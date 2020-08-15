#pragma once

#include "arch/arch.hpp"
#include "state.hpp"
#include <cassert>

namespace co {
template <class T> class Routine {
public:
  State state;
  State *main_state;
  T arg;
  bool finished;

public:
  Routine() : state(0), finished(true){};
  Routine(int stack_size, State *ms)
      : state(stack_size), main_state(ms), finished(true){};
  ~Routine() { assert(finished); };

  // exit this thread and signal that it is done, a thread function
  // must never return without calling exit.
  void exit() {
    finished = true;
    yield();
  };

  bool is_finished() { return finished; }

  // run the given function in this routine, with the provided argument
  void set(void (*fp)(Routine<T> *), T arg) {
    assert(finished);
    finished = false;
    this->arg = arg;
    ARCH_set_register_state(state.registers, (void *)(fp), state.stack,
                            state.stack_size);
  }

  // get the argument that thus routine was started with
  T get_argument() { return this->arg; };

  // yield yields from the currently running routine to the next routine
  // in the series, or back to the initiator.
  void yield() { ARCH_switch(main_state, this); };

  void resume() { ARCH_switch(this, main_state); };
};
} // namespace co
