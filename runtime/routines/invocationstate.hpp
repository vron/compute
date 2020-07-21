#pragma once
#include "routines.hpp"
#include <cassert>
#include <cstdint>
#include <cstdlib>

class InvocationState {
public:
  void *registers[ARCH_reg_state];
  void *stack;
  size_t stack_size;

public:
  InvocationState(int stack_size) {
    assert(stack_size > 0);

    // zero out the register state
    for (int i = 0; i < ARCH_reg_state; i++) {
      this->registers[i] = 0;
    }

    this->stack_size = stack_size;
    this->stack = malloc(stack_size);
    // TODO(vron): use guard page at end of stack to catch stack overflows
  }

  ~InvocationState() {
    if (!this->stack)
      return;
    free(this->stack);
  };
};
