#pragma once

#include "routines.hpp"
#include <condition_variable>
#include <mutex>
#include <thread>

template <class T> class WorkGroup {
public:
  WorkQueue<T> sync;
  int num_thread;
  WorkThread<T> **workers;
  int stack_size;
  bool full_sync;

public:
  WorkGroup(int num_thread, int stack_size)
      : num_thread(num_thread), stack_size(stack_size), full_sync(true) {
    assert(num_thread > 0);
    this->workers =
        (WorkThread<T> **)malloc(num_thread * sizeof(WorkThread<T> *));
    for (int i = 0; i < num_thread; i++) {
      this->workers[i] = new WorkThread<T>(&this->sync);
    }
  }

  Invocation<T> *create_thread() { return (new Invocation<T>(stack_size)); };

  void run(int no, Invocation<T> **threads, bool pause) {
    // if pause == true do not use busy waiting, since we are not expecting to get
    // another wg dispatched straight away. Else do use busy waiting for performance.

    // simply split the work equally and split across threads. Downside of this
    // approach is that one thread might finish before. Upside is that we will largely
    // avoid false sharing and likely improve cache usage since subsequent invocations will
    // likely access subsequent array elements etc. Consider these two effects when / if
    // the scheduling is redone to something smarter.
    int nt = no < this->num_thread ? no : num_thread;
    int si = 0;
    int e_each = (no + (nt - 1)) / nt;
    for (int ti = 0; ti < nt; ti++) {
      // TODO: we should be able to replace this loop with single expression?
      int e = e_each;
      if (e + si > no) {
        e = no - si;
      }
      if (e < 1) {
        nt = ti;
        break;
      }
      si += e;
    }

    si = 0;
    sync.barrier.set_count(nt);
    for (int ti = 0; ti < nt; ti++) {
      int e = e_each;
      if (e + si > no) {
        e = no - si;
      }
      assert(e > 0);
      if(full_sync) {
        // the worker needs a full sync since it is not busy-waiting.
        workers[ti]->send_work(WorkPiece<T>(e, pause, &threads[si]));
      } else {
        // a partial sync means that the thread is busy-waiting, and that it will
        // assume that it should operate on exactly the same threads as last time so
        // we actually do not need to provide it with anything apart from saying
        // "do your work again".
        assert(false);
      }
      si += e;
    }
    for (int ti = 0; ti < nt; ti++) {
      workers[ti]->wait_for_done();
    }
  };

  ~WorkGroup() {
    // signal threads and wait for them to quit
    for (int i = 0; i < this->num_thread; i++) {
      this->workers[i]->send_work(WorkPiece<T>(0, true, nullptr));
    }
    for (int i = 0; i < this->num_thread; i++) {
      delete workers[i];
    }
    free(this->workers);
  }
};
