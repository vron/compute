#pragma once

#include "routines.hpp"
#include <thread>

template <class T> class WorkThread {
public:
  InvocationState state;
  WorkQueue<T> *sync;
  std::thread t;

public:
  WorkThread(WorkQueue<T> *sync) : state(InvocationState(1)), sync(sync) {
    this->t = std::thread(&WorkThread::loop, this);
  };
  ~WorkThread() { this->t.join(); };

  void resume_thread(Invocation<T> *t) {
    ARCH_switch(&t->state.registers, &this->state);
  };

private:
  void loop() {
    for (WorkPiece<T> work = sync->wait_for_work(); work.threads != nullptr;
         work = sync->wait_for_work()) {
      for (int i = 0; i < work.no; i++) {
        work.threads[i]->wt = this;
      }

      // note that as per the glsl compute specification we can reply on all invocations
      // reaching the same barrier calls in the same order. If this is no ensured by the
      // author we will get strange results here...
      do {
        for (int i = 0; i < work.no; i++) {
          this->resume_thread(work.threads[i]);
        }
        if (work.threads[0]->finished)
          break;
        // barrier calls also need to be sync:ed with those in other threads.
        sync->barrier.wait();

      } while (!work.threads[0]->finished);

      sync->send_done(true);
    }
  }
};
