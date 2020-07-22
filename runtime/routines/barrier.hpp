#pragma once
#include <condition_variable>
#include <mutex>
#include <cstdint>

// Barrier uses busy-wait to improve performance of sync between threads in the same
// work group.
class Barrier {
private:
// we store both the count and the generation in the same variable
// such that we can load and store both in one atomic operation.
  uint32_t count;
  uint32_t original_count;

public:
  Barrier() : count(0), original_count(0){};

  // the caller must ensure that there us synchronization betwen a call to set and any call to wait.
  void set_count(int c) {
    __atomic_store_n(&count, c, __ATOMIC_SEQ_CST);
    __atomic_store_n(&original_count, c, __ATOMIC_SEQ_CST);
  }

  void wait() {
    //return;
    uint32_t val = __atomic_fetch_sub(&count, 1, __ATOMIC_SEQ_CST);
    uint32_t my_count = val & 0xFFFF;
    uint32_t my_gen = val & 0xFFFF0000;
    if ( my_count == 1) {
      // we are the last one to reach the barrier, to signal the others
      my_gen += 1 << 0x10;
      uint32_t oc = __atomic_load_n(&original_count, __ATOMIC_SEQ_CST);
      oc += my_gen;
      __atomic_store_n(&count, oc, __ATOMIC_SEQ_CST);
      return;
    } else {
      uint32_t wait_gen;
      do {
        #if defined(__x86_64__)
        __asm__ ("rep nop"); // TODO: move this into arch package?
        #else
        #error "unsupported architecture"
        #endif
        uint32_t wait_val = __atomic_load_n(&this->count, __ATOMIC_SEQ_CST);
        wait_gen = wait_val & 0xFFFF0000;
      } while(wait_gen == my_gen);

      return;
    }
  };
};
