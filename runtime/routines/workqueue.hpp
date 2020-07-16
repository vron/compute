#pragma once

#include "routines.hpp"
#include <condition_variable>
#include <mutex>
#include <vector>

template <class T> class WorkQueue {
public:
  std::mutex work_m;
  std::condition_variable work_cv;
  std::vector<WorkPiece<T>> work_queue;

  Barrier barrier;

  std::mutex done_m;
  std::condition_variable done_cv;
  std::vector<bool> done_queue;

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
};
