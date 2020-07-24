#pragma once
#include "arch/arch.hpp"
#include <cassert>
#include <cstdlib>

namespace co {
class State {
public:
  void *registers[ARCH_reg_state];
  void *stack;
  size_t stack_size;

public:
  State(int stack_size) : stack_size(stack_size) {
    stack = nullptr;
    if (stack_size > 0) {
      stack = malloc(stack_size);
      assert(stack != nullptr);
      // TODO(vron): use guard page at end of stack to catch stack overflows?
    }
  }

  ~State() {
    if (stack != nullptr)
      free(stack);
  };
};
} // namespace co