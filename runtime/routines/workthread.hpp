#pragma once

#include "routines.hpp"
#include <thread>


/*
  A WorkThread uses an underlying OS thread to schedule shader invocations. It can be in one
  of the following 4 states at any point in time:

  A) Currently executing shader code (switching between multiple different ones)

  B) Busy-waiting in sync->barrier for other WorkThreads to reach the same barrier.

  C) Waiting in sync->wait_for_work for a WorkPiece to be executed
     this is the state for a newly created WorkThread

  D) Busy-waiting for a signal to run the same invocations again 
     this state is an optimization to avoid having to do the expensive sync in C between WGs
     in the same dispatch.
*/

// TODO: Padd this one to cache line size!
template <class T> class WorkThread {
public:
  InvocationState state;
  WorkQueue<T> *sync;
  std::thread t;

  std::mutex work_m;
  std::condition_variable work_cv;
  std::vector<WorkPiece<T>> work_queue;

  std::mutex done_m;
  std::condition_variable done_cv;
  std::vector<bool> done_queue;

public:
  WorkThread(WorkQueue<T> *sync) : state(InvocationState(1)), sync(sync) {
    this->t = std::thread(&WorkThread::loop, this);
  };
  ~WorkThread() { this->t.join(); };

  void resume_thread(Invocation<T> *t) {
    ARCH_switch(&t->state.registers, &this->state);
  };

   WorkPiece<T> wait_for_work() {
    std::unique_lock<std::mutex> lk(this->work_m);
    this->work_cv.wait(lk, [this] { return this->work_queue.size() > 0; });
    auto w = this->work_queue.back();
    this->work_queue.pop_back();
    return w;
  };

  void send_work(WorkPiece<T> w) {
    {
      std::lock_guard<std::mutex> lk(this->work_m);
      this->work_queue.push_back(w);
    }
    this->work_cv.notify_one();
  };

  bool wait_for_done() {
    std::unique_lock<std::mutex> lk(this->done_m);
    this->done_cv.wait(lk, [this] { return this->done_queue.size() > 0; });
    auto w = this->done_queue.back();
    this->done_queue.pop_back();
    return w;
  };

  void send_done(bool v) {
    {
      std::lock_guard<std::mutex> lk(this->done_m);
      this->done_queue.push_back(v);
    }
    this->done_cv.notify_one();
  };

private:
  void loop() {
    for (WorkPiece<T> work = this->wait_for_work(); work.threads != nullptr;
         work = this->wait_for_work()) {
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

      this->send_done(true);
    }
  }
};
