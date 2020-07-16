#pragma once

// TODO: investigate false sharing and thread locations (i thin we should use 64
// byte separation all over the place...)

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

public:
  WorkGroup(int num_thread, int stack_size)
      : num_thread(num_thread), stack_size(stack_size) {
    assert(num_thread > 0);
    this->workers =
        (WorkThread<T> **)malloc(num_thread * sizeof(WorkThread<T> *));
    for (int i = 0; i < num_thread; i++) {
      this->workers[i] = new WorkThread<T>(&this->sync);
    }
  }

  Invocation<T> *create_thread() { return (new Invocation<T>(stack_size)); };

  void run(int no, Invocation<T> **threads) {
    // simply split the work equally and split across threads. Downside of this
    // approach is that one thread might finish before. Upside is that we will largely
    // avoid false sharing and likely improve cache usage since subsequent invocations will
    // likely access subsequent array elements etc. Consider these two effects when / if
    // the scheduling is redone to something smarter.
    int nt = no < this->num_thread ? no : num_thread;
    int si = 0;
    int e_each = (no + (nt - 1)) / nt;
    sync.barrier.set_count(nt);
    for (int ti = 0; ti < nt; ti++) {
      int e = e_each;
      if (e + si > no) {
        e = no - si;
      }
      sync.send_work(WorkPiece<T>(e, &threads[si]));
      si += e;
    }

    for (int ti = 0; ti < nt; ti++) {
      sync.wait_for_done();
    }
  };

  ~WorkGroup() {
    // signal threads and wait for them to quit
    for (int i = 0; i < this->num_thread; i++) {
      this->sync.send_work(WorkPiece<T>(0, nullptr));
    }
    for (int i = 0; i < this->num_thread; i++) {
      delete workers[i];
    }
    free(this->workers);
  }
};
