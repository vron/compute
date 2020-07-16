#pragma once

#include "routines.hpp"
#include <cassert>

template <class T> class Invocation {
  InvocationState state; // MUST be first member
  bool finished;
  T arg;
  WorkThread<T> *wt;

private:
  Invocation(int stack_size)
      : state(InvocationState(stack_size)), finished(true){};

public:
  // exit this thread and signal that it is done, a thread function
  // must never return without calling exit.
  void exit() {
    this->finished = true;
    this->barrier();
  };

  // re-configure the thread to run the given function with the given arg
  void set_function(void (*fp)(Invocation<T> *), T arg) {
    assert(this->finished);
    this->finished = false;
    this->arg = arg;
    ARCH_set_register_state(this->state.registers, (void *)fp,
                            this->state.stack, this->state.stack_size);
  }

  T get_argument() { return this->arg; };

  void barrier() {
    assert(this->wt != nullptr);
    ARCH_switch(&this->wt->state, &this->state);
  };

  friend class WorkThread<T>;
  friend class WorkGroup<T>;
};
